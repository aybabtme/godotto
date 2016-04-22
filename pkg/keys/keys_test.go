package keys_test

import (
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
)

func TestApply(t *testing.T) {
	vmtest.Run(t, `
var pkg = cloud.keys;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

func TestList(t *testing.T) {
	vmtest.Run(t, `
var pkg = cloud.keys;
var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var key = list[0];
assert(key.id, "should have a 'id' field");
assert(key.name, "should have a 'name' field");
assert(key.fingerprint, "should have a 'fingerprint' field");
assert(key.public_key, "should have a 'public_key' field");
    `)
}
