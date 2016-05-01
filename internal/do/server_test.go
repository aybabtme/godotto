package do

import (
	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/digitalocean/godo"
)

func stubClient() (cloud.Client, func()) {
	u, done := Stub()
	return cloud.New(cloud.UseGodoOpts(func(c *godo.Client) {
		c.BaseURL = u
	})), done

}
