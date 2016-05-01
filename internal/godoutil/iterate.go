package godoutil

import (
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

func IterateList(ctx context.Context, fn func(*godo.ListOptions) (*godo.Response, error)) error {
	opt := &godo.ListOptions{Page: 1, PerPage: 200}
	for {
		resp, err := fn(opt)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if resp.Links != nil && !resp.Links.IsLastPage() {
			opt.Page++
		} else {
			return nil
		}
	}
}
