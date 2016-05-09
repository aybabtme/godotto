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
	CreateFn func(name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error)
	GetFn    func(id int) (droplets.Droplet, error)
	DeleteFn func(id int) error
	ListFn   func(ctx context.Context) (<-chan droplets.Droplet, <-chan error)
}

func (mock *MockDroplets) Create(name, region, size, image string, opts ...droplets.CreateOpt) (droplets.Droplet, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(name, region, size, image, opts...)
	}
	return mock.wrap.Droplets().Create(name, region, size, image, opts...)
}
func (mock *MockDroplets) Get(id int) (droplets.Droplet, error) {
	if mock.GetFn != nil {
		return mock.GetFn(id)
	}
	return mock.wrap.Droplets().Get(id)
}
func (mock *MockDroplets) Delete(id int) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(id)
	}
	return mock.wrap.Droplets().Delete(id)
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
	GetFn func() (accounts.Account, error)
}

func (mock *MockAccounts) Get() (accounts.Account, error) {
	if mock.GetFn != nil {
		return mock.GetFn()
	}
	return mock.wrap.Accounts().Get()
}

// Actions

type MockActions struct {
	wrap   cloud.Client
	GetFn  func(id int) (actions.Action, error)
	ListFn func(ctx context.Context) (<-chan actions.Action, <-chan error)
}

func (mock *MockActions) Get(id int) (actions.Action, error) {
	if mock.GetFn != nil {
		return mock.GetFn(id)
	}
	return mock.wrap.Actions().Get(id)
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
	CreateFn       func(name, ip string, opts ...domains.CreateOpt) (domains.Domain, error)
	GetFn          func(id string) (domains.Domain, error)
	DeleteFn       func(id string) error
	ListFn         func(ctx context.Context) (<-chan domains.Domain, <-chan error)
	CreateRecordFn func(id string, opts ...domains.RecordOpt) (domains.Record, error)
	GetRecordFn    func(name string, id int) (domains.Record, error)
	UpdateRecordFn func(name string, id int, opts ...domains.RecordOpt) (domains.Record, error)
	DeleteRecordFn func(name string, id int) error
	ListRecordFn   func(ctx context.Context, name string) (<-chan domains.Record, <-chan error)
}

func (mock *MockDomains) Create(name, ip string, opts ...domains.CreateOpt) (domains.Domain, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(name, ip, opts...)
	}
	return mock.wrap.Domains().Create(name, ip, opts...)
}

func (mock *MockDomains) Get(id string) (domains.Domain, error) {
	if mock.GetFn != nil {
		return mock.GetFn(id)
	}
	return mock.wrap.Domains().Get(id)
}

func (mock *MockDomains) Delete(id string) error {
	if mock.DeleteFn != nil {
		return mock.DeleteFn(id)
	}
	return mock.wrap.Domains().Delete(id)
}

func (mock *MockDomains) List(ctx context.Context) (<-chan domains.Domain, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Domains().List(ctx)
}

func (mock *MockDomains) CreateRecord(id string, opts ...domains.RecordOpt) (domains.Record, error) {
	if mock.CreateRecordFn != nil {
		return mock.CreateRecordFn(id, opts...)
	}
	return mock.wrap.Domains().CreateRecord(id, opts...)
}

func (mock *MockDomains) GetRecord(name string, id int) (domains.Record, error) {
	if mock.GetRecordFn != nil {
		return mock.GetRecordFn(name, id)
	}
	return mock.wrap.Domains().GetRecord(name, id)
}

func (mock *MockDomains) UpdateRecord(name string, id int, opts ...domains.RecordOpt) (domains.Record, error) {
	if mock.UpdateRecordFn != nil {
		return mock.UpdateRecordFn(name, id, opts...)
	}
	return mock.wrap.Domains().UpdateRecord(name, id, opts...)
}

func (mock *MockDomains) DeleteRecord(name string, id int) error {
	if mock.DeleteRecordFn != nil {
		return mock.DeleteRecordFn(name, id)
	}
	return mock.wrap.Domains().DeleteRecord(name, id)
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
	GetByIDFn          func(int) (images.Image, error)
	GetBySlugFn        func(string) (images.Image, error)
	UpdateFn           func(int, ...images.UpdateOpt) (images.Image, error)
	DeleteFn           func(int) error
	ListFn             func(context.Context) (<-chan images.Image, <-chan error)
	ListApplicationFn  func(context.Context) (<-chan images.Image, <-chan error)
	ListDistributionFn func(context.Context) (<-chan images.Image, <-chan error)
	ListUserFn         func(context.Context) (<-chan images.Image, <-chan error)
}

