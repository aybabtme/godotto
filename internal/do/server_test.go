package do

import (
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/digitalocean/godo"
)

func stubClient() (cloud.Client, func()) {
	u, done := Stub()
	client := godo.NewClient(nil)
	client.BaseURL = u
	return cloud.New(cloud.UseGodo(client)), done
}
