package images_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/images"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.images;

assert(pkg != null, "package should be loaded");
assert(pkg.list != null, "list function should be defined");
assert(pkg.list_distribution != null, "list_distribution function should be defined");
assert(pkg.list_application != null, "list_application function should be defined");
assert(pkg.list_user != null, "list_user function should be defined");
assert(pkg.get != null, "get function should be defined");
assert(pkg.update != null, "update function should be defined");
assert(pkg.delete != null, "delete function should be defined");
    `)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	listfn := func(_ context.Context) (<-chan images.Image, <-chan error) {
		lc := make(chan images.Image)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}
	cloud.MockImages.ListFn = listfn
	cloud.MockImages.ListApplicationFn = listfn
	cloud.MockImages.ListDistributionFn = listfn
	cloud.MockImages.ListUserFn = listfn
	cloud.MockImages.GetByIDFn = func(_ context.Context, _ int) (images.Image, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockImages.GetBySlugFn = func(_ context.Context, _ string) (images.Image, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockImages.UpdateFn = func(_ context.Context, _ int, _ ...images.UpdateOpt) (images.Image, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockImages.DeleteFn = func(_ context.Context, _ int) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.images;

[
	{ name: "list",              fn: function() { pkg.list() } },
	{ name: "list_distribution", fn: function() { pkg.list_distribution() } },
	{ name: "list_application",  fn: function() { pkg.list_application() } },
	{ name: "list_user",         fn: function() { pkg.list_user() } },
	{ name: "get",               fn: function() { pkg.get(42) } },
	{ name: "update",            fn: function() { pkg.update(42, {name:"lol"}) } },
	{ name: "delete",            fn: function() { pkg.delete(42) } },

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

type image struct {
	*godo.Image
}

func (k *image) Struct() *godo.Image { return k.Image }

func TestList(t *testing.T) {
	listfn := func(_ context.Context) (<-chan images.Image, <-chan error) {
		lc := make(chan images.Image, 1)
		lc <- &image{&godo.Image{
			ID:           42,
			Name:         "derp",
			Type:         "herp",
			Distribution: "coreos",
			Slug:         "coreos-stable",
			Public:       true,
			Regions:      []string{"atlantis"},
		}}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	cloud := mockcloud.Client(nil)
	cloud.MockImages.ListFn = listfn
	cloud.MockImages.ListApplicationFn = listfn
	cloud.MockImages.ListDistributionFn = listfn
	cloud.MockImages.ListUserFn = listfn
	vmtest.Run(t, cloud, `
var pkg = cloud.images;

[
	pkg.list,
	pkg.list_distribution,
	pkg.list_application,
	pkg.list_user,
].forEach(function(fn) {
	var list = fn();
	assert(list != null, "should have received a list");
	assert(list.length > 0, "should have received some elements")
	var d = list[0];
	var want = {
		id:            42,
		name:          "derp",
		type:          "herp",
		distribution:  "coreos",
		slug:          "coreos-stable",
		public:        true,
		regions:       ["atlantis"],
		min_disk_size: 0,
	};
	equals(d, want, "should have proper object");
});
`)
}

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockImages.UpdateFn = func(_ context.Context, id int, _ ...images.UpdateOpt) (images.Image, error) {
		return &image{&godo.Image{
			ID:           id,
			Name:         "derp",
			Type:         "herp",
			Distribution: "coreos",
			Slug:         "coreos-stable",
			Public:       true,
			Regions:      []string{"atlantis"},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.images;

var want = {
	id:           42,
	name:         "derp",
	type:         "herp",
	distribution: "coreos",
	slug:         "coreos-stable",
	public:       true,
	regions:      ["atlantis"],
	min_disk_size: 0,
};
var d = pkg.update(42, {name:"lol"})
equals(d, want, "should have proper object");
`)
}

func TestUpdate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockImages.GetByIDFn = func(_ context.Context, id int) (images.Image, error) {
		return &image{&godo.Image{
			ID:           id,
			Name:         "derp",
			Type:         "herp",
			Distribution: "coreos",
			Slug:         "coreos-stable",
			Public:       true,
			Regions:      []string{"atlantis"},
		}}, nil
	}
	cloud.MockImages.GetBySlugFn = func(_ context.Context, slug string) (images.Image, error) {
		return &image{&godo.Image{
			ID:           42,
			Name:         "derp",
			Type:         "herp",
			Distribution: "coreos",
			Slug:         slug,
			Public:       true,
			Regions:      []string{"atlantis"},
		}}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.images;

var want = {
	id:           42,
	name:         "derp",
	type:         "herp",
	distribution: "coreos",
	slug:         "coreos-stable",
	public:       true,
	regions:      ["atlantis"],
	min_disk_size: 0,
};
var d = pkg.get(42)
equals(d, want, "should have proper object");
var d = pkg.get("coreos-stable")
equals(d, want, "should have proper object");
`)
}

func TestDelete(t *testing.T) {
	wantID := 42
	cloud := mockcloud.Client(nil)
	cloud.MockImages.DeleteFn = func(_ context.Context, gotID int) error {
		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}
		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.images;

pkg.delete(42);
`)
}
