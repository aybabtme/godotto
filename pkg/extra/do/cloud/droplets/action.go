package droplets

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/godoutil"
	"github.com/digitalocean/godo"
)

// An ActionClient can interact with the DigitalOcean DropletActions service.
type ActionClient interface {
	Shutdown(ctx context.Context, dropletID int) error
	PowerOff(ctx context.Context, dropletID int) error
	PowerOn(ctx context.Context, dropletID int) error
	PowerCycle(ctx context.Context, dropletID int) error
	Reboot(ctx context.Context, dropletID int) error
	Restore(ctx context.Context, dropletID, imageID int) error
	Resize(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error
	Rename(ctx context.Context, dropletID int, name string) error
	Snapshot(ctx context.Context, dropletID int, name string) error
	EnableBackups(ctx context.Context, dropletID int) error
	DisableBackups(ctx context.Context, dropletID int) error
	PasswordReset(ctx context.Context, dropletID int) error
	RebuildByImageID(ctx context.Context, dropletID int, imageID int) error
	RebuildByImageSlug(ctx context.Context, dropletID int, imageSlug string) error
	ChangeKernel(ctx context.Context, dropletID int, kernelID int) error
	EnableIPv6(ctx context.Context, dropletID int) error
	EnablePrivateNetworking(ctx context.Context, dropletID int) error
}

type actionClient struct {
	g *godo.Client
}

func (svc *actionClient) Shutdown(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.Shutdown(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) PowerOff(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.PowerOff(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) PowerOn(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.PowerOn(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) PowerCycle(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.PowerCycle(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Reboot(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.Reboot(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Restore(ctx context.Context, dropletID, imageID int) error {
	_, resp, err := svc.g.DropletActions.Restore(ctx, dropletID, imageID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Resize(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error {
	_, resp, err := svc.g.DropletActions.Resize(ctx, dropletID, sizeSlug, resizeDisk)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Rename(ctx context.Context, dropletID int, name string) error {
	_, resp, err := svc.g.DropletActions.Rename(ctx, dropletID, name)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) Snapshot(ctx context.Context, dropletID int, name string) error {
	_, resp, err := svc.g.DropletActions.Snapshot(ctx, dropletID, name)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) EnableBackups(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.EnableBackups(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) DisableBackups(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.DisableBackups(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) PasswordReset(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.PasswordReset(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) RebuildByImageID(ctx context.Context, dropletID int, imageID int) error {
	_, resp, err := svc.g.DropletActions.RebuildByImageID(ctx, dropletID, imageID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) RebuildByImageSlug(ctx context.Context, dropletID int, imageSlug string) error {
	_, resp, err := svc.g.DropletActions.RebuildByImageSlug(ctx, dropletID, imageSlug)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) ChangeKernel(ctx context.Context, dropletID int, kernelID int) error {
	_, resp, err := svc.g.DropletActions.ChangeKernel(ctx, dropletID, kernelID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) EnableIPv6(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.EnableIPv6(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}

func (svc *actionClient) EnablePrivateNetworking(ctx context.Context, dropletID int) error {
	_, resp, err := svc.g.DropletActions.EnablePrivateNetworking(ctx, dropletID)
	if err != nil {
		return err
	}
	return godoutil.WaitForActions(ctx, svc.g, resp.Links)
}
