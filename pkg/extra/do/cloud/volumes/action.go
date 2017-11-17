package volumes

import (
	"context"

	"github.com/aybabtme/godotto/internal/godoutil"
	"github.com/digitalocean/godo"
)

// An ActionClient can interact with the DigitalOcean StorageAction service.
type ActionClient interface {
	Attach(ctx context.Context, ip string, dropletID int) error
	DetachByDropletID(ctx context.Context, ip string, dropletID int) error
}

type actionClient struct {
	g *godo.Client
}

func (svc *actionClient) Attach(ctx context.Context, driveID string, dropletID int) error {
	action, _, err := svc.g.StorageActions.Attach(ctx, driveID, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForAction(ctx, svc.g, action)
}

func (svc *actionClient) DetachByDropletID(ctx context.Context, volumeID string, dropletID int) error {
	action, _, err := svc.g.StorageActions.DetachByDropletID(ctx, volumeID, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForAction(ctx, svc.g, action)
}
