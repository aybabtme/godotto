package regions_test

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/regions"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

type region struct {
	*godo.Region
}

func (k *region) Struct() *godo.Region { return k.Region }

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.regions;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)

	want := &godo.Region{Slug: "lol", Name: "lol", Available: true, Features: []string{"lol"}}

	cloud.MockRegions.ListFn = func(_ context.Context) (<-chan regions.Region, <-chan error) {
		rc, ec := make(chan regions.Region, 1), make(chan error, 0)
		rc <- &region{want}
		close(rc)
		close(ec)
		return rc, ec
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.regions;
var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var region = list[0];
assert(region.slug, "should have field 'slug'");
assert(region.name, "should have field 'name'");
assert(region.sizes, "should have field 'sizes'");
assert(region.available, "should have field 'available'");
assert(region.features, "should have field 'features'");
    `)
}
