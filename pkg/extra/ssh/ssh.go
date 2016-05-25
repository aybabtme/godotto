package ssh

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, auth ssh.AuthMethod) (v otto.Value, cleanup func(), err error) {
	var qdn = func() {}
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, qdn, err
	}

	svc := sshSvc{
		ctx:    ctx,
		auth:   auth,
		opened: make(map[*ssh.Client]struct{}),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"session", svc.session},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, qdn, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}
	cleanup = func() {
		svc.mu.Lock()
		defer svc.mu.Unlock()
		for client := range svc.opened {
			_ = client.Close()
		}
	}

	return root.Value(), cleanup, nil
}

type sshSvc struct {
	ctx  context.Context
	auth ssh.AuthMethod

	mu     sync.Mutex
	opened map[*ssh.Client]struct{}
}

type connectOpts struct {
	Hostname string
	Port     string
	Timeout  time.Duration
	Cfg      *ssh.ClientConfig
}

func (svc *sshSvc) connectArgs(vm *otto.Otto, v otto.Value) *connectOpts {
	var (
		host string
		user = "root"
	)
	switch {
	case v.IsString():
		host, _ = v.ToString()
		if host == "" {
			ottoutil.Throw(vm, "no hostname provided")
		}
	case v.IsObject():
		droplet := godojs.ArgDroplet(vm, v)
		var err error
		host, err = droplet.PublicIPv4()
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		if host == "" {
			ottoutil.Throw(vm, "provided Droplet has no public IPv4")
		}
		slug := droplet.Image.Slug
		switch {
		case strings.Contains(slug, "coreos"):
			user = "core"
		case strings.Contains(slug, "freebsd"):
			user = "freebsd"
		}
	default:
		ottoutil.Throw(vm, "argument must be a string or a Droplet")
	}
	return &connectOpts{
		Hostname: host,
		Port:     "22",
		Cfg: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{svc.auth},
		},
	}
}

func (svc *sshSvc) optionalConnectArgs(vm *otto.Otto, opts *connectOpts, v otto.Value) *connectOpts {
	if v.Object() == nil {
		ottoutil.Throw(vm, "optional arguments must be an Object")
	}
	if user := ottoutil.String(vm, ottoutil.GetObject(vm, v, "user", false)); user != "" {
		opts.Cfg.User = user
	}
	if port := ottoutil.String(vm, ottoutil.GetObject(vm, v, "port", false)); port != "" {
		opts.Port = port
	}
	if dur := ottoutil.GetObject(vm, v, "timeout", false); dur.IsDefined() {
		if timeout := ottoutil.Duration(vm, dur); timeout != 0 {
			opts.Timeout = timeout
		}
	}
	return opts
}

func (svc *sshSvc) connect(ctx context.Context, opts *connectOpts) (*ssh.Client, error) {
	if opts.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	addr := net.JoinHostPort(opts.Hostname, opts.Port)
	var err error
	for {
		conn, derr := net.DialTimeout("tcp", addr, 2*time.Second)
		if derr != nil {
			if retryable(derr) {
				continue
			}
			return nil, fmt.Errorf("unexpected network error: %v", derr)
		}
		sconn, sc, rr, cerr := ssh.NewClientConn(conn, addr, opts.Cfg)
		if cerr == nil {
			client := ssh.NewClient(sconn, sc, rr)
			svc.mu.Lock()
			defer svc.mu.Unlock()
			svc.opened[client] = struct{}{}
			return client, nil
		}
		_ = conn.Close()
		err = fmt.Errorf("can't ssh into address %q, %v", addr, cerr)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, err
		}
	}
}

func (svc *sshSvc) session(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	opts := svc.connectArgs(vm, all.Argument(0))
	switch len(all.ArgumentList) {
	case 1: // done
	case 2: // provided options
		opts = svc.optionalConnectArgs(vm, opts, all.Argument(1))
	default: // too many!
		ottoutil.Throw(vm, "too many arguments")
	}
	client, err := svc.connect(svc.ctx, opts)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"exec": func(all otto.FunctionCall) otto.Value {
			vm := all.Otto
			cmd := ottoutil.String(vm, all.Argument(0))

			ss, err := client.NewSession()
			if err != nil {
				_ = client.Close()
				ottoutil.Throw(vm, err.Error())
			}
			out, err := ss.CombinedOutput(cmd)
			if err != nil {
				ottoutil.Throw(vm, "%v: %s", err, string(out))
			}
			_ = ss.Close()
			return ottoutil.ToValue(vm, string(out))
		},
		"close": func(all otto.FunctionCall) otto.Value {
			vm := all.Otto
			if err := client.Close(); err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			svc.mu.Lock()
			defer svc.mu.Unlock()
			delete(svc.opened, client)
			return q
		},
	})
}

// errors

var knownFailureSuffixes = []string{
	"connection refused",
	"connection reset by peer.",
	"connection timed out.",
	"connection timed out", // inconsistent standard library
	"no such host",
	"remote error: handshake failure",
	"unexpected EOF.",
	"use of closed network connection",
	"request canceled while waiting for connection",
	"read/write on closed pipe",
	"unexpected EOF reading trailer",
}

func hasRetryableSuffix(err error) bool {
	s := err.Error()
	for _, suffix := range knownFailureSuffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func retryable(err error) bool {
	if hasRetryableSuffix(err) {
		return true
	}

	switch e := err.(type) {
	case temporaryAndTimeoutError:
		return e.Temporary() || e.Timeout()
	case timeoutError:
		return e.Timeout()
	case temporaryError:
		return e.Temporary()
	case retryableError:
		return e.Retry()
	case *url.Error:
		return retryable(e.Err)
	default:
		return false
	}
}

type temporaryAndTimeoutError interface {
	Temporary() bool
	Timeout() bool
	Error() string
}

type timeoutError interface {
	Timeout() bool
	Error() string
}

type temporaryError interface {
	Temporary() bool
	Error() string
}

type retryableError interface {
	Retry() bool
	Error() string
}
