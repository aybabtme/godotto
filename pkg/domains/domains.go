package domains

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/internal/godojs"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"

	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := domainSvc{
		ctx: ctx,
		svc: client.Domains(),
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"list", svc.list},
		{"get", svc.get},
		{"create", svc.create},
		{"delete", svc.delete},

		{"records", svc.records},
		{"record", svc.record},
		{"create_record", svc.createRecord},
		{"edit_record", svc.editRecord},
		{"delete_record", svc.deleteRecord},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type domainSvc struct {
	ctx context.Context
	svc domains.Client
}

func (svc *domainSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0)

	req := godojs.ArgDomainCreateRequest(vm, arg)
	d, err := svc.svc.Create(
		svc.ctx,
		req.Name, req.IPAddress,
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DomainToVM(vm, d.Struct())
}

func (svc *domainSvc) get(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = godojs.ArgDomainName(vm, all.Argument(0))
	)
	d, err := svc.svc.Get(svc.ctx, name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.DomainToVM(vm, d.Struct())
}

func (svc *domainSvc) delete(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = godojs.ArgDomainName(vm, all.Argument(0))
	)
	err := svc.svc.Delete(svc.ctx, name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *domainSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	var domains = make([]otto.Value, 0)
	domainc, errc := svc.svc.List(svc.ctx)
	for d := range domainc {
		domains = append(domains, godojs.DomainToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(domains)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) createRecord(all otto.FunctionCall) otto.Value {
	var (
		vm     = all.Otto
		name   = godojs.ArgDomainName(vm, all.Argument(0))
		record = godojs.ArgDomainRecord(vm, all.Argument(1))
	)
	d, err := svc.svc.CreateRecord(svc.ctx, name, domains.UseGodoRecord(record))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DomainRecordToVM(vm, d.Struct())
}

func (svc *domainSvc) record(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = godojs.ArgDomainName(vm, all.Argument(0))
		id   = godojs.ArgRecordID(vm, all.Argument(1))
	)
	d, err := svc.svc.GetRecord(svc.ctx, name, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return godojs.DomainRecordToVM(vm, d.Struct())
}

func (svc *domainSvc) editRecord(all otto.FunctionCall) otto.Value {
	var (
		vm     = all.Otto
		name   = godojs.ArgDomainName(vm, all.Argument(0))
		id     = godojs.ArgRecordID(vm, all.Argument(1))
		record = godojs.ArgDomainRecord(vm, all.Argument(1))
	)
	d, err := svc.svc.UpdateRecord(svc.ctx, name, id, domains.UseGodoRecord(record))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return godojs.DomainRecordToVM(vm, d.Struct())
}

func (svc *domainSvc) deleteRecord(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = godojs.ArgDomainName(vm, all.Argument(0))
		id   = godojs.ArgRecordID(vm, all.Argument(1))
	)
	err := svc.svc.DeleteRecord(svc.ctx, name, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *domainSvc) records(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = godojs.ArgDomainName(vm, all.Argument(0))
	)
	var records = make([]otto.Value, 0)
	recordc, errc := svc.svc.ListRecord(svc.ctx, name)
	for d := range recordc {
		records = append(records, godojs.DomainRecordToVM(vm, d.Struct()))
	}
	if err := <-errc; err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := vm.ToValue(records)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}
