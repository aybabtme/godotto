package sizes

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/aybabtme/godotto/internal/do/cloud/sizes"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := sizeSvc{
		svc: client.Sizes(),
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
	svc sizes.Client
}

func (svc *sizeSvc) list(all otto.FunctionCall) otto.Value {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vm := all.Otto

	var sizes = make([]otto.Value, 0)
	sizec, errc := svc.svc.List(ctx)
	for d := range sizec {
		v, err := svc.sizeToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		sizes = append(sizes, v)
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(sizes)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *sizeSvc) sizeToVM(vm *otto.Otto, v sizes.Size) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
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
