package tags_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/tags"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

var testTag *godo.Tag = &godo.Tag{
	Name: "test",
}

func TestTagApply(t *testing.T) {
	cloud := mockcloud.Client(nil)

	vmtest.Run(t, cloud, `
	var pkg = cloud.tags;
	assert(pkg != null, "package should be loaded");
	assert(pkg.create != null, "create function should be defined");
	assert(pkg.tag_resources != null, "tag_resources function should be defined");
	assert(pkg.untag_resources != null, "tag_resources function should be defined");
	assert(pkg.get != null, "get function should be defined");
	assert(pkg.list != null, "list function should be defined");
	assert(pkg.delete != null, "delete function should be defined");
	`)
}

func TestTagThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockTags.ListFn = func(_ context.Context) (<-chan tags.Tag, <-chan error) {
		lc := make(chan tags.Tag)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}

	cloud.MockTags.CreateFn = func(_ context.Context, _ string, _ ...tags.CreateOpt) (tags.Tag, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockTags.GetFn = func(_ context.Context, _ string) (tags.Tag, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockTags.TagFn = func(_ context.Context, _ string, _ []godo.Resource) error {
		return errors.New("throw me")
	}

	cloud.MockTags.UntagFn = func(_ context.Context, _ string, _ []godo.Resource) error {
		return errors.New("throw me")
	}

	cloud.MockTags.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;

		var name = "test";

		var tag = {
			name: name	
		};

		var testTag = {
			name: name,
			resources: []
		};

		[
		{ name: "create", fn: function() { pkg.create(tag) } },	
		{ name: "get", fn: function() { pkg.get(name) } },
		{ name: "list", fn: function() { pkg.list() } },
		{ name: "delete", fn: function() { pkg.delete() } },
		{ name: "tag_resources", fn: function() { pkg.tag_resources(testTag) } },
		{ name: "untag_resources", fn: function() { pkg.untag_resources(testTag) } }
		].forEach(function(kv) {
			var name = kv.name;
			var fn = kv.fn;
			try {
				fn(); throw "dont catch me";
			} catch (e) {
				equals("throw me", e.message, name + " should send the right exception");	
			}
		});
	`)
}

type tag struct {
	*godo.Tag
}

func (k *tag) Struct() *godo.Tag { return k.Tag }

func TestTagCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockTags.CreateFn = func(_ context.Context, name string, _ ...tags.CreateOpt) (tags.Tag, error) {
		return &tag{testTag}, nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;

		var tag = pkg.create({
			name: "test"	
		});

		var want = {
			name: "test"
		}

		equals(tag, want, "should have proper object");
	`)
}

func TestTagGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockTags.GetFn = func(_ context.Context, name string) (tags.Tag, error) {
		return &tag{testTag}, nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;

		var tag = pkg.get({
			name: "test"	
		});

		var want = {
			name: "test"
		}

		equals(tag, want, "should have proper object");
	`)
}

func TestTagList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockTags.ListFn = func(_ context.Context) (<-chan tags.Tag, <-chan error) {
		lc := make(chan tags.Tag, 1)
		lc <- &tag{testTag}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;

		var list = pkg.list();
		assert(list != null, "should have received a list");
		assert(list.length > 0, "should have received some elements")

		var want = {
			name: "test"	
		};

		var t = list[0];
		equals(t, want, "should have proper object");
	`)
}

func TestTagDelete(t *testing.T) {
	wantName := "test"
	cloud := mockcloud.Client(nil)
	cloud.MockTags.DeleteFn = func(_ context.Context, gotName string) error {
		if gotName != wantName {
			t.Fatalf("want %v got %v", wantName, gotName)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;
		pkg.delete("test");
	`)
}

func TestTagTagResources(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantName := "test"
	cloud.MockTags.TagFn = func(_ context.Context, gotName string, res []godo.Resource) error {
		if gotName != wantName {
			t.Fatalf("want %v got %v", wantName, gotName)
		}

		wantRes := godo.Resource{
			ID:   "1234567",
			Type: "droplet",
		}

		gotRes := res[0]
		if wantRes != gotRes {
			t.Fatalf("want %v got %v", wantRes, gotRes)
		}

		// Won't return an error
		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;
		pkg.tag_resources({
			name: "test",
			resources: [
				{
					id: "1234567",
					type: "droplet"
				}	
			]
		});
	`)
}

func TestTagUntagResources(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantName := "test"
	cloud.MockTags.UntagFn = func(_ context.Context, gotName string, res []godo.Resource) error {
		if gotName != wantName {
			t.Fatalf("want %v got %v", wantName, gotName)
		}

		wantRes := godo.Resource{
			ID:   "1234567",
			Type: "droplet",
		}

		gotRes := res[0]
		if wantRes != gotRes {
			t.Fatalf("want %v got %v", wantRes, gotRes)
		}
		// Won't return an error
		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;
		pkg.untag_resources({
			name: "test",
			resources: [
				{
					id: "1234567",
					type: "droplet"
				}	
			]
		});
	`)
}
