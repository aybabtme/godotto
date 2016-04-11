package sizes

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

	svc := sizeSvc{
		svc: client.Sizes,
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type sizeSvc struct {
	svc godo.SizesService
}

func (svc *sizeSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	opt := &godo.ListOptions{Page: 1, PerPage: 200}

	var sizes  = make([]otto.Value, 0)

	for {
		items, resp, err := svc.svc.List(opt)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}

		for _, a := range items {
			v, err := svc.sizeToVM(vm, a)
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			sizes = append(sizes, v)
		}

		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			break
		}
	}

	v, err := vm.ToValue(sizes)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *sizeSvc) sizeToVM(vm *otto.Otto, g godo.Size) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"slug", g.Slug},
		{"memory", g.Memory},
		{"vcpus", g.Vcpus},
		{"disk", g.Disk},
		{"price_monthly", g.PriceMonthly},
		{"price_hourly", g.PriceHourly},
		{"regions", g.Regions},
		{"available", g.Available},
		{"transfer", g.Transfer},
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
