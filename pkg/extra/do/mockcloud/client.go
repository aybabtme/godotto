/*
Package mockcloud allows mocking a cloud.Client.

    var real cloud.Client
    mock := mockcloud.Client(real)
    mock.MockDroplets.Get = func(id int) (droplets.Droplet, error) {
        panic("invoked the mock!")
    }

That's it!
*/
package mockcloud

import (
	"context"

	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/accounts"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/actions"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/images"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/loadbalancers"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/regions"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/sizes"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/tags"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/volumes"
	"github.com/digitalocean/godo"
)

// Client

type Mock struct {
	wrap              cloud.Client
	MockDroplets      *MockDroplets
	MockAccounts      *MockAccounts
	MockActions       *MockActions
	MockDomains       *MockDomains
	MockImages        *MockImages
	MockKeys          *MockKeys
	MockRegions       *MockRegions
	MockSizes         *MockSizes
	MockFloatingIPs   *MockFloatingIPs
	MockVolumes       *MockVolumes
	MockTags          *MockTags
	MockLoadBalancers *MockLoadBalancers
}

func Client(client cloud.Client) *Mock {
	return &Mock{wrap: client,
		MockDroplets:      &MockDroplets{wrap: client, MockDropletActions: &MockDropletActions{wrap: client}},
		MockAccounts:      &MockAccounts{wrap: client},
		MockActions:       &MockActions{wrap: client},
		MockDomains:       &MockDomains{wrap: client},
		MockImages:        &MockImages{wrap: client},
		MockKeys:          &MockKeys{wrap: client},
		MockRegions:       &MockRegions{wrap: client},
		MockSizes:         &MockSizes{wrap: client},
		MockFloatingIPs:   &MockFloatingIPs{wrap: client, MockFloatingIPActions: &MockFloatingIPActions{wrap: client}},
		MockVolumes:       &MockVolumes{wrap: client, MockVolumeActions: &MockVolumeActions{wrap: client}},
		MockTags:          &MockTags{wrap: client},
		MockLoadBalancers: &MockLoadBalancers{wrap: client},
	}
}

func (mock *Mock) Droplets() droplets.Client           { return mock.MockDroplets }
func (mock *Mock) Accounts() accounts.Client           { return mock.MockAccounts }
func (mock *Mock) Actions() actions.Client             { return mock.MockActions }
func (mock *Mock) Domains() domains.Client             { return mock.MockDomains }
func (mock *Mock) Images() images.Client               { return mock.MockImages }
func (mock *Mock) Keys() keys.Client                   { return mock.MockKeys }
func (mock *Mock) Regions() regions.Client             { return mock.MockRegions }
func (mock *Mock) Sizes() sizes.Client                 { return mock.MockSizes }
func (mock *Mock) FloatingIPs() floatingips.Client     { return mock.MockFloatingIPs }
func (mock *Mock) Volumes() volumes.Client             { return mock.MockVolumes }
func (mock *Mock) Tags() tags.Client                   { return mock.MockTags }
func (mock *Mock) LoadBalancers() loadbalancers.Client { return mock.MockLoadBalancers }

// Droplets

type MockDroplets struct {
	wrap               cloud.Client
	CreateFn           func(ctx context.Context, name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error)
	GetFn              func(ctx context.Context, id int) (droplets.Droplet, error)
	DeleteFn           func(ctx context.Context, id int) error
	ListFn             func(ctx context.Context) (<-chan droplets.Droplet, <-chan error)
	MockDropletActions *MockDropletActions
}

func (mock *MockDroplets) Create(ctx context.Context, name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, name, region, size, image, opts...)
	}
	return mock.wrap.Droplets().Create(ctx, name, region, size, image, opts...)
}
func (mock *MockDroplets) Get(ctx context.Context, id int) (droplets.Droplet, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx, id)
	}
	return mock.wrap.Droplets().Get(ctx, id)
}
func (mock *MockDroplets) Delete(ctx context.Context, id int) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, id)
	}
	return mock.wrap.Droplets().Delete(ctx, id)
}
func (mock *MockDroplets) List(ctx context.Context) (<-chan droplets.Droplet, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Droplets().List(ctx)
}
func (mock *MockDroplets) Actions() droplets.ActionClient {
	if mock.MockDropletActions != nil {
		return mock.MockDropletActions
	}
	return mock.wrap.Droplets().Actions()
}

