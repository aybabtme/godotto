package domains

import (
	"fmt"

	"github.com/aybabtme/godotto/internal/ottoutil"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto, client *godo.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := domainSvc{
		svc: client.Domains,
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
	svc godo.DomainsService
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

	opts := &godo.DomainCreateRequest{
		Name:      ottoutil.String(vm, ottoutil.GetObject(vm, arg, "name")),
		IPAddress: ottoutil.String(vm, ottoutil.GetObject(vm, arg, "ip_address")),
	}

	d, _, err := svc.svc.Create(opts)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainToVM(vm, *d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) get(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	name := svc.argDomainName(all, 0)

	d, _, err := svc.svc.Get(name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	v, err := svc.domainToVM(vm, *d)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) delete(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	name := svc.argDomainName(all, 0)

	_, err := svc.svc.Delete(name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *domainSvc) list(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	opt := &godo.ListOptions{Page: 1, PerPage: 200}

	var domains  = make([]otto.Value, 0)

	for {
		items, resp, err := svc.svc.List(opt)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}

		for _, d := range items {
			v, err := svc.domainToVM(vm, d)
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			domains = append(domains, v)
		}

		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			break
		}
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

	d, _, err := svc.svc.CreateRecord(name, record)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, *d)
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
	d, _, err := svc.svc.Record(name, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, *d)
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
	d, _, err := svc.svc.EditRecord(name, id, record)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	v, err := svc.domainRecordToVM(vm, *d)
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
	_, err := svc.svc.DeleteRecord(name, id)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *domainSvc) records(all otto.FunctionCall) otto.Value {

	vm := all.Otto
	name := svc.argDomainName(all, 0)

	opt := &godo.ListOptions{Page: 1, PerPage: 200}

	var records  = make([]otto.Value, 0)

	for {
		items, resp, err := svc.svc.Records(name, opt)
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}

		for _, d := range items {
			v, err := svc.domainRecordToVM(vm, d)
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			records = append(records, v)
		}

		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			break
		}
	}

	v, err := vm.ToValue(records)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return v
}

func (svc *domainSvc) domainToVM(vm *otto.Otto, g godo.Domain) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
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

func (svc *domainSvc) domainRecordToVM(vm *otto.Otto, g godo.DomainRecord) (otto.Value, error) {
	d, _ := vm.Object(`({})`)
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
