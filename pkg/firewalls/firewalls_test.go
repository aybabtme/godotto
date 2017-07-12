package firewalls_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/firewalls"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
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
	assert(pkg.list != null, "list function should be defined");
	assert(pkg.update != null, "update function should be defined");
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

	cloud.MockFirewalls.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	cloud.MockFirewalls.ListFn = func(_ context.Context) (<-chan firewalls.Firewall, <-chan error) {
		fc := make(chan firewalls.Firewall)
		close(fc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return fc, ec
	}

	cloud.MockFirewalls.UpdateFn = func(_ context.Context, _ string, _ ...firewalls.UpdateOpt) (firewalls.Firewall, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockFirewalls.AddTagsFn = func(_ context.Context, _ string, _ ...string) error {
		return errors.New("throw me")
	}

	cloud.MockFirewalls.RemoveTagsFn = func(_ context.Context, _ string, _ ...string) error {
		return errors.New("throw me")
	}

	cloud.MockFirewalls.AddDropletsFn = func(_ context.Context, _ string, _ ...int) error {
		return errors.New("throw me")
	}

	cloud.MockFirewalls.RemoveDropletsFn = func(_ context.Context, _ string, _ ...int) error {
		return errors.New("throw me")
	}

	cloud.MockFirewalls.AddRulesFn = func(_ context.Context, _ string, _ []godo.InboundRule, _ []godo.OutboundRule) error {
		return errors.New("throw me")

	}

	cloud.MockFirewalls.RemoveRulesFn = func(_ context.Context, _ string, _ []godo.InboundRule, _ []godo.OutboundRule) error {
		return errors.New("throw me")

	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.firewalls;
	var fw = {
		"id": "test-uuid",
		"name": "test-sg",
		"inbound_rules": [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
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
		123456
		],
		"tags": [
		"haproxy"
		]
	};

	var dropletId = 42;

	[
	{name: "create", fn: function() { pkg.create(fw) }},
	{name: "get", fn: function() { pkg.get(fw.id) }},
	{name: "list", fn: function() { pkg.list() }},
	{name: "delete", fn: function() { pkg.delete(fw.id) }},
	{name: "update", fn: function() { pkg.update(fw.id, fw) }},
	{name: "add_tags", fn: function() { pkg.add_tags(fw.id, ["test"]) }},
	{name: "remove_tags", fn: function() { pkg.remove_tags(fw.id, ["test"]) }},
	{name: "add_droplets", fn: function() { pkg.add_droplets(fw.id, [42]) }},
	{name: "remove_droplets", fn: function() { pkg.remove_droplets(fw.id, [42]) }},
	{name: "add_rules", fn: function() { pkg.add_rules(fw.id, [], []) }},
	{name: "remove_rules", fn: function() { pkg.remove_rules(fw.id, [], []) }},
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

type firewall struct {
	*godo.Firewall
}

func (k *firewall) Struct() *godo.Firewall { return k.Firewall }

var (
	fwInboundRules = []godo.InboundRule{
		{
			Protocol:  "icmp",
			PortRange: "0",
			Sources: &godo.Sources{
				LoadBalancerUIDs: []string{"test-lb-uuid"},
				Tags:             []string{"haproxy"},
				Addresses:        []string(nil),
				DropletIDs:       []int{},
			},
		},
		{
			Protocol:  "tcp",
			PortRange: "8000-9000",
			Sources: &godo.Sources{
				Addresses:        []string{"0.0.0.0/0"},
				Tags:             []string{},
				DropletIDs:       []int{},
				LoadBalancerUIDs: []string{},
			},
		},
	}

	fwOutboundRules = []godo.OutboundRule{
		{
			Protocol:  "icmp",
			PortRange: "0",
			Destinations: &godo.Destinations{
				Tags:             []string{"haproxy"},
				Addresses:        []string(nil),
				DropletIDs:       []int{},
				LoadBalancerUIDs: []string{},
			},
		},
		{
			Protocol:  "tcp",
			PortRange: "8000-9000",
			Destinations: &godo.Destinations{
				Addresses:        []string{"::/1"},
				DropletIDs:       []int{},
				LoadBalancerUIDs: []string{},
				Tags:             []string{},
			},
		},
	}

	f = &godo.Firewall{
		ID:            "test-uuid",
		Name:          "test-sg",
		InboundRules:  fwInboundRules,
		OutboundRules: fwOutboundRules,
		DropletIDs:    []int{123456},
		Tags:          []string{"haproxy"},
	}
)

func TestFirewallCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockFirewalls.CreateFn = func(_ context.Context, name string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule, opts ...firewalls.CreateOpt) (firewalls.Firewall, error) {
		return &firewall{f}, nil
	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.firewalls;

	var f = pkg.create(
		{
			"id": "test-uuid",
			"name": "test-sg",
			"inbound_rules": [
			{
				"protocol": "icmp",
				"ports": "0",
				"sources": {
					"load_balancer_uids": [
					"test-lb-uuid",
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
					],
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
			123456
			],
			"tags": [
			"haproxy"
			]
		}
	);

	var want = {
		"id": "test-uuid",
		"name": "test-sg",
		"inbound_rules": [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
				],
				"tags": [
				"haproxy"
				],
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
				],
			}
		},
		{
			"protocol": "tcp",
			"ports": "8000-9000",
			"destinations": {
				"addresses": [
				"::/1"
				],
			}
		}
		],
		"created_at": "",
		"droplet_ids": [
		123456
		],
		"tags": [
		"haproxy"
		],
		"status": "",
	};

	equals(f, want, "should have proper object");
	`)
}

func TestFirewallDelete(t *testing.T) {
	wantId := "test-uuid"
	cloud := mockcloud.Client(nil)

	cloud.MockFirewalls.DeleteFn = func(_ context.Context, gotId string) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
pkg.delete("test-uuid");
`)
}

func TestFirewallList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockFirewalls.ListFn = func(_ context.Context) (<-chan firewalls.Firewall, <-chan error) {
		fc := make(chan firewalls.Firewall, 1)
		fc <- &firewall{f}
		close(fc)
		ec := make(chan error)
		close(ec)
		return fc, ec
	}
	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
var list = pkg.list();
assert(list != null, "should have received a list");
assert(list.length > 0, "should have received some elements");

var want = {
		"id": "test-uuid",
		"name": "test-sg",
		"inbound_rules": [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
				],
				"tags": [
				"haproxy"
				],
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
				],
			}
		},
		{
			"protocol": "tcp",
			"ports": "8000-9000",
			"destinations": {
				"addresses": [
				"::/1"
				],
			}
		}
		],
		"created_at": "",
		"droplet_ids": [
		123456
		],
		"tags": [
		"haproxy"
		],
		"status": "",
	};