// Droplet Actions

type MockDropletActions struct {
	wrap                      cloud.Client
	ShutdownFn                func(ctx context.Context, dropletID int) error
	PowerOffFn                func(ctx context.Context, dropletID int) error
	PowerOnFn                 func(ctx context.Context, dropletID int) error
	PowerCycleFn              func(ctx context.Context, dropletID int) error
	RebootFn                  func(ctx context.Context, dropletID int) error
	RestoreFn                 func(ctx context.Context, dropletID, imageID int) error
	ResizeFn                  func(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error
	RenameFn                  func(ctx context.Context, dropletID int, name string) error
	SnapshotFn                func(ctx context.Context, dropletID int, name string) error
	EnableBackupsFn           func(ctx context.Context, dropletID int) error
	DisableBackupsFn          func(ctx context.Context, dropletID int) error
	PasswordResetFn           func(ctx context.Context, dropletID int) error
	RebuildByImageIDFn        func(ctx context.Context, dropletID int, imageID int) error
	RebuildByImageSlugFn      func(ctx context.Context, dropletID int, imageSlug string) error
	ChangeKernelFn            func(ctx context.Context, dropletID int, kernelID int) error
	EnableIPv6Fn              func(ctx context.Context, dropletID int) error
	EnablePrivateNetworkingFn func(ctx context.Context, dropletID int) error
	UpgradeFn                 func(ctx context.Context, dropletID int) error
}

func (mock *MockDropletActions) Shutdown(ctx context.Context, dropletID int) error {
	if mock.ShutdownFn != nil {
		return mock.ShutdownFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().Shutdown(ctx, dropletID)
}
func (mock *MockDropletActions) PowerOff(ctx context.Context, dropletID int) error {
	if mock.PowerOffFn != nil {
		return mock.PowerOffFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().PowerOff(ctx, dropletID)
}
func (mock *MockDropletActions) PowerOn(ctx context.Context, dropletID int) error {
	if mock.PowerOnFn != nil {
		return mock.PowerOnFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().PowerOn(ctx, dropletID)
}
func (mock *MockDropletActions) PowerCycle(ctx context.Context, dropletID int) error {
	if mock.PowerCycleFn != nil {
		return mock.PowerCycleFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().PowerCycle(ctx, dropletID)
}
func (mock *MockDropletActions) Reboot(ctx context.Context, dropletID int) error {
	if mock.RebootFn != nil {
		return mock.RebootFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().Reboot(ctx, dropletID)
}
func (mock *MockDropletActions) Restore(ctx context.Context, dropletID, imageID int) error {
	if mock.RestoreFn != nil {
		return mock.RestoreFn(ctx, dropletID, imageID)
	}
	return mock.wrap.Droplets().Actions().Restore(ctx, dropletID, imageID)
}
func (mock *MockDropletActions) Resize(ctx context.Context, dropletID int, sizeSlug string, resizeDisk bool) error {
	if mock.ResizeFn != nil {
		return mock.ResizeFn(ctx, dropletID, sizeSlug, resizeDisk)
	}
	return mock.wrap.Droplets().Actions().Resize(ctx, dropletID, sizeSlug, resizeDisk)
}
func (mock *MockDropletActions) Rename(ctx context.Context, dropletID int, name string) error {
	if mock.RenameFn != nil {
		return mock.RenameFn(ctx, dropletID, name)
	}
	return mock.wrap.Droplets().Actions().Rename(ctx, dropletID, name)
}
func (mock *MockDropletActions) Snapshot(ctx context.Context, dropletID int, name string) error {
	if mock.SnapshotFn != nil {
		return mock.SnapshotFn(ctx, dropletID, name)
	}
	return mock.wrap.Droplets().Actions().Snapshot(ctx, dropletID, name)
}
func (mock *MockDropletActions) EnableBackups(ctx context.Context, dropletID int) error {
	if mock.EnableBackupsFn != nil {
		return mock.EnableBackupsFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().EnableBackups(ctx, dropletID)
}
func (mock *MockDropletActions) DisableBackups(ctx context.Context, dropletID int) error {
	if mock.DisableBackupsFn != nil {
		return mock.DisableBackupsFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().DisableBackups(ctx, dropletID)
}
func (mock *MockDropletActions) PasswordReset(ctx context.Context, dropletID int) error {
	if mock.PasswordResetFn != nil {
		return mock.PasswordResetFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().PasswordReset(ctx, dropletID)
}
func (mock *MockDropletActions) RebuildByImageID(ctx context.Context, dropletID int, imageID int) error {
	if mock.RebuildByImageIDFn != nil {
		return mock.RebuildByImageIDFn(ctx, dropletID, imageID)
	}
	return mock.wrap.Droplets().Actions().RebuildByImageID(ctx, dropletID, imageID)
}
func (mock *MockDropletActions) RebuildByImageSlug(ctx context.Context, dropletID int, imageSlug string) error {
	if mock.RebuildByImageSlugFn != nil {
		return mock.RebuildByImageSlugFn(ctx, dropletID, imageSlug)
	}
	return mock.wrap.Droplets().Actions().RebuildByImageSlug(ctx, dropletID, imageSlug)
}
func (mock *MockDropletActions) ChangeKernel(ctx context.Context, dropletID int, kernelID int) error {
	if mock.ChangeKernelFn != nil {
		return mock.ChangeKernelFn(ctx, dropletID, kernelID)
	}
	return mock.wrap.Droplets().Actions().ChangeKernel(ctx, dropletID, kernelID)
}
func (mock *MockDropletActions) EnableIPv6(ctx context.Context, dropletID int) error {
	if mock.EnableIPv6Fn != nil {
		return mock.EnableIPv6Fn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().EnableIPv6(ctx, dropletID)
}
func (mock *MockDropletActions) EnablePrivateNetworking(ctx context.Context, dropletID int) error {
	if mock.EnablePrivateNetworkingFn != nil {
		return mock.EnablePrivateNetworkingFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().EnablePrivateNetworking(ctx, dropletID)
}
func (mock *MockDropletActions) Upgrade(ctx context.Context, dropletID int) error {
	if mock.UpgradeFn != nil {
		return mock.UpgradeFn(ctx, dropletID)
	}
	return mock.wrap.Droplets().Actions().Upgrade(ctx, dropletID)
}

// Accounts

type MockAccounts struct {
	wrap  cloud.Client
	GetFn func(context.Context) (accounts.Account, error)
}

func (mock *MockAccounts) Get(ctx context.Context) (accounts.Account, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx)
	}
	return mock.wrap.Accounts().Get(ctx)
}

// Actions

type MockActions struct {
	wrap   cloud.Client
	GetFn  func(ctx context.Context, id int) (actions.Action, error)
	ListFn func(ctx context.Context) (<-chan actions.Action, <-chan error)
}

func (mock *MockActions) Get(ctx context.Context, id int) (actions.Action, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx, id)
	}
	return mock.wrap.Actions().Get(ctx, id)
}

func (mock *MockActions) List(ctx context.Context) (<-chan actions.Action, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Actions().List(ctx)
}

// Domains

type MockDomains struct {
	wrap           cloud.Client
	CreateFn       func(ctx context.Context, name, ip string, opts ...domains.CreateOpt) (domains.Domain, error)
	GetFn          func(ctx context.Context, id string) (domains.Domain, error)
	DeleteFn       func(ctx context.Context, id string) error
	ListFn         func(ctx context.Context) (<-chan domains.Domain, <-chan error)
	CreateRecordFn func(ctx context.Context, id string, opts ...domains.RecordOpt) (domains.Record, error)
	GetRecordFn    func(ctx context.Context, name string, id int) (domains.Record, error)
	UpdateRecordFn func(ctx context.Context, name string, id int, opts ...domains.RecordOpt) (domains.Record, error)
	DeleteRecordFn func(ctx context.Context, name string, id int) error
	ListRecordFn   func(ctx context.Context, name string) (<-chan domains.Record, <-chan error)
}

func (mock *MockDomains) Create(ctx context.Context, name, ip string, opts ...domains.CreateOpt) (domains.Domain, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, name, ip, opts...)
	}
	return mock.wrap.Domains().Create(ctx, name, ip, opts...)
}

func (mock *MockDomains) Get(ctx context.Context, id string) (domains.Domain, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx, id)
	}
	return mock.wrap.Domains().Get(ctx, id)
}

func (mock *MockDomains) Delete(ctx context.Context, id string) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, id)
	}
	return mock.wrap.Domains().Delete(ctx, id)
}

