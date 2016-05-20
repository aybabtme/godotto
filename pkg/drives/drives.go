package drives

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/drives"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := driveSvc{
		ctx: ctx,
		svc: client.Drives(),
	}

	actions, err := applyAction(ctx, vm, client)
	if err != nil {
		return q, err
	}

	for _, applier := range []struct {
		Name   string
		Method interface{}
	}{
		{"list_drives", svc.listDrive},
		{"get_drive", svc.getDrive},
		{"create_drive", svc.createDrive},
		{"delete_drive", svc.deleteDrive},

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

type driveSvc struct {
	ctx context.Context
	svc drives.Client
}

func (svc *driveSvc) createDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	req := godojs.ArgDriveCreateRequest(vm, arg)
	d, err := svc.svc.CreateDrive(
		svc.ctx,
		req.Name, req.Region, req.SizeGibiBytes,
		drives.SetDriveDescription(req.Description),
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DriveToVM(vm, d.Struct())
}

func (svc *driveSvc) getDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := godojs.ArgDriveID(vm, all.Argument(0))

	d, err := svc.svc.GetDrive(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.DriveToVM(vm, d.Struct())
}

func (svc *driveSvc) deleteDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := godojs.ArgDriveID(vm, all.Argument(0))

	err := svc.svc.DeleteDrive(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *driveSvc) listDrive(all otto.FunctionCall) otto.Value {

	vm := all.Otto

	var drives = make([]otto.Value, 0)
	drivec, errc := svc.svc.ListDrives(svc.ctx)
	for d := range drivec {
		drives = append(drives, godojs.DriveToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(drives)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *driveSvc) createSnapshot(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	req := godojs.ArgSnapshotCreateRequest(vm, arg)

	d, err := svc.svc.CreateSnapshot(svc.ctx, req.DriveID, req.Name,
		drives.SetSnapshotDescription(req.Description),
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DriveSnapshotToVM(vm, d.Struct())
}

func (svc *driveSvc) getSnapshot(all otto.FunctionCall) otto.Value {
	var (
		vm = all.Otto
		id = godojs.ArgSnapshotID(vm, all.Argument(0))
	)
	d, err := svc.svc.GetSnapshot(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DriveSnapshotToVM(vm, d.Struct())
}

func (svc *driveSvc) deleteSnapshot(all otto.FunctionCall) otto.Value {
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

func (svc *driveSvc) listSnapshots(all otto.FunctionCall) otto.Value {

	var (
		vm      = all.Otto
		driveID = godojs.ArgDriveID(vm, all.Argument(0))
	)

	var Snapshots = make([]otto.Value, 0)
	snapshotc, errc := svc.svc.ListSnapshots(svc.ctx, driveID)
	for d := range snapshotc {
		Snapshots = append(Snapshots, godojs.DriveSnapshotToVM(vm, d.Struct()))
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
