package accounts

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/pkg/extra/godojs"
	"github.com/aybabtme/godotto/pkg/extra/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/accounts"
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
	return godojs.AccountToVM(vm, a.Struct())
}
