package firewalls

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/firewalls"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := firewallsSvc{
		ctx: ctx,
		svc: client.Firewalls(),
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"create", svc.create},
		{"get", svc.get},
		{"delete", svc.delete},
		{"list", svc.list},
	} {

		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type firewallsSvc struct {
	ctx context.Context
	svc firewalls.Client
}

func (svc *firewallsSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgFirewallCreate(vm, arg)

	f, err := svc.svc.Create(svc.ctx, req.Name, req.InboundRules, req.OutboundRules, firewalls.UseGodoCreate(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.FirewallToVM(vm, f.Struct())
}

func (svc *firewallsSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	fwID := godojs.ArgFirewallID(vm, arg)
	f, err := svc.svc.Get(svc.ctx, fwID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.FirewallToVM(vm, f.Struct())
}

func (svc *firewallsSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	fwID := godojs.ArgFirewallID(vm, arg)
	err := svc.svc.Delete(svc.ctx, fwID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var fws = make([]otto.Value, 0)
	fwc, errc := svc.svc.List(svc.ctx)
	for f := range fwc {
		fws = append(fws, godojs.FirewallToVM(vm, f.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(fws)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}