func (mock *MockDomains) List(ctx context.Context) (<-chan domains.Domain, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Domains().List(ctx)
}

func (mock *MockDomains) CreateRecord(ctx context.Context, id string, opts ...domains.RecordOpt) (domains.Record, error) {
	if mock.CreateRecordFn != nil {
		return mock.CreateRecordFn(ctx, id, opts...)
	}
	return mock.wrap.Domains().CreateRecord(ctx, id, opts...)
}

func (mock *MockDomains) GetRecord(ctx context.Context, name string, id int) (domains.Record, error) {
	if mock.GetRecordFn != nil {
		return mock.GetRecordFn(ctx, name, id)
	}
	return mock.wrap.Domains().GetRecord(ctx, name, id)
}

func (mock *MockDomains) UpdateRecord(ctx context.Context, name string, id int, opts ...domains.RecordOpt) (domains.Record, error) {
	if mock.UpdateRecordFn != nil {
		return mock.UpdateRecordFn(ctx, name, id, opts...)
	}
	return mock.wrap.Domains().UpdateRecord(ctx, name, id, opts...)
}

func (mock *MockDomains) DeleteRecord(ctx context.Context, name string, id int) error {
	if mock.DeleteRecordFn != nil {
		return mock.DeleteRecordFn(ctx, name, id)
	}
	return mock.wrap.Domains().DeleteRecord(ctx, name, id)
}

