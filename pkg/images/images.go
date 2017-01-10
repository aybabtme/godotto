package images

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/images"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := imageSvc{
		ctx: ctx,
		svc: client.Images(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
		{"list_distribution", svc.listDistribution},
		{"list_application", svc.listApplication},
		{"list_user", svc.listUser},
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

type imageSvc struct {
	ctx context.Context
	svc images.Client
}

func (svc *imageSvc) argImageUpdate(all otto.FunctionCall, i int) *godo.ImageUpdateRequest {
	vm := all.Otto
	arg := all.Argument(i)

	return &godo.ImageUpdateRequest{
		Name: godojs.ArgImageName(vm, arg),
	}
}

func (svc *imageSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var (
		img images.Image
		err error
	)
	arg := all.Argument(0)
	switch {
	case arg.IsNumber():
		id := godojs.ArgImageID(vm, all.Argument(0))
		img, err = svc.svc.GetByID(svc.ctx, id)
	case arg.IsString():
		slug := godojs.ArgImageSlug(vm, all.Argument(0))
		img, err = svc.svc.GetBySlug(svc.ctx, slug)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.ImageToVM(vm, img.Struct())
}

func (svc *imageSvc) update(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var (
		// they read the same arg, just different fields
		id  = godojs.ArgImageID(vm, all.Argument(0))
		req = svc.argImageUpdate(all, 1)
	)
	img, err := svc.svc.Update(svc.ctx, id, images.UseGodoImage(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.ImageToVM(vm, img.Struct())
}

func (svc *imageSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := godojs.ArgImageID(vm, all.Argument(0))

	err := svc.svc.Delete(svc.ctx, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *imageSvc) list(all otto.FunctionCall) otto.Value {
	return svc.listCommon(all, svc.svc.List)
}

func (svc *imageSvc) listDistribution(all otto.FunctionCall) otto.Value {
	return svc.listCommon(all, svc.svc.ListDistribution)
}

func (svc *imageSvc) listApplication(all otto.FunctionCall) otto.Value {
	return svc.listCommon(all, svc.svc.ListApplication)
}

func (svc *imageSvc) listUser(all otto.FunctionCall) otto.Value {
	return svc.listCommon(all, svc.svc.ListUser)
}

type listfunc func(context.Context) (<-chan images.Image, <-chan error)

func (svc *imageSvc) listCommon(all otto.FunctionCall, listfn listfunc) otto.Value {
	vm := all.Otto

	var images = make([]otto.Value, 0)
	imagec, errc := listfn(svc.ctx)
	for d := range imagec {
		images = append(images, godojs.ImageToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(images)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
