package loadbalancers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud/loadbalancers"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/digitalocean/godo"
)

var testLb *godo.LoadBalancer = &godo.LoadBalancer{}

func TestLoadBalancerApply(t *testing.T) {
	cloud := mockcloud.Client(nil)

	vmtest.Run(t, cloud, `
	var pkg = cloud.load_balancers;
	assert(pkg != null, "package should be loaded");
	assert(pkg.create != null, "create function should be defined");
	assert(pkg.delete != null, "delete function should be defined");
	assert(pkg.get != null, "get function shouled be defined");
	assert(pkg.add_droplets != null, "add_droplets function should be defined");
	assert(pkg.remove_droplets != null, "remove_droplets function should be defined");
	assert(pkg.add_forwarding_rules != null, "add_forwarding_rules function should be defined");
	assert(pkg.remove_forwarding_rules != null, "remove_forwarding_rules function should be defined");
	`)
}

func TestLoadBalancerThrows(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockLoadBalancers.ListFn = func(_ context.Context) (<-chan loadbalancers.LoadBalancer, <-chan error) {
		lc := make(chan loadbalancers.LoadBalancer)
		close(lc)
		ec := make(chan error, 1)
		ec <- errors.New("throw me")
		close(ec)
		return lc, ec
	}

	cloud.MockLoadBalancers.CreateFn = func(_ context.Context, _, _ string, _ []godo.ForwardingRule, _ ...loadbalancers.CreateOpt) (loadbalancers.LoadBalancer, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockLoadBalancers.GetFn = func(_ context.Context, _ string) (loadbalancers.LoadBalancer, error) {
		return nil, errors.New("throw me")
	}

	cloud.MockLoadBalancers.DeleteFn = func(_ context.Context, _ string) error {
		return errors.New("throw me")
	}

	cloud.MockLoadBalancers.AddDropletsFn = func(_ context.Context, _ string, _ ...int) error {
		return errors.New("throw me")
	}

	cloud.MockLoadBalancers.RemoveDropletsFn = func(_ context.Context, _ string, _ ...int) error {
		return errors.New("throw me")
	}

	cloud.MockLoadBalancers.AddForwardingRulesFn = func(_ context.Context, _ string, _ ...godo.ForwardingRule) error {
		return errors.New("throw me")
	}

	cloud.MockLoadBalancers.RemoveForwardingRulesFn = func(_ context.Context, _ string, _ ...godo.ForwardingRule) error {
		return errors.New("throw me")
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;

		var lb = {
				id: "12-34-56-78",
				name: "example-lb-01",
				ip: "138.197.50.73",
				algorithm: "least_connections",
				status: "active",
				created_at: "2017-02-01T22:22:58Z",
				forwarding_rules: [
				{
					entry_protocol: "http",
					entry_port: 80,
					target_protocol: "http",
					target_port: 80,
					certificate_id: "",
					tls_passthrough: false
				},
				{
					entry_protocol: "https",
					entry_port: 444,
					target_protocol: "https",
					target_port: 443,
					certificate_id: "",
					tls_passthrough: true
				}
				],
				health_check: {
					protocol: "http",
					port: 80,
					path: "/",
					check_interval_seconds: 10,
					response_timeout_seconds: 5,
					healthy_threshold: 5,
					unhealthy_threshold: 3
				},
				sticky_sessions: {
					type: "cookies",
					cookie_name: "DO_LB",
					cookie_ttl_seconds: 300
				},
				region: {
					name: "New York 3",
					slug: "nyc3",
					sizes: [
					"512mb",
					"1gb",
					"2gb",
					"4gb",
					"8gb",
					"16gb",
					"m-16gb",
					"32gb",
					"m-32gb",
					"48gb",
					"m-64gb",
					"64gb",
					"m-128gb",
					"m-224gb"
					],
					features: [
					"private_networking",
					"backups",
					"ipv6",
					"metadata",
					"install_agent"
					],
					available: true
				},
				tag: "",
				droplet_ids: [
				34153248,
				34153250
				],
				redirect_http_to_https: false
		};

		var dropletId = 42;

		[
			{name: "create", fn: function() { pkg.create(lb) }},
			{name: "get", fn: function() { pkg.get(lb.id) }},
			{name: "delete", fn: function() { pkg.delete(lb) }},
			{name: "list", fn: function() { pkg.list() }},
			{name: "add_droplets", fn: function() { pkg.add_droplets(lb, [ dropletId ]) }},
			{name: "remove_droplets", fn: function() { pkg.remove_droplets(lb, [ dropletId ]) }},
			{name: "add_forwarding_rules", fn: function() { pkg.add_forwarding_rules(lb, lb.forwarding_rules) }},
			{name: "remove_forwarding_rules", fn: function() { pkg.remove_forwarding_rules(lb, lb.forwarding_rules) }}
		 ].forEach(function(kv) {
			var name = kv.name;
			var fn = kv.fn;

			try {
				fn(); throw "don't catch me";
			} catch(e) {
				equals("throw me", e.message, name + "should send the right exception!");
			}
		 });
	`)
}

var (
	region          = &godo.Region{Name: "newyork3", Slug: "nyc3", Sizes: []string{"small"}, Available: true, Features: []string{"all"}}
	forwardingRules = []godo.ForwardingRule{
		{
			EntryProtocol:  "http",
			EntryPort:      80,
			TargetProtocol: "http",
			TargetPort:     80,
			CertificateID:  "",
			TlsPassthrough: false,
		},
	}

	healthCheck = &godo.HealthCheck{
		Protocol:               "http",
		Port:                   80,
		Path:                   "/",
		CheckIntervalSeconds:   10,
		ResponseTimeoutSeconds: 5,
		HealthyThreshold:       5,
		UnhealthyThreshold:     3,
	}

	stickySessions = &godo.StickySessions{
		Type: "none",
	}

	l = &godo.LoadBalancer{ID: "test-uuid", Name: "example-lb-01", HealthCheck: healthCheck, StickySessions: stickySessions, ForwardingRules: forwardingRules, DropletIDs: []int{3164444, 3164445}, Status: "new", Algorithm: "round_robin", Region: region}
)

type loadBalancer struct {
	*godo.LoadBalancer
}

func (k *loadBalancer) Struct() *godo.LoadBalancer { return k.LoadBalancer }

func TestLoadBalancerCreate(t *testing.T) {
	cloud := mockcloud.Client(nil)

	cloud.MockLoadBalancers.CreateFn = func(_ context.Context, name, region string, forwardingRules []godo.ForwardingRule, opts ...loadbalancers.CreateOpt) (loadbalancers.LoadBalancer, error) {
		return &loadBalancer{l}, nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;

		var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

		var l = pkg.create({
			"name": "example-lb-01",
			"region": "nyc3",
			"forwarding_rules": [
			{
				"entry_protocol": "http",
				"entry_port": 80,
				"target_protocol": "http",
				"target_port": 80,
				"certificate_id": "",
				"tls_passthrough": false
			}
			],
			"health_check": {
				"protocol": "http",
				"port": 80,
				"path": "/",
				"check_interval_seconds": 10,
				"response_timeout_seconds": 5,
				"healthy_threshold": 5,
				"unhealthy_threshold": 3
			},
			"sticky_sessions": {
				"type": "none"
			},
			"droplet_ids": [
			3164444,
			3164445
			]
		});

		var want = {
			"id": "test-uuid",
			"name": "example-lb-01",
			"ip": "",
			"algorithm": "round_robin",
			"status": "new",
			"created_at": "",
			"forwarding_rules": [
			{
				"entry_protocol": "http",
				"entry_port": 80,
				"target_protocol": "http",
				"target_port": 80,
				"certificate_id": "",
				"tls_passthrough": false
			}
			],

			"health_check": {
				"protocol": "http",
				"port": 80,
				"path": "/",
				"check_interval_seconds": 10,
				"response_timeout_seconds": 5,
				"healthy_threshold": 5,
				"unhealthy_threshold": 3
			},
			"sticky_sessions": {
				"type": "none",
				cookie_name: "",
				cookie_ttl_seconds: 0,
			},
			"region": region,
			"tag": "",
			"droplet_ids": [
			3164444,
			3164445
			],
			"redirect_http_to_https": false
		};

		equals(l, want, "should have proper object");
	`)
}

func TestLoadBalancersList(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockLoadBalancers.ListFn = func(_ context.Context) (<-chan loadbalancers.LoadBalancer, <-chan error) {
		lc := make(chan loadbalancers.LoadBalancer, 1)
		lc <- &loadBalancer{l}
		close(lc)
		ec := make(chan error)
		close(ec)
		return lc, ec
	}

	vmtest.Run(t, cloud, `
	var pkg = cloud.load_balancers;
	var list = pkg.list();
	assert(list != null, "should have received a list");
	assert(list.length > 0, "should have received some elements");

	var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };
	var want = {
		"id": "test-uuid",
		"name": "example-lb-01",
		"ip": "",
		"algorithm": "round_robin",
		"status": "new",
		"created_at": "",
		"forwarding_rules": [
		{
			"entry_protocol": "http",
			"entry_port": 80,
			"target_protocol": "http",
			"target_port": 80,
			"certificate_id": "",
			"tls_passthrough": false
		}
		],

		"health_check": {
			"protocol": "http",
			"port": 80,
			"path": "/",
			"check_interval_seconds": 10,
			"response_timeout_seconds": 5,
			"healthy_threshold": 5,
			"unhealthy_threshold": 3
		},
		"sticky_sessions": {
			"type": "none",
			cookie_name: "",
			cookie_ttl_seconds: 0,
		},
		"region": region,
		"tag": "",
		"droplet_ids": [
		3164444,
		3164445
		],
		"redirect_http_to_https": false
	};	

	var l = list[0];

	equals(l, want, "should have proper object");
	`)
}

func TestLoadBalancerDelete(t *testing.T) {
	wantId := "test-uuid"
	cloud := mockcloud.Client(nil)

	cloud.MockLoadBalancers.DeleteFn = func(_ context.Context, gotId string) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
			var pkg = cloud.load_balancers;
			pkg.delete("test-uuid");
	`)
}

func TestLoadBalancerGet(t *testing.T) {
	cloud := mockcloud.Client(nil)
	cloud.MockLoadBalancers.GetFn = func(_ context.Context, id string) (loadbalancers.LoadBalancer, error) {
		return &loadBalancer{l}, nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;
		var region = { name: "newyork3", slug: "nyc3", sizes: ["small"], available: true, features: ["all"] };

		var want = {
			"id": "test-uuid",
			"name": "example-lb-01",
			"ip": "",
			"algorithm": "round_robin",
			"status": "new",
			"created_at": "",
		"forwarding_rules": [
			{
				"entry_protocol": "http",
				"entry_port": 80,
				"target_protocol": "http",
				"target_port": 80,
				"certificate_id": "",
				"tls_passthrough": false
			}
			],

			"health_check": {
				"protocol": "http",
				"port": 80,
				"path": "/",
				"check_interval_seconds": 10,
				"response_timeout_seconds": 5,
				"healthy_threshold": 5,
				"unhealthy_threshold": 3
			},
			"sticky_sessions": {
				"type": "none",
				cookie_name: "",
				cookie_ttl_seconds: 0,
			},
			"region": region,
			"tag": "",
			"droplet_ids": [
			3164444,
			3164445
			],
			"redirect_http_to_https": false
		};

		var l = pkg.get('test-uuid');

		equals(l, want, "should have proper object");
	`)
}

