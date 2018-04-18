package accounts_test

import (
	"errors"
	"testing"
	"context"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/accounts"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

func TestApply(t *testing.T) {
	cloud := mockcloud.Client(nil)
	vmtest.Run(t, cloud, `
var pkg = cloud.accounts;

assert(pkg != null, "package should be loaded");
assert(pkg.get != null, "get function should be defined");
    `)
}

type account struct {
	*godo.Account
}

func (k *account) Struct() *godo.Account { return k.Account }

func TestGet(t *testing.T) {
	cloud := mockcloud.Client(nil)

	want := &godo.Account{
		DropletLimit:    9000,
		FloatingIPLimit: 9000,
		Email:           "hello@example.com",
		UUID:            "deadbeef",
		EmailVerified:   true,
		Status:          "it works",
		StatusMessage:   "no really",
	}

	cloud.MockAccounts.GetFn = func(_ context.Context) (accounts.Account, error) {
		return &account{want}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.accounts;

var a = pkg.get();
var want = {
	droplet_limit: 9000,
	floating_ip_limit: 9000,
	email: "hello@example.com",
	uuid: "deadbeef",
	email_verified: true,
	status: "it works",
	status_message: "no really"
};
equals(a, want, "should get proper object");
    `)
}

func TestGetThrow(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockAccounts.GetFn = func(_ context.Context) (accounts.Account, error) {
		return nil, errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.accounts;

try {
	pkg.get();
	throw "dont catch me";
} catch (e) {
	equals("throw me", e.message, "should send the right exception");
}
    `)
}
