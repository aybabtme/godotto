package images

import (
	"fmt"

	"golang.org/x/net/context"

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

func (svc *imageSvc) argImageID(all otto.FunctionCall, i int) int {
	vm := all.Otto
	arg := all.Argument(i)

	var id int
	switch {
	case arg.IsNumber():
		id = ottoutil.Int(vm, arg)
	case arg.IsObject():
		id = ottoutil.Int(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be a Image or a ImageID")
	}
	return id
}

func (svc *imageSvc) argImageSlug(all otto.FunctionCall, i int) string {
	vm := all.Otto
	arg := all.Argument(i)

	var slug string
	switch {
	case arg.IsString():
		slug = ottoutil.String(vm, arg)
	case arg.IsObject():
		slug = ottoutil.String(vm, ottoutil.GetObject(vm, arg.Object(), "slug"))
	default:
		ottoutil.Throw(vm, "argument must be a Image or a ImageSlug")
	}
	return slug
}

func (svc *imageSvc) argImageUpdate(all otto.FunctionCall, i int) *godo.ImageUpdateRequest {
	vm := all.Otto
	arg := all.Argument(i).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a ImageRecord")
	}
	return &godo.ImageUpdateRequest{
		Name: ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
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
		id := svc.argImageID(all, 0)
		img, err = svc.svc.GetByID(svc.ctx, id)
	case arg.IsString():
		slug := svc.argImageSlug(all, 0)
		img, err = svc.svc.GetBySlug(svc.ctx, slug)
	}
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.imageToVM(vm, img)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *imageSvc) update(all otto.FunctionCall) otto.Value {
	vm := all.Otto

	var (
		// they read the same arg, just different fields
		id  = svc.argImageID(all, 0)
		req = svc.argImageUpdate(all, 1)
	)
	img, err := svc.svc.Update(svc.ctx, id, images.UseGodoImage(req))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.imageToVM(vm, img)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *imageSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	id := svc.argImageID(all, 0)

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
		v, err := svc.imageToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		images = append(images, v)
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

func (svc *imageSvc) imageToVM(vm *otto.Otto, v images.Image) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", int64(g.ID)},
		{"name", g.Name},
		{"type", g.Type},
		{"distribution", g.Distribution},
		{"slug", g.Slug},
		{"public", g.Public},
		{"regions", g.Regions},
		{"min_disk_size", int64(g.MinDiskSize)},
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
