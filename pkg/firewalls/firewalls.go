package firewalls

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/pkg/extra/godojs"
	"github.com/aybabtme/godotto/pkg/extra/ottoutil"
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
		{"update", svc.update},
		{"add_tags", svc.addTags},
		{"remove_tags", svc.removeTags},
		{"add_droplets", svc.addDroplets},
		{"remove_droplets", svc.removeDroplets},
		{"add_rules", svc.addRules},
		{"remove_rules", svc.removeRules},
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

func (svc *firewallsSvc) update(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	req := godojs.ArgFirewallUpdate(vm, all.Argument(1))
	f, err := svc.svc.Update(svc.ctx, fwID, firewalls.UseGodoFirewall(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.FirewallToVM(vm, f.Struct())
}

func (svc *firewallsSvc) addTags(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	tags := godojs.ArgTags(vm, all.Argument(1))

	err := svc.svc.AddTags(svc.ctx, fwID, tags...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) removeTags(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	tags := godojs.ArgTags(vm, all.Argument(1))

	err := svc.svc.RemoveTags(svc.ctx, fwID, tags...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) addDroplets(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	dropletIDs := godojs.ArgDropletIDs(vm, all.Argument(1))

	err := svc.svc.AddDroplets(svc.ctx, fwID, dropletIDs...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) removeDroplets(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	dropletIDs := godojs.ArgDropletIDs(vm, all.Argument(1))

	err := svc.svc.RemoveDroplets(svc.ctx, fwID, dropletIDs...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) addRules(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	inboundRules := godojs.ArgInboundRules(vm, all.Argument(1))
	outboundRules := godojs.ArgOutboundRules(vm, all.Argument(2))

	err := svc.svc.AddRules(svc.ctx, fwID, inboundRules, outboundRules)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *firewallsSvc) removeRules(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	fwID := godojs.ArgFirewallID(vm, all.Argument(0))
	inboundRules := godojs.ArgInboundRules(vm, all.Argument(1))
	outboundRules := godojs.ArgOutboundRules(vm, all.Argument(2))

	err := svc.svc.RemoveRules(svc.ctx, fwID, inboundRules, outboundRules)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}