func (mock *MockDomains) ListRecord(ctx context.Context, name string) (<-chan domains.Record, <-chan error) {
	if mock.ListRecordFn != nil {
		return mock.ListRecordFn(ctx, name)
	}
	return mock.wrap.Domains().ListRecord(ctx, name)
}

// Images

type MockImages struct {
	wrap               cloud.Client
	GetByIDFn          func(context.Context, int) (images.Image, error)
	GetBySlugFn        func(context.Context, string) (images.Image, error)
	UpdateFn           func(context.Context, int, ...images.UpdateOpt) (images.Image, error)
	DeleteFn           func(context.Context, int) error
	ListFn             func(context.Context) (<-chan images.Image, <-chan error)
	ListApplicationFn  func(context.Context) (<-chan images.Image, <-chan error)
	ListDistributionFn func(context.Context) (<-chan images.Image, <-chan error)
	ListUserFn         func(context.Context) (<-chan images.Image, <-chan error)
}

func (mock *MockImages) GetByID(ctx context.Context, id int) (images.Image, error) {
	if mock.GetByIDFn != nil {
		return mock.GetByIDFn(ctx, id)
	}
	return mock.wrap.Images().GetByID(ctx, id)
}
func (mock *MockImages) GetBySlug(ctx context.Context, slug string) (images.Image, error) {
	if mock.GetBySlugFn != nil {
		return mock.GetBySlugFn(ctx, slug)
	}
	return mock.wrap.Images().GetBySlug(ctx, slug)
}
func (mock *MockImages) Update(ctx context.Context, id int, opts ...images.UpdateOpt) (images.Image, error) {
	if mock.UpdateFn != nil {
		return mock.UpdateFn(ctx, id, opts...)
	}
	return mock.wrap.Images().Update(ctx, id, opts...)
}
func (mock *MockImages) Delete(ctx context.Context, id int) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, id)
	}
	return mock.wrap.Images().Delete(ctx, id)
}
func (mock *MockImages) List(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Images().List(ctx)
}
func (mock *MockImages) ListApplication(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListApplicationFn != nil {
		return mock.ListApplicationFn(ctx)
	}
	return mock.wrap.Images().ListApplication(ctx)
}
func (mock *MockImages) ListDistribution(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListDistributionFn != nil {
		return mock.ListDistributionFn(ctx)
	}
	return mock.wrap.Images().ListDistribution(ctx)
}
func (mock *MockImages) ListUser(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListUserFn != nil {
		return mock.ListUserFn(ctx)
	}
	return mock.wrap.Images().ListUser(ctx)
}

// Keys

