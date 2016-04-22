package regions_test

import (
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
)

func TestApply(t *testing.T) {
	vmtest.Run(t, `
var pkg = cloud.regions;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

func TestList(t *testing.T) {
	vmtest.Run(t, `
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
