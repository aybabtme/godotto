package floatingips

import (
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
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
	_, _, err := svc.g.FloatingIPActions.Assign(ip, dropletID)
	return err
}

func (svc *actionClient) Unassign(ctx context.Context, ip string) error {
	_, _, err := svc.g.FloatingIPActions.Unassign(ip)
	return err
}
