package godoutil

import (
	"context"

	"github.com/digitalocean/godo"
)

func IterateList(ctx context.Context, fn func(context.Context, *godo.ListOptions) (*godo.Response, error)) error {
	opt := &godo.ListOptions{Page: 1, PerPage: 200}
	for {
		resp, err := fn(ctx, opt)
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
