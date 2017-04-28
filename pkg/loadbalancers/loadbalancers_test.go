package loadbalancers

import (
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
)

func TestTagApply(t *testing.T) {
	cloud := mockcloud.Client(nil)

	vmtest.Run(t, cloud, `
	var pkg = cloud.load_balancers;
	assert(pkg != null, "package should be loaded");
	assert(pkg.create != null, "create function should be defined");
	`)
}
