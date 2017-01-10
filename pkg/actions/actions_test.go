package actions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/actions"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.actions;

assert(pkg != null, "package should be loaded");
assert(pkg.get != null, "get function should be defined");
assert(pkg.list != null, "list function should be defined");
    `)
}

type action struct {
	*godo.Action
}

func (k *action) Struct() *godo.Action { return k.Action }

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)

	want := &godo.Action{
		ID:           42,
		Status:       "my_status",
		Type:         "my_type",
		StartedAt:    &godo.Timestamp{time.Date(1987, 03, 24, 10, 30, 00, 00, time.UTC)},
		CompletedAt:  &godo.Timestamp{time.Date(1988, 03, 24, 10, 30, 00, 00, time.UTC)},
		ResourceID:   9000,
		ResourceType: "my_resource_type",
		Region:       &godo.Region{Slug: "my_region"},
		RegionSlug:   "my_region_slug",
	}

	cloud.MockActions.GetFn = func(_ context.Context, id int) (actions.Action, error) {
		return &action{want}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.actions;

var a = pkg.get(42);

var want = {
	"id": 42,
	"status": "my_status",
	"type": "my_type",
	"started_at": "1987-03-24T10:30:00Z",
	"completed_at": "1988-03-24T10:30:00Z",
	"resource_id": 9000,
	"resource_type": "my_resource_type",
	"region_slug": "my_region_slug"
};
equals(a, want, "should get proper object");
    `)
}

func TestList(t *testing.T) {

	want := &godo.Action{
		ID:           42,
		Status:       "my_status",
		Type:         "my_type",
		StartedAt:    &godo.Timestamp{time.Date(1987, 03, 24, 10, 30, 00, 00, time.UTC)},
		CompletedAt:  &godo.Timestamp{time.Date(1988, 03, 24, 10, 30, 00, 00, time.UTC)},
		ResourceID:   9000,
		ResourceType: "my_resource_type",
		Region:       &godo.Region{Slug: "my_region"},
		RegionSlug:   "my_region_slug",
	}

	cloud := mockcloud.Client(nil)
	cloud.MockActions.ListFn = func(_ context.Context) (<-chan actions.Action, <-chan error) {
		lc := make(chan actions.Action, 1)
		lc <- &action{want}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.actions;

var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements")

var a = list[0];

var want = {
	"id": 42,
	"status": "my_status",
	"type": "my_type",
	"started_at": "1987-03-24T10:30:00Z",
	"completed_at": "1988-03-24T10:30:00Z",
	"resource_id": 9000,
	"resource_type": "my_resource_type",
	"region_slug": "my_region_slug"
};
equals(a, want, "should get proper object");
`)
}

func TestThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockActions.GetFn = func(_ context.Context, _ int) (actions.Action, error) {
		return nil, errors.New("throw me")
	}
	cloud.MockActions.ListFn = func(_ context.Context) (<-chan actions.Action, <-chan error) {
		lc := make(chan actions.Action)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.actions;

try {
	pkg.get(1);
	throw "dont catch me";
} catch (e) {
	equals("throw me", e.message, "should send the right exception");
}

try {
	pkg.list();
	throw "dont catch me";
} catch (e) {
	equals("throw me", e.message, "should send the right exception again");
}
    `)
}
