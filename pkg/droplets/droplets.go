package droplets

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"

	"github.com/digitalocean/godo"
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

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
		{"get", svc.get},
		{"create", svc.create},
		{"delete", svc.delete},
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
	arg := all.Argument(0).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a object")
	}

	imgArg := ottoutil.GetObject(vm, arg, "image")
	if imgArg.IsUndefined() {
		ottoutil.Throw(vm, "object must contain an 'image' field")
	}

	var (
		name   = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name"))
		region = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "region"))
		size   = ottoutil.String(vm, ottoutil.GetObject(vm, arg, "size"))
		image  string
	)

	switch {
	case imgArg.IsString():
		image = ottoutil.String(vm, imgArg)
	case imgArg.IsObject():
		image = ottoutil.String(vm, ottoutil.GetObject(vm, imgArg.Object(), "slug"))
	}

	opts := &godo.DropletCreateRequest{
		Backups:           ottoutil.Bool(vm, ottoutil.GetObject(vm, arg, "backups")),
		IPv6:              ottoutil.Bool(vm, ottoutil.GetObject(vm, arg, "ipv6")),
		PrivateNetworking: ottoutil.Bool(vm, ottoutil.GetObject(vm, arg, "private_networking")),
		UserData:          ottoutil.String(vm, ottoutil.GetObject(vm, arg, "user_data")),
	}

	sshArgs := ottoutil.GetObject(vm, arg, "ssh_keys").Object()
	if sshArgs != nil {
		for _, k := range sshArgs.Keys() {
			sshArg := ottoutil.GetObject(vm, sshArgs, k).Object()
			if sshArg == nil {
				ottoutil.Throw(vm, "'ssh_keys' field must be an object")
			}
			opts.SSHKeys = append(opts.SSHKeys, godo.DropletCreateSSHKey{
				ID:          int(ottoutil.Int(vm, ottoutil.GetObject(vm, sshArg, "id"))),
				Fingerprint: ottoutil.String(vm, ottoutil.GetObject(vm, sshArg, "fingerprint")),
			})
		}
	}

	d, err := svc.svc.Create(svc.ctx, name, region, size, image, droplets.UseGodoCreate(opts))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.dropletToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *dropletSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	var did int
	switch {
	case arg.IsNumber():
		did = ottoutil.Int(vm, arg)
	case arg.IsObject():
		did = ottoutil.Int(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Droplet or a DropletID")
	}

	d, err := svc.svc.Get(svc.ctx, did)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.dropletToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *dropletSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	var did int
	switch {
	case arg.IsNumber():
		did = ottoutil.Int(vm, arg)
	case arg.IsObject():
		did = ottoutil.Int(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Droplet or a DropletID")
	}

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
		v, err := svc.dropletToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		droplets = append(droplets, v)
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

// helpers

func (svc *dropletSvc) dropletToVM(vm *otto.Otto, v droplets.Droplet) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	publicIPv4, _ := g.PublicIPv4()

	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", int64(g.ID)},
		{"name", g.Name},
		{"memory", int64(g.Memory)},
		{"vcpus", int64(g.Vcpus)},
		{"disk", int64(g.Disk)},
		{"region_slug", g.Region.Slug},
		{"image_id", int64(g.Image.ID)},
		{"image_slug", g.Image.Slug},
		{"size_slug", g.Size.Slug},
		{"backup_ids", intsToInt64s(g.BackupIDs)},
		{"snapshot_ids", intsToInt64s(g.SnapshotIDs)},
		{"locked", g.Locked},
		{"status", g.Status},
		{"public_ipv4", publicIPv4},
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

func intsToInt64s(in []int) (out []int64) {
	for _, i := range in {
		out = append(out, int64(i))
	}
	return out
}
