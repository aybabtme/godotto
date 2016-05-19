package domains_test

import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

assert(pkg != null, "package should be loaded");
assert(pkg.get != null, "get function should be defined");
assert(pkg.list != null, "list function should be defined");
assert(pkg.create != null, "create function should be defined");
assert(pkg.delete != null, "delete function should be defined");

assert(pkg.record != null, "record function should be defined");
assert(pkg.records != null, "records function should be defined");
assert(pkg.create_record != null, "create_record function should be defined");
assert(pkg.edit_record != null, "edit_record function should be defined");
assert(pkg.delete_record != null, "delete_record function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockDomains.ListFn = func(_ context.Context) (<-chan domains.Domain, <-chan error) {
		lc := make(chan domains.Domain)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockDomains.GetFn = func(_ context.Context, _ string) (domains.Domain, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDomains.CreateFn = func(_ context.Context, _, _ string, _ ...domains.CreateOpt) (domains.Domain, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDomains.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}
	cloud.MockDomains.ListRecordFn = func(_ context.Context, _ string) (<-chan domains.Record, <-chan error) {
		lc := make(chan domains.Record)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockDomains.GetRecordFn = func(_ context.Context, _ string, _ int) (domains.Record, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDomains.CreateRecordFn = func(_ context.Context, _ string, _ ...domains.RecordOpt) (domains.Record, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDomains.UpdateRecordFn = func(_ context.Context, _ string, _ int, _ ...domains.RecordOpt) (domains.Record, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockDomains.DeleteRecordFn = func(_ context.Context, _ string, _ int) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var dr = {
	type: "",
	name: "",
	data: "",
	priority: "",
	port: "",
	weight: "",
};

[
	{ name: "list",          fn: function() { pkg.list() } },
	{ name: "get",           fn: function() { pkg.get("hello.com") } },
	{ name: "create",        fn: function() { pkg.create({}) } },
	{ name: "delete",        fn: function() { pkg.delete({}) } },

	{ name: "records",       fn: function() { pkg.records("hello.com") } },
	{ name: "record",        fn: function() { pkg.record("hello.com", 1) } },
	{ name: "create_record", fn: function() { pkg.create_record("hello.com", {}) } },
	{ name: "edit_record",   fn: function() { pkg.edit_record("hello.com", dr) } },
	{ name: "delete_record", fn: function() { pkg.delete_record("hello.com", 1)  } }

].forEach(function(kv) {
	var name = kv.name;
	var fn = kv.fn;
	try {
		fn(); throw "dont catch me";
	} catch (e) {
		equals("throw me", e.message, name +" should send the right exception");
	};
})`)
}

type domain struct {
	*godo.Domain
}

func (k *domain) Struct() *godo.Domain { return k.Domain }

type record struct {
	*godo.DomainRecord
}

func (k *record) Struct() *godo.DomainRecord { return k.DomainRecord }

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.ListFn = func(_ context.Context) (<-chan domains.Domain, <-chan error) {
		lc := make(chan domains.Domain, 1)
		lc <- &domain{&godo.Domain{Name: "my_name", TTL: 42, ZoneFile: "my_zone_file"}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var d = list[0];
var want = {
	name: "my_name",
	ttl: 42,
	zone_file: "my_zone_file"
};
equals(d, want, "should have proper object");
`)
}

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.GetFn = func(_ context.Context, _ string) (domains.Domain, error) {
		return &domain{&godo.Domain{Name: "my_name", TTL: 42, ZoneFile: "my_zone_file"}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var d = pkg.get("my_name")
var want = {
	name: "my_name",
	ttl: 42,
	zone_file: "my_zone_file"
};
equals(d, want, "should have proper object");
`)
}

func TestCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.CreateFn = func(_ context.Context, name, ip string, opt ...domains.CreateOpt) (domains.Domain, error) {
		return &domain{&godo.Domain{Name: name, TTL: 42, ZoneFile: "my_zone_file"}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var d = pkg.create({
	name: "my_name",
	ip: "127.0.0.1"
});
var want = {
	name: "my_name",
	ttl: 42,
	zone_file: "my_zone_file"
};
equals(d, want, "should have proper object");
`)
}

func TestDelete(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.DeleteFn = func(_ context.Context, gotName string) error {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

pkg.delete("my name");
`)
}

func TestListRecord(t *testing.T) {
	wantName := "my name"
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.ListRecordFn = func(_ context.Context, gotName string) (<-chan domains.Record, <-chan error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		lc := make(chan domains.Record, 1)
		lc <- &record{&godo.DomainRecord{
			ID: 42, Type: "srv", Name: wantName, Data: "derp", Priority: 43, Port: 8080, Weight: 9000,
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var records = pkg.records("my name");
assert(records != null, "should have received a records");
assert(records.length > 0, "should have received some elements")

var d = records[0];
var want = {
	id: 42,
	type: "srv",
	name: "my name",
	data: "derp",
	priority: 43,
	port: 8080,
	weight: 9000
};
equals(d, want, "should have proper object");
`)
}

func TestGetRecord(t *testing.T) {
	wantName := "my name"
	wantID := 42

	cloud := mockcloud.Client(nil)
	cloud.MockDomains.GetRecordFn = func(_ context.Context, gotName string, gotID int) (domains.Record, error) {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return &record{&godo.DomainRecord{
			ID: 42, Type: "srv", Name: wantName, Data: "derp", Priority: 43, Port: 8080, Weight: 9000,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var d = pkg.record("my name", 42)
var want = {
	id: 42,
	type: "srv",
	name: "my name",
	data: "derp",
	priority: 43,
	port: 8080,
	weight: 9000
};
equals(d, want, "should have proper object");
`)
}

func TestCreateRecord(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.CreateRecordFn = func(_ context.Context, name string, _ ...domains.RecordOpt) (domains.Record, error) {
		return &record{&godo.DomainRecord{
			ID: 42, Type: "srv", Name: name, Data: "derp", Priority: 43, Port: 8080, Weight: 9000,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var want = {
	id: 42,
	type: "srv",
	name: "my name",
	data: "derp",
	priority: 43,
	port: 8080,
	weight: 9000
};

var d = pkg.create_record(want.name, want);
equals(d, want, "should have proper object");
`)
}

func TestUpdateRecord(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockDomains.UpdateRecordFn = func(_ context.Context, name string, id int, _ ...domains.RecordOpt) (domains.Record, error) {
		return &record{&godo.DomainRecord{
			ID: id, Type: "srv", Name: name, Data: "derp", Priority: 43, Port: 8080, Weight: 9000,
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

var want = {
	id: 42,
	type: "srv",
	name: "my name",
	data: "derp",
	priority: 43,
	port: 8080,
	weight: 9000
};

var d = pkg.edit_record(want.name, want);
equals(d, want, "should have proper object");
`)
}

func TestDeleteRecord(t *testing.T) {

	wantName := "my name"
	wantID := 42

	cloud := mockcloud.Client(nil)
	cloud.MockDomains.DeleteRecordFn = func(_ context.Context, gotName string, gotID int) error {
		if gotName != wantName {
			t.Fatalf("want %q got %q", wantName, gotName)
		}
		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.domains;

pkg.delete_record("my name", 42);
`)
}