var fw = list[0];

equals(fw, want, "should have proper object");
`)
}

func TestFirewallUpdate(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantID := "test-uuid"

	cloud.MockFirewalls.UpdateFn = func(_ context.Context, gotID string, opts ...firewalls.UpdateOpt) (firewalls.Firewall, error) {
		if gotID != wantID {
			t.Fatalf("want %v got %v", wantID, gotID)
		}

		return &firewall{f}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;

var want = {
		"id": "test-uuid",
		"name": "test-sg",
		"inbound_rules": [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
				],
				"tags": [
				"haproxy"
				],
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
				],
			}
		},
		{
			"protocol": "tcp",
			"ports": "8000-9000",
			"destinations": {
				"addresses": [
				"::/1"
				],
			}
		}
		],
		"created_at": "",
		"droplet_ids": [
		123456
		],
		"tags": [
		"haproxy"
		],
		"status": "",
	};

var f = pkg.update("test-uuid", want);
equals(f, want, "should have proper object");
`)
}

func TestFirewallGet(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockFirewalls.GetFn = func(_ context.Context, id string) (firewalls.Firewall, error) {
		return &firewall{f}, nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
var want = {
		"id": "test-uuid",
		"name": "test-sg",
		"inbound_rules": [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
				],
				"tags": [
				"haproxy"
				],
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
				],
			}
		},
		{
			"protocol": "tcp",
			"ports": "8000-9000",
			"destinations": {
				"addresses": [
				"::/1"
				],
			}
		}
		],
		"created_at": "",
		"droplet_ids": [
		123456
		],
		"tags": [
		"haproxy"
		],
		"status": "",
};

var f = pkg.get('test-uuid');

equals(f, want, "should have proper object");
`)
}

func TestFirewallAddTags(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.AddTagsFn = func(_ context.Context, gotId string, tags ...string) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantTag := "test"
		gotTag := tags[0]

		if wantTag != gotTag {
			t.Fatalf("want %v got %v", wantTag, gotTag)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
pkg.add_tags("test-uuid", ["test"]);
`)
}

func TestFirewallRemoveTags(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.RemoveTagsFn = func(_ context.Context, gotId string, tags ...string) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantTag := "test"
		gotTag := tags[0]

		if wantTag != gotTag {
			t.Fatalf("want %v got %v", wantTag, gotTag)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
pkg.remove_tags("test-uuid", ["test"]);
`)
}

func TestFirewallAddDroplets(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.AddDropletsFn = func(_ context.Context, gotId string, dids ...int) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantDropletID := 42
		gotDropletID := dids[0]

		if wantDropletID != gotDropletID {
			t.Fatalf("want %v got %v", wantDropletID, gotDropletID)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
pkg.add_droplets("test-uuid", [42]);
`)
}

func TestFirewallRemoveDroplets(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.RemoveDropletsFn = func(_ context.Context, gotId string, dids ...int) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantDropletID := 42
		gotDropletID := dids[0]

		if wantDropletID != gotDropletID {
			t.Fatalf("want %v got %v", wantDropletID, gotDropletID)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;
pkg.remove_droplets("test-uuid", [42]);
`)
}

func TestFirewallAddRules(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.AddRulesFn = func(_ context.Context, gotId string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantInboundRule := fwInboundRules[0]
		gotInboundRule := inboundRules[0]
		assert.Equal(t, wantInboundRule, gotInboundRule, "inbound rules should match")

		wantOutboundRule := fwOutboundRules[0]
		gotOutboundRule := outboundRules[0]
		assert.Equal(t, wantOutboundRule, gotOutboundRule, "outbound rules should match")

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;

var inbound_rules = [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
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
	];

var outbound_rules = [
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
	];

pkg.add_rules("test-uuid", inbound_rules, outbound_rules);
`)
}

func TestFirewallRemoveRules(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockFirewalls.RemoveRulesFn = func(_ context.Context, gotId string, inboundRules []godo.InboundRule, outboundRules []godo.OutboundRule) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantInboundRule := fwInboundRules[0]
		gotInboundRule := inboundRules[0]
		assert.Equal(t, wantInboundRule, gotInboundRule, "inbound rules should match")

		wantOutboundRule := fwOutboundRules[0]
		gotOutboundRule := outboundRules[0]
		assert.Equal(t, wantOutboundRule, gotOutboundRule, "outbound rules should match")

		return nil
	}

	vmtest.Run(t, cloud, `
var pkg = cloud.firewalls;

var inbound_rules = [
		{
			"protocol": "icmp",
			"ports": "0",
			"sources": {
				"load_balancer_uids": [
				"test-lb-uuid"
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
	];

var outbound_rules = [
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
	];

pkg.remove_rules("test-uuid", inbound_rules, outbound_rules);
`)
}
