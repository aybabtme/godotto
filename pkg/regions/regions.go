package regions

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/aybabtme/godotto/internal/do/cloud/regions"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := regionSvc{
		svc: client.Regions(),
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
	svc regions.Client
}

func (svc *regionSvc) list(all otto.FunctionCall) otto.Value {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vm := all.Otto

	var regions = make([]otto.Value, 0)
	regionc, errc := svc.svc.List(ctx)
	for d := range regionc {
		v, err := svc.regionToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		regions = append(regions, v)
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(regions)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *regionSvc) regionToVM(vm *otto.Otto, v regions.Region) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
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