type MockKeys struct {
	wrap                  cloud.Client
	CreateFn              func(ctx context.Context, name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error)
	GetByIDFn             func(context.Context, int) (keys.Key, error)
	GetByFingerprintFn    func(context.Context, string) (keys.Key, error)
	UpdateByIDFn          func(context.Context, int, ...keys.UpdateOpt) (keys.Key, error)
	UpdateByFingerprintFn func(context.Context, string, ...keys.UpdateOpt) (keys.Key, error)
	DeleteByIDFn          func(context.Context, int) error
	DeleteByFingerprintFn func(context.Context, string) error
	ListFn                func(context.Context) (<-chan keys.Key, <-chan error)
}

func (mock *MockKeys) Create(ctx context.Context, name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, name, publicKey, opts...)
	}
	return mock.wrap.Keys().Create(ctx, name, publicKey, opts...)
}

func (mock *MockKeys) GetByID(ctx context.Context, id int) (keys.Key, error) {
	if mock.GetByIDFn != nil {
		return mock.GetByIDFn(ctx, id)
	}
	return mock.wrap.Keys().GetByID(ctx, id)
}

func (mock *MockKeys) GetByFingerprint(ctx context.Context, fp string) (keys.Key, error) {
	if mock.GetByFingerprintFn != nil {
		return mock.GetByFingerprintFn(ctx, fp)
	}
	return mock.wrap.Keys().GetByFingerprint(ctx, fp)
}

func (mock *MockKeys) UpdateByID(ctx context.Context, id int, opts ...keys.UpdateOpt) (keys.Key, error) {
	if mock.UpdateByIDFn != nil {
		return mock.UpdateByIDFn(ctx, id, opts...)
	}
	return mock.wrap.Keys().UpdateByID(ctx, id, opts...)
}

func (mock *MockKeys) UpdateByFingerprint(ctx context.Context, fp string, opts ...keys.UpdateOpt) (keys.Key, error) {
	if mock.UpdateByFingerprintFn != nil {
		return mock.UpdateByFingerprintFn(ctx, fp, opts...)
	}
	return mock.wrap.Keys().UpdateByFingerprint(ctx, fp, opts...)
}

func (mock *MockKeys) DeleteByID(ctx context.Context, id int) error {
	if mock.DeleteByIDFn != nil {
		return mock.DeleteByIDFn(ctx, id)
	}
	return mock.wrap.Keys().DeleteByID(ctx, id)
}

func (mock *MockKeys) DeleteByFingerprint(ctx context.Context, fp string) error {
	if mock.DeleteByFingerprintFn != nil {
		return mock.DeleteByFingerprintFn(ctx, fp)
	}
	return mock.wrap.Keys().DeleteByFingerprint(ctx, fp)
}

func (mock *MockKeys) List(ctx context.Context) (<-chan keys.Key, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Keys().List(ctx)
}

// Regions

type MockRegions struct {
	wrap   cloud.Client
	ListFn func(ctx context.Context) (<-chan regions.Region, <-chan error)
}

func (mock *MockRegions) List(ctx context.Context) (<-chan regions.Region, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Regions().List(ctx)
}

// Sizes

type MockSizes struct {
	wrap   cloud.Client
	ListFn func(ctx context.Context) (<-chan sizes.Size, <-chan error)
}

func (mock *MockSizes) List(ctx context.Context) (<-chan sizes.Size, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Sizes().List(ctx)
}

// FloatingIPs

type MockFloatingIPs struct {
	wrap                  cloud.Client
	MockFloatingIPActions *MockFloatingIPActions
	CreateFn              func(ctx context.Context, region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error)
	GetFn                 func(ctx context.Context, ip string) (floatingips.FloatingIP, error)
	DeleteFn              func(ctx context.Context, ip string) error
	ListFn                func(ctx context.Context) (<-chan floatingips.FloatingIP, <-chan error)
}

func (mock *MockFloatingIPs) Create(ctx context.Context, region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, region, opts...)
	}
	return mock.wrap.FloatingIPs().Create(ctx, region, opts...)
}
func (mock *MockFloatingIPs) Get(ctx context.Context, ip string) (floatingips.FloatingIP, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx, ip)
	}
	return mock.wrap.FloatingIPs().Get(ctx, ip)
}
func (mock *MockFloatingIPs) Delete(ctx context.Context, ip string) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, ip)
	}
	return mock.wrap.FloatingIPs().Delete(ctx, ip)
}
func (mock *MockFloatingIPs) List(ctx context.Context) (<-chan floatingips.FloatingIP, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.FloatingIPs().List(ctx)
}
func (mock *MockFloatingIPs) Actions() floatingips.ActionClient {
	if mock.MockFloatingIPActions != nil {
		return mock.MockFloatingIPActions
	}
	return mock.wrap.FloatingIPs().Actions()
}

