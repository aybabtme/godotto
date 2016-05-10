package drives

import (
	"fmt"

	"golang.org/x/net/context"

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

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list_drives", svc.listDrive},
		{"get_drive", svc.getDrive},
		{"create_drive", svc.createDrive},
		{"delete_drive", svc.deleteDrive},

		{"list_snapshots", svc.listSnapshots},
		{"get_snapshot", svc.getSnapshot},
		{"delete_snapshot", svc.deleteSnapshot},
		{"create_snapshot", svc.createSnapshot},
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

func (svc *driveSvc) argDriveID(all otto.FunctionCall, i int) string {
	vm := all.Otto
	arg := all.Argument(i)

	var id string
	switch {
	case arg.IsString():
		id = ottoutil.String(vm, arg)
	case arg.IsObject():
		id = ottoutil.String(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Drive or a DriveID")
	}
	return id
}

func (svc *driveSvc) argSnapshotID(all otto.FunctionCall, i int) string {
	vm := all.Otto
	arg := all.Argument(i)

	var id string
	switch {
	case arg.IsString():
		id = ottoutil.String(vm, arg)
	case arg.IsObject():
		id = ottoutil.String(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Snapshot or a SnapshotID")
	}
	return id
}

func (svc *driveSvc) createDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a object")
	}

	var (
		name   = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name"))
		region = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "region"))
		size   = int64(ottoutil.Int(vm, ottoutil.GetObject(vm, arg, "size")))
		desc   = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "desc"))
		opt    []drives.CreateOpt
	)
	if desc != "" {
		opt = append(opt, drives.SetDriveDescription(desc))
	}
	d, err := svc.svc.CreateDrive(
		svc.ctx,
		name, region, size, opt...,
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.driveToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *driveSvc) getDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := svc.argDriveID(all, 0)

	d, err := svc.svc.GetDrive(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.driveToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *driveSvc) deleteDrive(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := svc.argDriveID(all, 0)

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
		v, err := svc.driveToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		drives = append(drives, v)
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
	snapArg := all.Argument(1).Object()
	var (
		driveID = svc.argDriveID(all, 0)
		name    = ottoutil.String(vm, ottoutil.GetObject(vm, snapArg, "name"))
		desc    = ottoutil.String(vm, ottoutil.GetObject(vm, snapArg, "desc"))
		opt     []drives.SnapshotOpt
	)
	if desc != "" {
		opt = append(opt, drives.SetSnapshotDescription(desc))
	}

	d, err := svc.svc.CreateSnapshot(svc.ctx, driveID, name, opt...)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.driveSnapshotToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *driveSvc) getSnapshot(all otto.FunctionCall) otto.Value {
	var (
		vm = all.Otto
		id = svc.argSnapshotID(all, 0)
	)
	d, err := svc.svc.GetSnapshot(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.driveSnapshotToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *driveSvc) deleteSnapshot(all otto.FunctionCall) otto.Value {
	var (
		vm = all.Otto
		id = svc.argSnapshotID(all, 0)
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
		driveID = svc.argDriveID(all, 0)
	)

	var Snapshots = make([]otto.Value, 0)
	snapshotc, errc := svc.svc.ListSnapshots(svc.ctx, driveID)
	for d := range snapshotc {
		v, err := svc.driveSnapshotToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		Snapshots = append(Snapshots, v)
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

func (svc *driveSvc) driveToVM(vm *otto.Otto, v drives.Drive) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", g.ID},
		{"name", g.Name},
		{"region", g.Region.Slug},
		{"size", g.SizeGigaBytes},
		{"description", g.Description},
	} {
		v, err := vm.ToValue(field.v)
		if err != nil {
			return q, fmt.Errorf("can't prepare field %q: %v", field.name, err)
		}
		if err := d.Set(field.name, v); err != nil {
			return q, fmt.Errorf("can't set field %q: %v", field.name, err)
		}
	}
	return d.Value(), nil
}

func (svc *driveSvc) driveSnapshotToVM(vm *otto.Otto, v drives.Snapshot) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", g.ID},
		{"drive_id", g.DriveID},
		{"name", g.Name},
		{"region", g.Region.Slug},
		{"size", g.SizeGibiBytes},
		{"description", g.Description},
	} {
		v, err := vm.ToValue(field.v)
		if err != nil {
			return q, fmt.Errorf("can't prepare field %q: %v", field.name, err)
		}
		if err := d.Set(field.name, v); err != nil {
			return q, fmt.Errorf("can't set field %q: %v", field.name, err)
		}
	}
	return d.Value(), nil
}
