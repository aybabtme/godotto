package sizes_test

import (
	"context"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/sizes"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.sizes;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

type size struct {
	*godo.Size
}

func (k *size) Struct() *godo.Size { return k.Size }

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)

	want := &godo.Size{
		Slug: "lol", Memory: 1, Vcpus: 2, Disk: 2, PriceMonthly: 1.0, PriceHourly: 0.1,
		Regions: []string{"lol"}, Available: true, Transfer: 1.0,
	}

	cloud.MockSizes.ListFn = func(_ context.Context) (<-chan sizes.Size, <-chan error) {
		rc, ec := make(chan sizes.Size, 1), make(chan error, 0)
		rc <- &size{want}
		close(rc)
		close(ec)
		return rc, ec
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.sizes;
var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var d = list[0];
var want = {
	slug: "lol",
	memory: 1,
	vcpus: 2,
	disk: 2,
	price_monthly: 1.0,
	price_hourly: 0.1,
	regions: ["lol"],
	available: true,
	transfer: 1.0,
}
equals(d, want, "should have proper object");
`)
}
