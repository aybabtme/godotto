package keys

import (
	"fmt"
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godojs"
	"github.com/aybabtme/godotto/pkg/extra/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := keySvc{
		ctx: ctx,
		svc: client.Keys(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
		{"create", svc.create},
		{"get", svc.get},
		{"update", svc.update},
		{"delete", svc.delete},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type keySvc struct {
	ctx context.Context
	svc keys.Client
}

func (svc *keySvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	req := godojs.ArgKeyCreate(vm, all.Argument(0))
	key, err := svc.svc.Create(svc.ctx, req.Name, req.PublicKey)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.KeyToVM(vm, key.Struct())
}

func (svc *keySvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var (
		key keys.Key
		err error
	)
	arg := all.Argument(0)
	switch {
	case arg.IsNumber():
		id := godojs.ArgKeyID(vm, all.Argument(0))
		key, err = svc.svc.GetByID(svc.ctx, id)
	case arg.IsString():
		fp := godojs.ArgKeyFingerprint(vm, all.Argument(0))
		key, err = svc.svc.GetByFingerprint(svc.ctx, fp)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.KeyToVM(vm, key.Struct())
}

func (svc *keySvc) update(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	var (
		key keys.Key
		err error
	)
	switch {
	case arg.IsNumber():
		id := godojs.ArgKeyID(vm, all.Argument(0))
		req := godojs.ArgKeyUpdate(vm, all.Argument(1))
		key, err = svc.svc.UpdateByID(svc.ctx, id, keys.UseGodoKey(req))
	case arg.IsString():
		fp := godojs.ArgKeyFingerprint(vm, all.Argument(0))
		req := godojs.ArgKeyUpdate(vm, all.Argument(1))
		key, err = svc.svc.UpdateByFingerprint(svc.ctx, fp, keys.UseGodoKey(req))
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.KeyToVM(vm, key.Struct())
}

func (svc *keySvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	var err error
	switch {
	case arg.IsNumber():
		id := godojs.ArgKeyID(vm, all.Argument(0))
		err = svc.svc.DeleteByID(svc.ctx, id)
	case arg.IsString():
		fp := godojs.ArgKeyFingerprint(vm, all.Argument(0))
		err = svc.svc.DeleteByFingerprint(svc.ctx, fp)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *keySvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var keys = make([]otto.Value, 0)
	keyc, errc := svc.svc.List(svc.ctx)
	for d := range keyc {
		keys = append(keys, godojs.KeyToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(keys)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
