package droplets

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/pkg/extra/godojs"
	"github.com/aybabtme/godotto/pkg/extra/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := dropletSvc{
		ctx: ctx,
		svc: client.Droplets(),
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
		{"get", svc.get},
		{"create", svc.create},
		{"create_multiple", svc.createMultiple},
		{"delete", svc.delete},
		{"actions", actions},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type dropletSvc struct {
	ctx context.Context
	svc droplets.Client
}

func (svc *dropletSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgDropletCreateRequest(vm, arg)

	d, err := svc.svc.Create(svc.ctx, req.Name, req.Region, req.Size, req.Image.Slug, droplets.UseGodoCreate(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DropletToVM(vm, d.Struct())
}

func (svc *dropletSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	did := godojs.ArgDropletID(vm, arg)

	d, err := svc.svc.Get(svc.ctx, did)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.DropletToVM(vm, d.Struct())
}

func (svc *dropletSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	did := godojs.ArgDropletID(vm, arg)

	err := svc.svc.Delete(svc.ctx, did)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *dropletSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var droplets = make([]otto.Value, 0)
	dropletc, errc := svc.svc.List(svc.ctx)
	for d := range dropletc {
		droplets = append(droplets, godojs.DropletToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(droplets)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *dropletSvc) createMultiple(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	arg := all.Argument(0)

	req := godojs.ArgDropletMultiCreateRequest(vm, arg)

	droplets, err := svc.svc.CreateMultiple(svc.ctx, req.Names, req.Region, req.Size, req.Image.Slug, droplets.UseGodoMultiCreate(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	var d = make([]otto.Value, 0, len(droplets))
	for _, droplet := range droplets {
		d = append(d, godojs.DropletToVM(vm, droplet.Struct()))
	}

	v, err := vm.ToValue(d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}