func TestLoadBalancerAddForwardingRules(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockLoadBalancers.AddForwardingRulesFn = func(_ context.Context, gotId string, rules ...godo.ForwardingRule) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantRule := godo.ForwardingRule{
			EntryProtocol:  "http",
			EntryPort:      81,
			TargetProtocol: "http",
			TargetPort:     81,
			CertificateID:  "",
			TlsPassthrough: false,
		}

		gotRule := rules[0]

		if wantRule != gotRule {
			t.Fatalf("want %v got %v", wantRule, gotRule)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;
		pkg.add_forwarding_rules("test-uuid", [{
			entry_protocol: "http", 
			entry_port: 81, 
			target_protocol: "http", 
			target_port: 81
		}
		]);
	`)
}

func TestLoadBalancerRemoveForwardingRules(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockLoadBalancers.RemoveForwardingRulesFn = func(_ context.Context, gotId string, rules ...godo.ForwardingRule) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantRule := godo.ForwardingRule{
			EntryProtocol:  "http",
			EntryPort:      81,
			TargetProtocol: "http",
			TargetPort:     81,
			CertificateID:  "",
			TlsPassthrough: false,
		}

		gotRule := rules[0]

		if wantRule != gotRule {
			t.Fatalf("want %v got %v", wantRule, gotRule)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;
		pkg.remove_forwarding_rules("test-uuid", [{
			entry_protocol: "http", 
			entry_port: 81, 
			target_protocol: "http", 
			target_port: 81
		}
		]);
	`)
}

func TestLoadBalancerAddDroplets(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockLoadBalancers.AddDropletsFn = func(_ context.Context, gotId string, dropletIds ...int) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantDropletId := 42

		gotDropletId := dropletIds[0]
		if wantDropletId != gotDropletId {
			t.Fatalf("want %v got %v", wantDropletId, gotDropletId)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;

		pkg.add_droplets("test-uuid", [42]);
	`)
}

func TestLoadBalancerRemoveDroplets(t *testing.T) {
	cloud := mockcloud.Client(nil)
	wantId := "test-uuid"
	cloud.MockLoadBalancers.RemoveDropletsFn = func(_ context.Context, gotId string, dropletIds ...int) error {
		if gotId != wantId {
			t.Fatalf("want %v got %v", wantId, gotId)
		}

		wantDropletId := 42

		gotDropletId := dropletIds[0]
		if wantDropletId != gotDropletId {
			t.Fatalf("want %v got %v", wantDropletId, gotDropletId)
		}

		return nil
	}

	vmtest.Run(t, cloud, `
		var pkg = cloud.load_balancers;

		pkg.remove_droplets("test-uuid", [42]);
	`)
}
