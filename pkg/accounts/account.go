package accounts

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/accounts"
	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := accountSvc{
		ctx: ctx,
		svc: client.Accounts(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"get", svc.get},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type accountSvc struct {
	ctx context.Context
	svc accounts.Client
}

func (svc *accountSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	a, err := svc.svc.Get(svc.ctx)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.accountToVM(vm, *a.Struct())
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *accountSvc) accountToVM(vm *otto.Otto, g godo.Account) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"droplet_limit", g.DropletLimit},
		{"floating_ip_limit", g.FloatingIPLimit},
		{"email", g.Email},
		{"uuid", g.UUID},
		{"email_verified", g.EmailVerified},
		{"status", g.Status},
		{"status_message", g.StatusMessage},
	} {
		v, err := vm.ToValue(field.v)
		if err != nil {
			return q, fmt.Errorf("can't prepare field %q: %v", field.name, err)
		}
		if err := d.Set(field.name, v); err != nil {
			return q, fmt.Errorf("can't set field %q: %v", field.name, err)
		}
	}
	return d.Value(), nil
}