// FloatingIP Actions

type MockFloatingIPActions struct {
	wrap       cloud.Client
	AssignFn   func(ctx context.Context, ip string, did int) error
	UnassignFn func(ctx context.Context, ip string) error
}

func (mock *MockFloatingIPActions) Assign(ctx context.Context, ip string, did int) error {
	if mock.AssignFn != nil {
		return mock.AssignFn(ctx, ip, did)
	}
	return mock.wrap.FloatingIPs().Actions().Assign(ctx, ip, did)
}
func (mock *MockFloatingIPActions) Unassign(ctx context.Context, ip string) error {
	if mock.UnassignFn != nil {
		return mock.UnassignFn(ctx, ip)
	}
	return mock.wrap.FloatingIPs().Actions().Unassign(ctx, ip)
}

// Volumes

type MockVolumes struct {
	wrap              cloud.Client
	MockVolumeActions *MockVolumeActions
	CreateVolumeFn    func(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...volumes.CreateOpt) (volumes.Volume, error)
	GetVolumeFn       func(context.Context, string) (volumes.Volume, error)
	DeleteVolumeFn    func(context.Context, string) error
	ListVolumesFn     func(context.Context) (<-chan volumes.Volume, <-chan error)
	CreateSnapshotFn  func(ctx context.Context, volumeID, name string, opts ...volumes.SnapshotOpt) (volumes.Snapshot, error)
	GetSnapshotFn     func(context.Context, string) (volumes.Snapshot, error)
	DeleteSnapshotFn  func(context.Context, string) error
	ListSnapshotsFn   func(ctx context.Context, volumeID string) (<-chan volumes.Snapshot, <-chan error)
}

func (mock *MockVolumes) CreateVolume(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...volumes.CreateOpt) (volumes.Volume, error) {
	if mock.CreateVolumeFn != nil {
		return mock.CreateVolumeFn(ctx, name, region, sizeGibiBytes, opts...)
	}
	return mock.wrap.Volumes().CreateVolume(ctx, name, region, sizeGibiBytes, opts...)
}
func (mock *MockVolumes) GetVolume(ctx context.Context, id string) (volumes.Volume, error) {
	if mock.GetVolumeFn != nil {
		return mock.GetVolumeFn(ctx, id)
	}
	return mock.wrap.Volumes().GetVolume(ctx, id)
}
func (mock *MockVolumes) DeleteVolume(ctx context.Context, id string) error {
	if mock.DeleteVolumeFn != nil {
		return mock.DeleteVolumeFn(ctx, id)
	}
	return mock.wrap.Volumes().DeleteVolume(ctx, id)
}
func (mock *MockVolumes) ListVolumes(ctx context.Context) (<-chan volumes.Volume, <-chan error) {
	if mock.ListVolumesFn != nil {
		return mock.ListVolumesFn(ctx)
	}
	return mock.wrap.Volumes().ListVolumes(ctx)
}
func (mock *MockVolumes) CreateSnapshot(ctx context.Context, volumeID, name string, opts ...volumes.SnapshotOpt) (volumes.Snapshot, error) {
	if mock.CreateSnapshotFn != nil {
		return mock.CreateSnapshotFn(ctx, volumeID, name, opts...)
	}
	return mock.wrap.Volumes().CreateSnapshot(ctx, volumeID, name, opts...)
}
func (mock *MockVolumes) GetSnapshot(ctx context.Context, id string) (volumes.Snapshot, error) {
	if mock.GetSnapshotFn != nil {
		return mock.GetSnapshotFn(ctx, id)
	}
	return mock.wrap.Volumes().GetSnapshot(ctx, id)
}
func (mock *MockVolumes) DeleteSnapshot(ctx context.Context, id string) error {
	if mock.DeleteSnapshotFn != nil {
		return mock.DeleteSnapshotFn(ctx, id)
	}
	return mock.wrap.Volumes().DeleteSnapshot(ctx, id)
}
func (mock *MockVolumes) ListSnapshots(ctx context.Context, volumeID string) (<-chan volumes.Snapshot, <-chan error) {
	if mock.ListSnapshotsFn != nil {
		return mock.ListSnapshotsFn(ctx, volumeID)
	}
	return mock.wrap.Volumes().ListSnapshots(ctx, volumeID)
}