func (mock *MockImages) GetByID(id int) (images.Image, error) {
	if mock.GetByIDFn != nil {
		return mock.GetByID(id)
	}
	return mock.wrap.Images().GetByID(id)
}
func (mock *MockImages) GetBySlug(slug string) (images.Image, error) {
	if mock.GetBySlugFn != nil {
		return mock.GetBySlug(slug)
	}
	return mock.wrap.Images().GetBySlug(slug)
}
func (mock *MockImages) Update(id int, opts ...images.UpdateOpt) (images.Image, error) {
	if mock.UpdateFn != nil {
		return mock.Update(id, opts...)
	}
	return mock.wrap.Images().Update(id, opts...)
}
func (mock *MockImages) Delete(id int) error {
	if mock.DeleteFn != nil {
		return mock.Delete(id)
	}
	return mock.wrap.Images().Delete(id)
}
func (mock *MockImages) List(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListFn != nil {
		return mock.ListFn(ctx)
	}
	return mock.wrap.Images().List(ctx)
}
func (mock *MockImages) ListApplication(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListApplicationFn != nil {
		return mock.ListApplication(ctx)
	}
	return mock.wrap.Images().ListApplication(ctx)
}
func (mock *MockImages) ListDistribution(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListDistributionFn != nil {
		return mock.ListDistribution(ctx)
	}
	return mock.wrap.Images().ListDistribution(ctx)
}
func (mock *MockImages) ListUser(ctx context.Context) (<-chan images.Image, <-chan error) {
	if mock.ListUserFn != nil {
		return mock.ListUser(ctx)
	}
	return mock.wrap.Images().ListUser(ctx)
}

// Keys

type MockKeys struct {
	wrap                  cloud.Client
	CreateFn              func(name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error)
	GetByIDFn             func(int) (keys.Key, error)
	GetByFingerprintFn    func(string) (keys.Key, error)
	UpdateByIDFn          func(int, ...keys.UpdateOpt) (keys.Key, error)
	UpdateByFingerprintFn func(string, ...keys.UpdateOpt) (keys.Key, error)
	DeleteByIDFn          func(int) error
	DeleteByFingerprintFn func(string) error
	ListFn                func(context.Context) (<-chan keys.Key, <-chan error)
}

func (mock *MockKeys) Create(name, publicKey string, opts ...keys.CreateOpt) (keys.Key, error) {
	if mock.CreateFn != nil {
		return mock.CreateFn(name, publicKey, opts...)
	}
	return mock.wrap.Keys().Create(name, publicKey, opts...)
}

func (mock *MockKeys) GetByID(id int) (keys.Key, error) {
	if mock.GetByIDFn != nil {
		return mock.GetByIDFn(id)
	}
	return mock.wrap.Keys().GetByID(id)
}

func (mock *MockKeys) GetByFingerprint(fp string) (keys.Key, error) {
	if mock.GetByFingerprintFn != nil {
		return mock.GetByFingerprintFn(fp)
	}
	return mock.wrap.Keys().GetByFingerprint(fp)
}

func (mock *MockKeys) UpdateByID(id int, opts ...keys.UpdateOpt) (keys.Key, error) {
	if mock.UpdateByIDFn != nil {
		return mock.UpdateByIDFn(id, opts...)
	}
	return mock.wrap.Keys().UpdateByID(id, opts...)
}

func (mock *MockKeys) UpdateByFingerprint(fp string, opts ...keys.UpdateOpt) (keys.Key, error) {
	if mock.UpdateByFingerprintFn != nil {
		return mock.UpdateByFingerprintFn(fp, opts...)
	}
	return mock.wrap.Keys().UpdateByFingerprint(fp, opts...)
}

func (mock *MockKeys) DeleteByID(id int) error {
	if mock.DeleteByIDFn != nil {
		return mock.DeleteByIDFn(id)
	}
	return mock.wrap.Keys().DeleteByID(id)
}

