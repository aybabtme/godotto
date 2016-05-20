package floatingips

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := floatingIPSvc{
		ctx: ctx,
		svc: client.FloatingIPs(),
	}

	actions, err := applyAction(ctx, vm, client)
	if err != nil {
		return q, err
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"list", svc.list},
		{"create", svc.create},
		{"get", svc.get},
		{"delete", svc.delete},
		{"actions", actions},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type floatingIPSvc struct {
	ctx context.Context
	svc floatingips.Client
}

func (svc *floatingIPSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	req := godojs.ArgFloatingIPCreateRequest(vm, all.Argument(0))
	fip, err := svc.svc.Create(svc.ctx, req.Region, floatingips.UseGodoFloatingIP(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.FloatingIPToVM(vm, fip.Struct())
}

func (svc *floatingIPSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	ip := godojs.ArgFloatingIPActualIP(vm, all.Argument(0))
	fip, err := svc.svc.Get(svc.ctx, ip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.FloatingIPToVM(vm, fip.Struct())
}

func (svc *floatingIPSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	ip := godojs.ArgFloatingIPActualIP(vm, all.Argument(0))

	err := svc.svc.Delete(svc.ctx, ip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *floatingIPSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var floatingIPs = make([]otto.Value, 0)
	floatingIPc, errc := svc.svc.List(svc.ctx)
	for d := range floatingIPc {
		floatingIPs = append(floatingIPs, godojs.FloatingIPToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(floatingIPs)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
