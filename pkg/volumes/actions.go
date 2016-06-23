package volumes

import (
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/volumes"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/context"
)

func applyAction(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := actionSvc{
		ctx: ctx,
		svc: client.Volumes().Actions(),
	}
	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"attach", svc.attach},
		{"detach", svc.detach},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type actionSvc struct {
	ctx context.Context
	svc volumes.ActionClient
}

func (svc *actionSvc) attach(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	ip := godojs.ArgVolumeID(vm, all.Argument(0))
	dropletID := godojs.ArgDropletID(vm, all.Argument(1))
	err := svc.svc.Attach(svc.ctx, ip, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) detach(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := godojs.ArgVolumeID(vm, all.Argument(0))
	err := svc.svc.Detach(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}
