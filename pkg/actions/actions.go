package actions

import (
	"fmt"
	"context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/actions"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := actionSvc{
		ctx: ctx,
		svc: client.Actions(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"get", svc.get},
		{"list", svc.list},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type actionSvc struct {
	ctx context.Context
	svc actions.Client
}

func (svc *actionSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	aid := godojs.ArgActionID(vm, arg)
	a, err := svc.svc.Get(svc.ctx, aid)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.ActionToVM(vm, a.Struct())
}

func (svc *actionSvc) list(all otto.FunctionCall) otto.Value {

	vm := all.Otto
	var actions = make([]otto.Value, 0)
	actionc, errc := svc.svc.List(svc.ctx)
	for action := range actionc {
		actions = append(actions, godojs.ActionToVM(vm, action.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(actions)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