func (mock *MockVolumes) Actions() volumes.ActionClient {
	if mock.MockVolumeActions != nil {
		return mock.MockVolumeActions
	}
	return mock.wrap.Volumes().Actions()
}

// Volume Actions

type MockVolumeActions struct {
	wrap     cloud.Client
	AttachFn func(ctx context.Context, volumeID string, dropletID int) error
	DetachFn func(ctx context.Context, volumeID string) error
}

func (mock *MockVolumeActions) Attach(ctx context.Context, volumeID string, dropletID int) error {
	if mock.AttachFn != nil {
		return mock.AttachFn(ctx, volumeID, dropletID)
	}
	return mock.wrap.Volumes().Actions().Attach(ctx, volumeID, dropletID)
}

func (mock *MockVolumeActions) Detach(ctx context.Context, volumeID string) error {
	if mock.DetachFn != nil {
		return mock.DetachFn(ctx, volumeID)
	}
	return mock.wrap.Volumes().Actions().Detach(ctx, volumeID)
}

// Tags

type MockTags struct {
	wrap     cloud.Client
	CreateFn func(ctx context.Context, name string, opt ...tags.CreateOpt) (tags.Tag, error)
	GetFn    func(ctx context.Context, name string) (tags.Tag, error)
	ListFn   func(ctx context.Context) (<-chan tags.Tag, <-chan error)
	DeleteFn func(ctx context.Context, name string) error
	TagFn    func(ctx context.Context, name string, res []godo.Resource) error
	UntagFn  func(ctx context.Context, name string, res []godo.Resource) error
}

func (mock *MockTags) Create(ctx context.Context, name string, opts ...tags.CreateOpt) (tags.Tag, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, name, opts...)
	}

	return mock.wrap.Tags().Create(ctx, name, opts...)
}

func (mock *MockTags) Get(ctx context.Context, name string) (tags.Tag, error) {
	if mock.GetFn != nil {
		return mock.GetFn(ctx, name)
	}

	return mock.wrap.Tags().Get(ctx, name)
}

func (mock *MockTags) List(ctx context.Context) (<-chan tags.Tag, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}

	return mock.wrap.Tags().List(ctx)
}

func (mock *MockTags) Delete(ctx context.Context, name string) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, name)
	}

	return mock.wrap.Tags().Delete(ctx, name)
}

func (mock *MockTags) TagResources(ctx context.Context, name string, res []godo.Resource) error {
	if mock.TagFn != nil {
		return mock.TagFn(ctx, name, res)
	}

	return mock.wrap.Tags().TagResources(ctx, name, res)
}

func (mock *MockTags) UntagResources(ctx context.Context, name string, res []godo.Resource) error {
	if mock.UntagFn != nil {
		return mock.UntagFn(ctx, name, res)
	}

	return mock.wrap.Tags().UntagResources(ctx, name, res)
}

// Load Balancers

type MockLoadBalancers struct {
	wrap     cloud.Client
	CreateFn func(ctx context.Context, name, region string, forwardingRules []godo.ForwardingRule, opt ...loadbalancers.CreateOpt) (loadbalancers.LoadBalancer, error)
	DeleteFn func(ctx context.Context, id string) error
	ListFn   func(ctx context.Context) (<-chan loadbalancers.LoadBalancer, <-chan error)
}

func (mock *MockLoadBalancers) Create(ctx context.Context, name, region string, forwardingRules []godo.ForwardingRule, opts ...loadbalancers.CreateOpt) (loadbalancers.LoadBalancer, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(ctx, name, region, forwardingRules, opts...)
	}

	return mock.wrap.LoadBalancers().Create(ctx, name, region, forwardingRules, opts...)
}

func (mock *MockLoadBalancers) Delete(ctx context.Context, id string) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(ctx, id)
	}

	return mock.wrap.LoadBalancers().Delete(ctx, id)
}

func (mock *MockLoadBalancers) List(ctx context.Context) (<-chan loadbalancers.LoadBalancer, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}

	return mock.wrap.LoadBalancers().List(ctx)
}
