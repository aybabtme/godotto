package volumes

import (
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

// An ActionClient can interact with the DigitalOcean StorageAction service.
type ActionClient interface {
	Attach(ctx context.Context, ip string, dropletID int) error
	Detach(ctx context.Context, ip string) error
}

type actionClient struct {
	g *godo.Client
}

func (svc *actionClient) Attach(ctx context.Context, driveID string, dropletID int) error {
	_, _, err := svc.g.StorageActions.Attach(driveID, dropletID)
	return err
}

func (svc *actionClient) Detach(ctx context.Context, driveID string) error {
	_, _, err := svc.g.StorageActions.Detach(driveID)
	return err
}
