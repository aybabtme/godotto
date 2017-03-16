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
