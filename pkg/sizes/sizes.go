package sizes

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/sizes"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := sizeSvc{
		ctx: ctx,
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
	ctx context.Context
	svc sizes.Client
}

func (svc *sizeSvc) list(all otto.FunctionCall) otto.Value {

	vm := all.Otto

	var sizes = make([]otto.Value, 0)
	sizec, errc := svc.svc.List(svc.ctx)
	for d := range sizec {
		sizes = append(sizes, godojs.SizeToVM(vm, d.Struct()))
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
