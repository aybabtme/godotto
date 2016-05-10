package domains

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"

	"github.com/digitalocean/godo"
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
		{"delete_record", svc.deleteRecord},
		{"edit_record", svc.editRecord},
		{"create_record", svc.createRecord},
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

func (svc *domainSvc) argDomainName(all otto.FunctionCall, i int) string {
	vm := all.Otto
	arg := all.Argument(i)

	var name string
	switch {
	case arg.IsString():
		name = ottoutil.String(vm, arg)
	case arg.IsObject():
		name = ottoutil.String(vm, ottoutil.GetObject(vm, arg.Object(), "name"))
	default:
		ottoutil.Throw(vm, "argument must be a Domain or a DomainName")
	}
	return name
}

func (svc *domainSvc) argRecordID(all otto.FunctionCall, i int) int {
	return ottoutil.Int(all.Otto, all.Argument(i))
}

func (svc *domainSvc) argDomainRecord(all otto.FunctionCall, i int) *godo.DomainRecordEditRequest {
	vm := all.Otto
	arg := all.Argument(i).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a DomainRecord")
	}
	return &godo.DomainRecordEditRequest{
		Type:     ottoutil.String(vm, ottoutil.GetObject(vm, arg, "type")),
		Name:     ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
		Data:     ottoutil.String(vm, ottoutil.GetObject(vm, arg, "data")),
		Priority: ottoutil.Int(vm, ottoutil.GetObject(vm, arg, "priority")),
		Port:     ottoutil.Int(vm, ottoutil.GetObject(vm, arg, "port")),
		Weight:   ottoutil.Int(vm, ottoutil.GetObject(vm, arg, "weight")),
	}
}

func (svc *domainSvc) create(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	arg := all.Argument(0).Object()
	if arg == nil {
		ottoutil.Throw(vm, "argument must be a object")
	}

	d, err := svc.svc.Create(
		svc.ctx,
		ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
		ottoutil.String(vm, ottoutil.GetObject(vm, arg, "ip_address")),
	)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	name := svc.argDomainName(all, 0)

	d, err := svc.svc.Get(svc.ctx, name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.domainToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	name := svc.argDomainName(all, 0)

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
		v, err := svc.domainToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		domains = append(domains, v)
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
	vm := all.Otto
	name := svc.argDomainName(all, 0)
	record := svc.argDomainRecord(all, 1)

	d, err := svc.svc.CreateRecord(svc.ctx, name, domains.UseGodoRecord(record))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) record(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = svc.argDomainName(all, 0)
		id   = svc.argRecordID(all, 1)
	)
	d, err := svc.svc.GetRecord(svc.ctx, name, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) editRecord(all otto.FunctionCall) otto.Value {
	var (
		vm     = all.Otto
		name   = svc.argDomainName(all, 0)
		id     = svc.argRecordID(all, 1)
		record = svc.argDomainRecord(all, 2)
	)
	d, err := svc.svc.UpdateRecord(svc.ctx, name, id, domains.UseGodoRecord(record))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) deleteRecord(all otto.FunctionCall) otto.Value {
	var (
		vm   = all.Otto
		name = svc.argDomainName(all, 0)
		id   = svc.argRecordID(all, 1)
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
		name = svc.argDomainName(all, 0)
	)

	var records = make([]otto.Value, 0)
	recordc, errc := svc.svc.ListRecord(svc.ctx, name)
	for d := range recordc {
		v, err := svc.domainRecordToVM(vm, d)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		records = append(records, v)
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

func (svc *domainSvc) domainToVM(vm *otto.Otto, v domains.Domain) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"name", g.Name},
		{"ttl", g.TTL},
		{"zone_file", g.ZoneFile},
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

func (svc *domainSvc) domainRecordToVM(vm *otto.Otto, v domains.Record) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
	g := v.Struct()
	for _, field := range []struct {
		name string
		v    interface{}
	}{
		{"id", g.ID},
		{"type", g.Type},
		{"name", g.Name},
		{"data", g.Data},
		{"priority", g.Priority},
		{"port", g.Port},
		{"weight", g.Weight},
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
