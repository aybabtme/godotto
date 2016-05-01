package actions

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/aybabtme/godotto/internal/do/cloud/actions"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := actionSvc{
		svc: client.Actions(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"get", svc.get},
		{"list", svc.list},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type actionSvc struct {
	svc actions.Client
}

func (svc *actionSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	var aid int
	switch {
	case arg.IsNumber():
		aid = ottoutil.Int(vm, arg)
	case arg.IsObject():
		aid = ottoutil.Int(vm, ottoutil.GetObject(vm, arg.Object(), "id"))
	default:
		ottoutil.Throw(vm, "argument must be an Action or an ActionID")
	}

	a, err := svc.svc.Get(aid)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.actionToVM(vm, a)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *actionSvc) list(all otto.FunctionCall) otto.Value {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vm := all.Otto
	actions := []actions.Action{}
	actionc, errc := svc.svc.List(ctx)
	for action := range actionc {
		actions = append(actions, action)
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(actions)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *actionSvc) actionToVM(vm *otto.Otto, a actions.Action) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := a.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", g.ID},
		{"status", g.Status},
		{"type", g.Type},
		{"started_at", g.StartedAt},
		{"completed_at", g.CompletedAt},
		{"resource_id", g.ResourceID},
		{"resource_type", g.ResourceType},
		{"region", g.Region},
		{"region_slug", g.RegionSlug},
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
