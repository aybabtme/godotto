package do

import "testing"

func TestStubDropletCreate(t *testing.T) {
	cloud, done := stubClient()
	defer done()

	var (
		wantName   = "myname"
		wantRegion = "nyc3"
		wantSize   = "4gb"
		wantImage  = "coreos-stable"
	)

	d, err := cloud.Droplets().Create(wantName, wantRegion, wantSize, wantImage)
	if err != nil {
		t.Fatal(err)
	}
	drop := d.Struct()
	if want, got := wantName, drop.Name; want != got {
		t.Fatalf("want %v got %v", want, got)
	}
	if want, got := wantRegion, drop.Region.Slug; want != got {
		t.Fatalf("want %v got %v", want, got)
	}
	if want, got := wantSize, drop.SizeSlug; want != got {
		t.Fatalf("want %v got %v", want, got)
	}
	if want, got := wantImage, drop.Image.Slug; want != got {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestStubDropletDelete(t *testing.T) {
	cloud, done := stubClient()
	defer done()

	d, err := cloud.Droplets().Create("myname", "nyc3", "4gb", "coreos-stable")
	if err != nil {
		t.Fatal(err)
	}
	if err := cloud.Droplets().Delete(d.Struct().ID); err != nil {
		t.Fatal(err)
	}
}
