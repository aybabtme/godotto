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
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/accounts"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/actions"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/domains"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/drives"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/floatingips"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/images"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/keys"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/regions"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/sizes"
	"golang.org/x/net/context"
)

// Client

type Mock struct {
	wrap            cloud.Client
	MockDroplets    *MockDroplets
	MockAccounts    *MockAccounts
	MockActions     *MockActions
	MockDomains     *MockDomains
	MockImages      *MockImages
	MockKeys        *MockKeys
	MockRegions     *MockRegions
	MockSizes       *MockSizes
	MockFloatingIPs *MockFloatingIPs
	MockDrives      *MockDrives
}

func Client(client cloud.Client) *Mock {
	return &Mock{wrap: client,
		MockDroplets:    &MockDroplets{wrap: client},
		MockAccounts:    &MockAccounts{wrap: client},
		MockActions:     &MockActions{wrap: client},
		MockDomains:     &MockDomains{wrap: client},
		MockImages:      &MockImages{wrap: client},
		MockKeys:        &MockKeys{wrap: client},
		MockRegions:     &MockRegions{wrap: client},
		MockSizes:       &MockSizes{wrap: client},
		MockFloatingIPs: &MockFloatingIPs{wrap: client},
		MockDrives:      &MockDrives{wrap: client},
	}
}

func (mock *Mock) Droplets() droplets.Client       { return mock.MockDroplets }
func (mock *Mock) Accounts() accounts.Client       { return mock.MockAccounts }
func (mock *Mock) Actions() actions.Client         { return mock.MockActions }
func (mock *Mock) Domains() domains.Client         { return mock.MockDomains }
func (mock *Mock) Images() images.Client           { return mock.MockImages }
func (mock *Mock) Keys() keys.Client               { return mock.MockKeys }
func (mock *Mock) Regions() regions.Client         { return mock.MockRegions }
func (mock *Mock) Sizes() sizes.Client             { return mock.MockSizes }
func (mock *Mock) FloatingIPs() floatingips.Client { return mock.MockFloatingIPs }
func (mock *Mock) Drives() drives.Client           { return mock.MockDrives }

// Droplets

type MockDroplets struct {
	wrap     cloud.Client
	CreateFn func(ctx context.Context, name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error)
	GetFn    func(ctx context.Context, id int) (droplets.Droplet, error)
	DeleteFn func(ctx context.Context, id int) error
	ListFn   func(ctx context.Context) (<-chan droplets.Droplet, <-chan error)
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
	wrap     cloud.Client
	CreateFn func(ctx context.Context, region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error)
	GetFn    func(ctx context.Context, ip string) (floatingips.FloatingIP, error)
	DeleteFn func(ctx context.Context, ip string) error
	ListFn   func(ctx context.Context) (<-chan floatingips.FloatingIP, <-chan error)
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

// Drives

type MockDrives struct {
	wrap             cloud.Client
	CreateDriveFn    func(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...drives.CreateOpt) (drives.Drive, error)
	GetDriveFn       func(context.Context, string) (drives.Drive, error)
	DeleteDriveFn    func(context.Context, string) error
	ListDrivesFn     func(context.Context) (<-chan drives.Drive, <-chan error)
	CreateSnapshotFn func(ctx context.Context, driveID, name string, opts ...drives.SnapshotOpt) (drives.Snapshot, error)
	GetSnapshotFn    func(context.Context, string) (drives.Snapshot, error)
	DeleteSnapshotFn func(context.Context, string) error
	ListSnapshotsFn  func(ctx context.Context, driveID string) (<-chan drives.Snapshot, <-chan error)
}

func (mock *MockDrives) CreateDrive(ctx context.Context, name, region string, sizeGibiBytes int64, opts ...drives.CreateOpt) (drives.Drive, error) {
	if mock.CreateDriveFn != nil {
		return mock.CreateDriveFn(ctx, name, region, sizeGibiBytes, opts...)
	}
	return mock.wrap.Drives().CreateDrive(ctx, name, region, sizeGibiBytes, opts...)
}
func (mock *MockDrives) GetDrive(ctx context.Context, id string) (drives.Drive, error) {
	if mock.GetDriveFn != nil {
		return mock.GetDriveFn(ctx, id)
	}
	return mock.wrap.Drives().GetDrive(ctx, id)
}
func (mock *MockDrives) DeleteDrive(ctx context.Context, id string) error {
	if mock.DeleteDriveFn != nil {
		return mock.DeleteDriveFn(ctx, id)
	}
	return mock.wrap.Drives().DeleteDrive(ctx, id)
}
func (mock *MockDrives) ListDrives(ctx context.Context) (<-chan drives.Drive, <-chan error) {
	if mock.ListDrivesFn != nil {
		return mock.ListDrivesFn(ctx)
	}
	return mock.wrap.Drives().ListDrives(ctx)
}
func (mock *MockDrives) CreateSnapshot(ctx context.Context, driveID, name string, opts ...drives.SnapshotOpt) (drives.Snapshot, error) {
	if mock.CreateSnapshotFn != nil {
		return mock.CreateSnapshotFn(ctx, driveID, name, opts...)
	}
	return mock.wrap.Drives().CreateSnapshot(ctx, driveID, name, opts...)
}
func (mock *MockDrives) GetSnapshot(ctx context.Context, id string) (drives.Snapshot, error) {
	if mock.GetSnapshotFn != nil {
		return mock.GetSnapshotFn(ctx, id)
	}
	return mock.wrap.Drives().GetSnapshot(ctx, id)
}
func (mock *MockDrives) DeleteSnapshot(ctx context.Context, id string) error {
	if mock.DeleteSnapshotFn != nil {
		return mock.DeleteSnapshotFn(ctx, id)
	}
	return mock.wrap.Drives().DeleteSnapshot(ctx, id)
}
func (mock *MockDrives) ListSnapshots(ctx context.Context, driveID string) (<-chan drives.Snapshot, <-chan error) {
	if mock.ListSnapshotsFn != nil {
		return mock.ListSnapshotsFn(ctx, driveID)
	}
	return mock.wrap.Drives().ListSnapshots(ctx, driveID)
}
