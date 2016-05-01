package floatingips

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/aybabtme/godotto/internal/do/cloud/floatingips"
	"github.com/aybabtme/godotto/internal/ottoutil"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := floatingIPSvc{
		svc: client.FloatingIPs(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
		{"create", svc.create},
		{"get", svc.get},
		{"delete", svc.delete},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type floatingIPSvc struct {
	svc floatingips.Client
}

func (svc *floatingIPSvc) argCreateFloatingIP(all otto.FunctionCall, i int) *godo.FloatingIPCreateRequest {
	vm := all.Otto
	arg := all.Argument(i).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a Floating IP")
	}
	return &godo.FloatingIPCreateRequest{
		Region:    ottoutil.String(vm, ottoutil.GetObject(vm, arg, "region")),
		DropletID: ottoutil.Int(vm, ottoutil.GetObject(vm, arg, "droplet_id")),
	}
}

func (svc *floatingIPSvc) argFloatingIP(all otto.FunctionCall, i int) string {
	return ottoutil.String(all.Otto, all.Argument(i))
}

func (svc *floatingIPSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	req := svc.argCreateFloatingIP(all, 0)
	fip, err := svc.svc.Create(req.Region, floatingips.UseGodoFloatingIP(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.floatingIPToVM(vm, fip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	ip := svc.argFloatingIP(all, 0)
	fip, err := svc.svc.Get(ip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.floatingIPToVM(vm, fip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	ip := svc.argFloatingIP(all, 0)

	err := svc.svc.Delete(ip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *floatingIPSvc) list(all otto.FunctionCall) otto.Value {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vm := all.Otto

	var floatingIPs = make([]otto.Value, 0)
	floatingIPc, errc := svc.svc.List(ctx)
	for d := range floatingIPc {
		v, err := svc.floatingIPToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		floatingIPs = append(floatingIPs, v)
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(floatingIPs)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) floatingIPToVM(vm *otto.Otto, v floatingips.FloatingIP) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"region", g.Region},
		{"droplet", g.Droplet},
		{"ip", g.IP},
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
