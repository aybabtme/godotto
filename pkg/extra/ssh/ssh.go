package ssh

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, auth ssh.AuthMethod) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := sshSvc{
		ctx:  ctx,
		auth: auth,
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"session", svc.session},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}
	return root.Value(), nil
}

type sshSvc struct {
	ctx  context.Context
	auth ssh.AuthMethod
}

type connectOpts struct {
	Hostname string
	Port     string
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
	case v.IsObject():
		host = ottoutil.String(vm, ottoutil.GetObject(vm, v.Object(), "public_ipv4"))
		if host == "" {
			ottoutil.Throw(vm, "provided Droplet has no public IPv4")
		}
		slug := ottoutil.String(vm, ottoutil.GetObject(vm, v.Object(), "image_slug"))
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
	obj := v.Object()
	if user := ottoutil.String(vm, ottoutil.GetObject(vm, obj, "user")); user != "" {
		opts.Cfg.User = user
	}
	if port := ottoutil.String(vm, ottoutil.GetObject(vm, obj, "port")); port != "" {
		opts.Port = port
	}
	return opts
}

func (svc *sshSvc) connect(ctx context.Context, opts *connectOpts) (*ssh.Client, error) {
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
			return ssh.NewClient(sconn, sc, rr), nil
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

	return ottoutil.ToPkg(vm, map[string]func(otto.FunctionCall) otto.Value{
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
