package ssh

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, auth ssh.AuthMethod) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := sshSvc{
		auth: auth,
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"shell", svc.shell},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}
	return root.Value(), nil
}

type sshSvc struct {
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
		host = ottoutil.GetObject(vm, v.Object(), "public_ipv4")

		slug := ottoutil.GetObject(vm, v.Object(), "image_slug")
		switch {
		case strings.Contains(slug, "coreos"):
			user = "core"
		case strings.Contains(slug, "freebsd"):
			user = "freebsd"
		}
	default:
		ottoutil.Throw(vm, "argument must be a string or a Droplet")
	}

	return &connectDropletOpts{
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

func (svc *sshSvc) shell(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	opts := svc.connectArgs(vm, all.Argument(0))
	switch len(all.ArgumentList) {
	case 1: // done
	case 2: // provided options
		opts = svc.optionalConnectArgs(vm, opts, all.Argument(1))
	default: // too many!
		ottoutil.Throw(vm, "too many arguments")
	}
	client, err := svc.connect(context.TODO(), opts)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	run = func(all otto.FunctionCall) otto.Value {
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
		return ottoutil.ToValue(string(out))
	}
	return ottoutil.ToAnonFunc(vm, run)
}
