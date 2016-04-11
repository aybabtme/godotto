package regions

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

	svc := regionSvc{
		svc: client.Regions,
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

type regionSvc struct {
	svc godo.RegionsService
}

func (svc *regionSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	opt := &godo.ListOptions{Page: 1, PerPage: 200}

	var regions  = make([]otto.Value, 0)

	for {
		items, resp, err := svc.svc.List(opt)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}

		for _, a := range items {
			v, err := svc.regionToVM(vm, a)
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			regions = append(regions, v)
		}

		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			break
		}
	}

	v, err := vm.ToValue(regions)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *regionSvc) regionToVM(vm *otto.Otto, g godo.Region) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"slug", g.Slug},
		{"name", g.Name},
		{"sizes", g.Sizes},
		{"available", g.Available},
		{"features", g.Features},
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
