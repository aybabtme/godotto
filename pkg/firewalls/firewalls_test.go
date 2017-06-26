package firewalls_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/firewalls"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

var testFw *godo.Firewall = &godo.Firewall{}

func TestFirewallApply(t *testing.T) {
	cloud := mockcloud.Client(nil)

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
assert(pkg != null, "package should be loaded");
assert(pkg.create != null, "create function should be defined");
assert(pkg.get != null, "get function should be defined");
assert(pkg.delete != null, "delete function should be defined");
assert(pkg.update != null, "update function should be defined");
assert(pkg.list != null, "list function should be defined");
assert(pkg.add_tags != null, "add_tags function should be defined");
assert(pkg.remove_tags != null, "remove_tags function should be defined");
assert(pkg.add_droplets != null, "add_droplets function should be defined");
assert(pkg.remove_droplets != null, "remove_droplets function should be defined");
assert(pkg.add_rules != null, "add_rules function should be defined");
assert(pkg.remove_rules != null, "remove_rules function should be defined");
`)
}

func TestFirewallThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockFirewalls.CreateFn = func(_ context.Context, _ string, rules []godo.InboundRule, _ []godo.OutboundRule, _ ...firewalls.CreateOpt) (firewalls.Firewall, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockFirewalls.GetFn = func(_ context.Context, _ string) (firewalls.Firewall, error) {
		return nil, errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
	 var pkg = cloud.firewalls;
	 var fw = {
		"id": "84d87802-df17-4f8f-a691-58e408570c12",
		"name": "test-sg",
		"inbound_rules": [
		  {
			"protocol": "icmp",
			"ports": "0",
			"sources": {
			  "load_balancer_uids": [
				"d2d3920a-9d45-41b0-b018-d15e18ec60a4"
			  ],
			  "tags": [
				"haproxy"
			  ]
			}
		  },
		  {
			"protocol": "tcp",
			"ports": "8000-9000",
			"sources": {
			  "addresses": [
				"0.0.0.0/0"
			  ]
			}
		  }
		],
		"outbound_rules": [
		  {
			"protocol": "icmp",
			"ports": "0",
			"destinations": {
			  "tags": [
				"haproxy"
			  ]
			}
		  },
		  {
			"protocol": "tcp",
			"ports": "8000-9000",
			"destinations": {
			  "addresses": [
				"::/1"
			  ]
			}
		  }
		],
		"created_at": "2017-05-08T16:25:21Z",
		"droplet_ids": [
		  46298047
		],
		"tags": [
		  "haproxy"
		]
	  };

	var dropletId = 42;

[
  {name: "create", fn: function() { pkg.create(fw) }},
  {name: "get", fn: function() { pkg.get(fw.id) }},
].forEach(function(kv) {
var name = kv.name;
var fn = kv.fn;
try {
				fn(); throw "don't catch me";
			} catch(e) {
				equals("throw me", e.message, name + " should send the right exception!");
			}
});

`)
}
