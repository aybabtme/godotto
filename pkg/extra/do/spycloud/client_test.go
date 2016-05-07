package spycloud_test

import (
	"reflect"
	"testing"

	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/aybabtme/godotto/pkg/extra/do/spycloud"
	"github.com/digitalocean/godo"
)

type droplet struct{ *godo.Droplet }

func (d *droplet) Struct() *godo.Droplet {
	return d.Droplet
}

func TestSpy(t *testing.T) {
	want := &godo.Droplet{ID: 42, Name: "hello"}

	mock := mockcloud.Client(nil)
	mock.MockDroplets.CreateFn = func(name, _, _, _ string, _ ...droplets.CreateOpt) (droplets.Droplet, error) {
		return &droplet{&godo.Droplet{ID: want.ID, Name: name}}, nil
	}
	mock.MockDroplets.DeleteFn = func(id int) error { return nil }

	cloud, spy := spycloud.Client(mock)
	d, err := cloud.Droplets().Create(want.Name, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if want, got := want, d.Struct(); !reflect.DeepEqual(want, got) {
		t.Errorf("want=%#v", want)
		t.Errorf("got =%#v", got)
	}

	seen := 0
	spy(spycloud.Droplets(func(got *godo.Droplet) {
		seen++
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want=%#v", want)
			t.Errorf("got =%#v", got)
		}
	}))
	if seen != 1 {
		t.Errorf("want seen %d, got %d", 1, seen)
	}

	if err := cloud.Droplets().Delete(d.Struct().ID); err != nil {
		t.Fatal(err)
	}

	spy(spycloud.Droplets(func(got *godo.Droplet) {
		t.Errorf("should not see %#v", got)
	}))
}
