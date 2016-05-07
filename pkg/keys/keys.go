package keys

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := keySvc{
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
	svc keys.Client
}

func (svc *keySvc) argKeyID(all otto.FunctionCall, i int) int {
	vm := all.Otto
	arg := all.Argument(i)

	var id int
	switch {
	case arg.IsNumber():
		id = ottoutil.Int(vm, arg)
	case arg.IsObject():
		id = ottoutil.Int(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Key or a KeyID")
	}
	return id
}

func (svc *keySvc) argKeyFingerprint(all otto.FunctionCall, i int) string {
	vm := all.Otto
	arg := all.Argument(i)

	var fp string
	switch {
	case arg.IsString():
		fp = ottoutil.String(vm, arg)
	case arg.IsObject():
		fp = ottoutil.String(vm, ottoutil.GetObject(vm, arg.Object(), "fp"))
	default:
		ottoutil.Throw(vm, "argument must be a Key or a KeyFingerprint")
	}
	return fp
}

func (svc *keySvc) argKeyCreate(all otto.FunctionCall, i int) *godo.KeyCreateRequest {
	vm := all.Otto
	arg := all.Argument(i).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a Key")
	}
	return &godo.KeyCreateRequest{
		Name:      ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
		PublicKey: ottoutil.String(vm, ottoutil.GetObject(vm, arg, "public_key")),
	}
}

func (svc *keySvc) argKeyUpdate(all otto.FunctionCall, i int) *godo.KeyUpdateRequest {
	vm := all.Otto
	arg := all.Argument(i).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a Key")
	}
	return &godo.KeyUpdateRequest{
		Name: ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
	}
}

func (svc *keySvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	req := svc.argKeyCreate(all, 0)
	key, err := svc.svc.Create(req.Name, req.PublicKey)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.keyToVM(vm, key)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
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
		id := svc.argKeyID(all, 0)
		key, err = svc.svc.GetByID(id)
	case arg.IsString():
		fp := svc.argKeyFingerprint(all, 0)
		key, err = svc.svc.GetByFingerprint(fp)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.keyToVM(vm, key)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
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
		id := svc.argKeyID(all, 0)
		req := svc.argKeyUpdate(all, 1)
		key, err = svc.svc.UpdateByID(id, keys.UseGodoKey(req))
	case arg.IsString():
		fp := svc.argKeyFingerprint(all, 0)
		req := svc.argKeyUpdate(all, 1)
		key, err = svc.svc.UpdateByFingerprint(fp, keys.UseGodoKey(req))
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.keyToVM(vm, key)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *keySvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)
	var err error
	switch {
	case arg.IsNumber():
		id := svc.argKeyID(all, 0)
		err = svc.svc.DeleteByID(id)
	case arg.IsString():
		fp := svc.argKeyFingerprint(all, 0)
		err = svc.svc.DeleteByFingerprint(fp)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *keySvc) list(all otto.FunctionCall) otto.Value {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vm := all.Otto

	var keys = make([]otto.Value, 0)
	keyc, errc := svc.svc.List(ctx)
	for d := range keyc {
		v, err := svc.keyToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		keys = append(keys, v)
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

func (svc *keySvc) keyToVM(vm *otto.Otto, v keys.Key) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", g.ID},
		{"name", g.Name},
		{"fingerprint", g.Fingerprint},
		{"public_key", g.PublicKey},
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
