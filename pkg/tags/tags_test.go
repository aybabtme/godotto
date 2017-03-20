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
	assert(pkg.tag_resources != null, "tag_resources function should be defined.");
	assert(pkg.untag_resources != null, "tag_resources function should be defined.");
	assert(pkg.get != null, "get function should be defined");
	`)
}

func TestTagThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockTags.CreateFn = func(_ context.Context, _ string, _ ...tags.CreateOpt) (tags.Tag, error) {
		return nil, errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;

		var name = "test";

		var tag = {
			name: name	
		};

		[
		{ name: "create", fn: function() { pkg.create(tag) } }	
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

func TestTagTagResources(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockTags.TagFn = func(_ context.Context, name string, res []godo.Resource) error {
		// Won't return an error
		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;
		pkg.tag_resources({
			name: "test",
			resources: [
				{
					id: "12345567",
					type: "droplet"
				}	
			]
		});
	`)
}

func TestTagUntagResources(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockTags.UntagFn = func(_ context.Context, name string, res []godo.Resource) error {
		// Won't return an error
		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.tags;
		pkg.untag_resources({
			name: "test",
			resources: [
				{
					id: "12345567",
					type: "droplet"
				}	
			]
		});
	`)
}
