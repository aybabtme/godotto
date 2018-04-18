package keys_test

import (
	"context"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

type key struct {
	*godo.Key
}

func (k *key) Struct() *godo.Key { return k.Key }

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.keys;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
    `)
}

func TestList(t *testing.T) {
	cloud := mockcloud.Client(nil)

	want := &godo.Key{ID: 1, Name: "lol", Fingerprint: "lol", PublicKey: "lol"}

	cloud.MockKeys.ListFn = func(_ context.Context) (<-chan keys.Key, <-chan error) {
		kc, ec := make(chan keys.Key, 1), make(chan error, 0)
		kc <- &key{want}
		close(kc)
		close(ec)
		return kc, ec
	}

	vmtest.Run(t, cloud, `
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
