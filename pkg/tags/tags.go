package tags

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/tags"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := tagSvc{
		ctx: ctx,
		svc: client.Tags(),
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"create", svc.create},
		{"get", svc.get},
		{"list", svc.list},
		{"delete", svc.delete},
		{"tag_resources", svc.tag_resources},
		{"untag_resources", svc.untag_resources},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type tagSvc struct {
	ctx context.Context
	svc tags.Client
}

func (svc *tagSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgTagCreateRequest(vm, arg)

	t, err := svc.svc.Create(svc.ctx, req.Name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.TagToVM(vm, t.Struct())
}

func (svc *tagSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	tag := ottoutil.String(vm, arg)

	t, err := svc.svc.Get(svc.ctx, tag)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.TagToVM(vm, t.Struct())
}

func (svc *tagSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var tags = make([]otto.Value, 0)
	tagc, errc := svc.svc.List(svc.ctx)
	for t := range tagc {
		tags = append(tags, godojs.TagToVM(vm, t.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(tags)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}

func (svc *tagSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	tag := ottoutil.String(vm, arg)
	err := svc.svc.Delete(svc.ctx, tag)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *tagSvc) tag_resources(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgTagTagResourcesRequest(vm, arg)
	name, err := ottoutil.GetObject(vm, arg, "name", true).ToString()
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	err = svc.svc.TagResources(svc.ctx, name, req.Resources)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}

func (svc *tagSvc) untag_resources(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgTagUntagResourcesRequest(vm, arg)
	name, err := ottoutil.GetObject(vm, arg, "name", true).ToString()
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	err = svc.svc.UntagResources(svc.ctx, name, req.Resources)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return q
}
