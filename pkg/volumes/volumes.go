package volumes

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/volumes"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := volumeSvc{
		ctx: ctx,
		svc: client.Volumes(),
	}

	actions, err := applyAction(ctx, vm, client)
	if err != nil {
		return q, err
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"list_volumes", svc.listVolume},
		{"get_volume", svc.getVolume},
		{"create_volume", svc.createVolume},
		{"delete_volume", svc.deleteVolume},

		{"list_snapshots", svc.listSnapshots},
		{"get_snapshot", svc.getSnapshot},
		{"delete_snapshot", svc.deleteSnapshot},
		{"create_snapshot", svc.createSnapshot},

		{"actions", actions},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type volumeSvc struct {
	ctx context.Context
	svc volumes.Client
}

func (svc *volumeSvc) createVolume(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	req := godojs.ArgVolumeCreateRequest(vm, arg)
	d, err := svc.svc.CreateVolume(
		svc.ctx,
		req.Name, req.Region, req.SizeGigaBytes,
		volumes.SetVolumeDescription(req.Description),
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.VolumeToVM(vm, d.Struct())
}

func (svc *volumeSvc) getVolume(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := godojs.ArgVolumeID(vm, all.Argument(0))

	d, err := svc.svc.GetVolume(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.VolumeToVM(vm, d.Struct())
}

func (svc *volumeSvc) deleteVolume(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := godojs.ArgVolumeID(vm, all.Argument(0))

	err := svc.svc.DeleteVolume(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *volumeSvc) listVolume(all otto.FunctionCall) otto.Value {

	vm := all.Otto

	var volumes = make([]otto.Value, 0)
	volumec, errc := svc.svc.ListVolumes(svc.ctx)
	for d := range volumec {
		volumes = append(volumes, godojs.VolumeToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(volumes)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *volumeSvc) createSnapshot(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	req := godojs.ArgSnapshotCreateRequest(vm, arg)
	d, err := svc.svc.CreateSnapshot(svc.ctx, req.VolumeID, req.Name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.VolumeSnapshotToVM(vm, d.Struct())
}

func (svc *volumeSvc) getSnapshot(all otto.FunctionCall) otto.Value {
	var (
		vm = all.Otto
		id = godojs.ArgSnapshotID(vm, all.Argument(0))
	)
	d, err := svc.svc.GetSnapshot(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.VolumeSnapshotToVM(vm, d.Struct())
}

func (svc *volumeSvc) deleteSnapshot(all otto.FunctionCall) otto.Value {
	var (
		vm = all.Otto
		id = godojs.ArgSnapshotID(vm, all.Argument(0))
	)
	err := svc.svc.DeleteSnapshot(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *volumeSvc) listSnapshots(all otto.FunctionCall) otto.Value {

	var (
		vm       = all.Otto
		volumeID = godojs.ArgVolumeID(vm, all.Argument(0))
	)

	var Snapshots = make([]otto.Value, 0)
	snapshotc, errc := svc.svc.ListSnapshots(svc.ctx, volumeID)
	for d := range snapshotc {
		Snapshots = append(Snapshots, godojs.VolumeSnapshotToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(Snapshots)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
