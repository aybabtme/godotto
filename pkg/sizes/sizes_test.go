package sizes_test

import (
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
)

func TestApply(t *testing.T) {
	vmtest.Run(t, `
var pkg = cloud.sizes;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

func TestList(t *testing.T) {
	vmtest.Run(t, `
var pkg = cloud.sizes;
var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var size = list[0];
assert(size.slug, "should have a slug")
assert(size.memory, "should have a memory")
assert(size.vcpus, "should have a vcpus")
assert(size.disk, "should have a disk")
assert(size.price_monthly, "should have a price_monthly")
assert(size.price_hourly, "should have a price_hourly")
assert(size.regions, "should have a regions")
assert(size.available, "should have a available")
assert(size.transfer, "should have a transfer")
    `)
}