func (mock *MockKeys) DeleteByFingerprint(fp string) error {
	if mock.DeleteByFingerprintFn != nil {
		return mock.DeleteByFingerprintFn(fp)
	}
	return mock.wrap.Keys().DeleteByFingerprint(fp)
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
	CreateFn func(region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error)
	GetFn    func(ip string) (floatingips.FloatingIP, error)
	DeleteFn func(ip string) error
	ListFn   func(ctx context.Context) (<-chan floatingips.FloatingIP, <-chan error)
}

func (mock *MockFloatingIPs) Create(region string, opts ...floatingips.CreateOpt) (floatingips.FloatingIP, error) {
	if mock.CreateFn != nil {
		return mock.Create(region, opts...)
	}
	return mock.wrap.FloatingIPs().Create(region, opts...)
}
func (mock *MockFloatingIPs) Get(ip string) (floatingips.FloatingIP, error) {
	if mock.GetFn != nil {
		return mock.Get(ip)
	}
	return mock.wrap.FloatingIPs().Get(ip)
}
func (mock *MockFloatingIPs) Delete(ip string) error {
	if mock.DeleteFn != nil {
		return mock.Delete(ip)
	}
	return mock.wrap.FloatingIPs().Delete(ip)
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
	CreateDriveFn    func(name, region string, sizeGibiBytes int64, opts ...drives.CreateOpt) (drives.Drive, error)
	GetDriveFn       func(string) (drives.Drive, error)
	DeleteDriveFn    func(string) error
	ListDrivesFn     func(context.Context) (<-chan drives.Drive, <-chan error)
	CreateSnapshotFn func(driveID, name string, opts ...drives.SnapshotOpt) (drives.Snapshot, error)
	GetSnapshotFn    func(string) (drives.Snapshot, error)
	DeleteSnapshotFn func(string) error
	ListSnapshotsFn  func(ctx context.Context, driveID string) (<-chan drives.Snapshot, <-chan error)
}

func (mock *MockDrives) CreateDrive(name, region string, sizeGibiBytes int64, opts ...drives.CreateOpt) (drives.Drive, error) {
	if mock.CreateDriveFn != nil {
		return mock.CreateDriveFn(name, region, sizeGibiBytes, opts...)
	}
	return mock.wrap.Drives().CreateDrive(name, region, sizeGibiBytes, opts...)
}
func (mock *MockDrives) GetDrive(id string) (drives.Drive, error) {
	if mock.GetDriveFn != nil {
		return mock.GetDriveFn(id)
	}
	return mock.wrap.Drives().GetDrive(id)
}
func (mock *MockDrives) DeleteDrive(id string) error {
	if mock.DeleteDriveFn != nil {
		return mock.DeleteDriveFn(id)
	}
	return mock.wrap.Drives().DeleteDrive(id)
}
func (mock *MockDrives) ListDrives(ctx context.Context) (<-chan drives.Drive, <-chan error) {
	if mock.ListDrivesFn != nil {
		return mock.ListDrivesFn(ctx)
	}
	return mock.wrap.Drives().ListDrives(ctx)
}
func (mock *MockDrives) CreateSnapshot(driveID, name string, opts ...drives.SnapshotOpt) (drives.Snapshot, error) {
	if mock.CreateSnapshotFn != nil {
		return mock.CreateSnapshotFn(driveID, name, opts...)
	}
	return mock.wrap.Drives().CreateSnapshot(driveID, name, opts...)
}
func (mock *MockDrives) GetSnapshot(id string) (drives.Snapshot, error) {
	if mock.GetSnapshotFn != nil {
		return mock.GetSnapshotFn(id)
	}
	return mock.wrap.Drives().GetSnapshot(id)
}
func (mock *MockDrives) DeleteSnapshot(id string) error {
	if mock.DeleteSnapshotFn != nil {
		return mock.DeleteSnapshotFn(id)
	}
	return mock.wrap.Drives().DeleteSnapshot(id)
}
func (mock *MockDrives) ListSnapshots(ctx context.Context, driveID string) (<-chan drives.Snapshot, <-chan error) {
	if mock.ListSnapshotsFn != nil {
		return mock.ListSnapshotsFn(ctx, driveID)
	}
	return mock.wrap.Drives().ListSnapshots(ctx, driveID)
}
