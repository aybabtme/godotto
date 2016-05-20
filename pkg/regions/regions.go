package regions

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/regions"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := regionSvc{
		ctx: ctx,
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
	ctx context.Context
	svc regions.Client
}

func (svc *regionSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var regions = make([]otto.Value, 0)
	regionc, errc := svc.svc.List(svc.ctx)
	for d := range regionc {
		regions = append(regions, godojs.RegionToVM(vm, d.Struct()))
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
