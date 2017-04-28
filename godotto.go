package godotto

import (
	"context"
	"fmt"

	"github.com/aybabtme/godotto/pkg/accounts"
	"github.com/aybabtme/godotto/pkg/actions"
	"github.com/aybabtme/godotto/pkg/domains"
	"github.com/aybabtme/godotto/pkg/droplets"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/floatingips"
	"github.com/aybabtme/godotto/pkg/images"
	"github.com/aybabtme/godotto/pkg/keys"
	"github.com/aybabtme/godotto/pkg/loadbalancers"
	"github.com/aybabtme/godotto/pkg/regions"
	"github.com/aybabtme/godotto/pkg/sizes"
	"github.com/aybabtme/godotto/pkg/tags"
	"github.com/aybabtme/godotto/pkg/volumes"
	"github.com/robertkrimen/otto"
)

var q = otto.NullValue()

func Apply(ctx context.Context, vm *otto.Otto, client cloud.Client) (otto.Value, error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	for _, applier := range []struct {
		Name  string
		Apply func(context.Context, *otto.Otto, cloud.Client) (otto.Value, error)
	}{
		{"accounts", accounts.Apply},
		{"actions", actions.Apply},
		{"domains", domains.Apply},
		{"droplets", droplets.Apply},
		{"images", images.Apply},
		{"keys", keys.Apply},
		{"regions", regions.Apply},
		{"floating_ips", floatingips.Apply},
		{"sizes", sizes.Apply},
		{"volumes", volumes.Apply},
		{"tags", tags.Apply},
		{"load_balancers", loadbalancers.Apply},
	} {
		svc, err := applier.Apply(ctx, vm, client)
		if err != nil {
			return q, fmt.Errorf("preparing godo %s service: %v", applier.Name, err)
		}
		if err := root.Set(applier.Name, svc); err != nil {
			return q, fmt.Errorf("adding godo %s service: %v", applier.Name, err)
		}
	}
	return root.Value(), err
}
