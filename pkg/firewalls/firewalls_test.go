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

	cloud.MockFirewalls.DeleteFn = func(_ context.Context, _ string) error {
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
	f = &godo.Firewall{
		ID:   "test-uuid",
		Name: "test-sg",
		InboundRules: []godo.InboundRule{
			{
				Protocol:  "icmp",
				PortRange: "0",
				Sources: &godo.Sources{
					LoadBalancerUIDs: []string{"test-lb-uuid"},
					Tags:             []string{"haproxy"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8000-9000",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0"},
				},
			},
		},
		OutboundRules: []godo.OutboundRule{
			{
				Protocol:  "icmp",
				PortRange: "0",
				Destinations: &godo.Destinations{
					Tags: []string{"haproxy"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8000-9000",
				Destinations: &godo.Destinations{
					Addresses: []string{"::/1"},
				},
			},
		},
		DropletIDs: []int{123456},
		Tags:       []string{"haproxy"},
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
