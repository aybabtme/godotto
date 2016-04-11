package floatingips

import (
	"fmt"

	"github.com/aybabtme/godotto/internal/ottoutil"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client *godo.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := floatingIPSvc{
		svc: client.FloatingIPs,
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
	svc godo.FloatingIPsService
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
	fip, _, err := svc.svc.Create(req)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.floatingIPToVM(vm, *fip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	id := svc.argFloatingIP(all, 0)
	fip, _, err := svc.svc.Get(id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.floatingIPToVM(vm, *fip)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := svc.argFloatingIP(all, 0)

	_, err := svc.svc.Delete(id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *floatingIPSvc) list(all otto.FunctionCall) otto.Value {

	vm := all.Otto
	opt := &godo.ListOptions{Page: 1, PerPage: 200}

	var floatingIPs  = make([]otto.Value, 0)

	for {
		items, resp, err := svc.svc.List(opt)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}

		for _, d := range items {
			v, err := svc.floatingIPToVM(vm, d)
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			floatingIPs = append(floatingIPs, v)
		}

		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			break
		}
	}

	v, err := vm.ToValue(floatingIPs)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *floatingIPSvc) floatingIPToVM(vm *otto.Otto, g godo.FloatingIP) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
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
