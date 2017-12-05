package floatingips

import (
	"context"

	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
)

// An ActionClient can interact with the DigitalOcean FloatingIPActions service.
type ActionClient interface {
	Assign(ctx context.Context, ip string, dropletID int) error
	Unassign(ctx context.Context, ip string) error
}

type actionClient struct {
	g *godo.Client
}

func (svc *actionClient) Assign(ctx context.Context, ip string, dropletID int) error {
	_, resp, err := svc.g.FloatingIPActions.Assign(ctx, ip, dropletID)
	if err != nil {
		return err
	}

	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Unassign(ctx context.Context, ip string) error {
	_, resp, err := svc.g.FloatingIPActions.Unassign(ctx, ip)
	if err != nil {
		return err
	}

	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}
