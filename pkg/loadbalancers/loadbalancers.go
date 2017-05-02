package loadbalancers

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/loadbalancers"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := loadBalancersSvc{
		ctx: ctx,
		svc: client.LoadBalancers(),
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"create", svc.create},
		{"get", svc.get},
		{"list", svc.list},
		{"delete", svc.delete},
		{"add_droplets", svc.addDroplets},
		{"remove_droplets", svc.removeDroplets},
		{"add_forwarding_rules", svc.addForwardingRules},
		{"remove_forwarding_rules", svc.removeForwardingRules},
	} {

		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type loadBalancersSvc struct {
	ctx context.Context
	svc loadbalancers.Client
}

func (svc *loadBalancersSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgLoadBalancerCreateRequest(vm, arg)

	l, err := svc.svc.Create(svc.ctx, req.Name, req.Region, req.ForwardingRules, loadbalancers.UseGodoCreate(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.LoadBalancerToVM(vm, l.Struct())
}

func (svc *loadBalancersSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	lbId := godojs.ArgLoadBalancerID(vm, arg)
	l, err := svc.svc.Get(svc.ctx, lbId)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.LoadBalancerToVM(vm, l.Struct())
}

func (svc *loadBalancersSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	lbId := godojs.ArgLoadBalancerID(vm, arg)

	err := svc.svc.Delete(svc.ctx, lbId)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *loadBalancersSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var lbs = make([]otto.Value, 0)
	lbc, errc := svc.svc.List(svc.ctx)
	for l := range lbc {
		lbs = append(lbs, godojs.LoadBalancerToVM(vm, l.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(lbs)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}

func (svc *loadBalancersSvc) addDroplets(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	lbId := godojs.ArgLoadBalancerID(vm, all.Argument(0))
	dropletIds := godojs.ArgDropletIDs(vm, all.Argument(1))

	err := svc.svc.AddDroplets(svc.ctx, lbId, dropletIds...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *loadBalancersSvc) removeDroplets(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	lbId := godojs.ArgLoadBalancerID(vm, all.Argument(0))
	dropletIds := godojs.ArgDropletIDs(vm, all.Argument(1))

	err := svc.svc.RemoveDroplets(svc.ctx, lbId, dropletIds...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *loadBalancersSvc) addForwardingRules(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	lbId := godojs.ArgLoadBalancerID(vm, all.Argument(0))
	rules := godojs.ArgForwardingRules(vm, all.Argument(1))

	err := svc.svc.AddForwardingRules(svc.ctx, lbId, rules...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *loadBalancersSvc) removeForwardingRules(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	lbId := godojs.ArgLoadBalancerID(vm, all.Argument(0))
	rules := godojs.ArgForwardingRules(vm, all.Argument(1))

	err := svc.svc.RemoveForwardingRules(svc.ctx, lbId, rules...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}
