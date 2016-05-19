package droplets

import (
	"fmt"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/context"
)

func applyAction(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	svc := actionSvc{
		ctx: ctx,
		svc: client.Droplets().Actions(),
	}
	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"shutdown", svc.shutdown},
		{"power_off", svc.powerOff},
		{"power_on", svc.powerOn},
		{"power_cycle", svc.powerCycle},
		{"reboot", svc.reboot},
		{"restore", svc.restore},
		{"resize", svc.resize},
		{"rename", svc.rename},
		{"snapshot", svc.snapshot},
		{"enable_backups", svc.enableBackups},
		{"disable_backups", svc.disableBackups},
		{"password_reset", svc.passwordReset},
		{"change_kernel", svc.changeKernel},
		{"enable_ipv6", svc.enableIPv6},
		{"enable_private_networking", svc.enablePrivateNetworking},
		{"upgrade", svc.upgrade},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

type actionSvc struct {
	ctx context.Context
	svc droplets.ActionClient
}

func (svc *actionSvc) shutdown(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.Shutdown(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) powerOff(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.PowerOff(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) powerOn(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.PowerOn(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) powerCycle(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.PowerCycle(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) reboot(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.Reboot(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) restore(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	imageID := argImageID(vm, all.Argument(1))
	err := svc.svc.Restore(svc.ctx, dropletID, imageID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) resize(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	sizeSlug := argSizeSlug(vm, all.Argument(1))
	resizeDisk := ottoutil.Bool(vm, all.Argument(2))
	err := svc.svc.Resize(svc.ctx, dropletID, sizeSlug, resizeDisk)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) rename(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	name := ottoutil.String(vm, all.Argument(1))
	err := svc.svc.Rename(svc.ctx, dropletID, name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) snapshot(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	name := ottoutil.String(vm, all.Argument(1))
	err := svc.svc.Snapshot(svc.ctx, dropletID, name)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) enableBackups(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.EnableBackups(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) disableBackups(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.DisableBackups(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) passwordReset(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.PasswordReset(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) changeKernel(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	kernelID := argKernelID(vm, all.Argument(1))
	err := svc.svc.ChangeKernel(svc.ctx, dropletID, kernelID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) enableIPv6(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.EnableIPv6(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) enablePrivateNetworking(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.EnablePrivateNetworking(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func (svc *actionSvc) upgrade(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dropletID := argDropletID(vm, all.Argument(0))
	err := svc.svc.Upgrade(svc.ctx, dropletID)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}
