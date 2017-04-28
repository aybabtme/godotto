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
	root, err := vm.Object(`{}`)
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
