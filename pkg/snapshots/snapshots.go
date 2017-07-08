package snapshots

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/snapshots"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := snapshotsSvc{
		ctx: ctx,
		svc: client.Snapshots(),
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"get", svc.get},
		{"list", svc.list},
		{"list_droplet", svc.listDroplet},
		{"list_volume", svc.listVolume},
		{"delete", svc.delete},
	} {

		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type snapshotsSvc struct {
	ctx context.Context
	svc snapshots.Client
}

func (svc *snapshotsSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	sId := godojs.ArgSnapshotID(vm, arg)
	s, err := svc.svc.Get(svc.ctx, sId)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.SnapshotToVM(vm, s.Struct())
}

func (svc *snapshotsSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	sId := godojs.ArgSnapshotID(vm, arg)

	err := svc.svc.Delete(svc.ctx, sId)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *snapshotsSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var ss = make([]otto.Value, 0)
	sc, errc := svc.svc.List(svc.ctx)
	for s := range sc {
		ss = append(ss, godojs.SnapshotToVM(vm, s.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(ss)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}

func (svc *snapshotsSvc) listDroplet(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var ss = make([]otto.Value, 0)
	sc, errc := svc.svc.ListDroplet(svc.ctx)
	for s := range sc {
		ss = append(ss, godojs.SnapshotToVM(vm, s.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(ss)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}

func (svc *snapshotsSvc) listVolume(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var ss = make([]otto.Value, 0)
	sc, errc := svc.svc.ListVolume(svc.ctx)
	for s := range sc {
		ss = append(ss, godojs.SnapshotToVM(vm, s.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(ss)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}
