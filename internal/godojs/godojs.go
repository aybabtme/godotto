package godojs

import (
	"strconv"
	"time"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

// from VM

func ArgActionID(vm *otto.Otto, v otto.Value) int {
	var aid int
	switch {
	case v.IsNumber():
		aid = ottoutil.Int(vm, v)
	case v.IsObject():
		aid = ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false))
	default:
		ottoutil.Throw(vm, "argument must be an Action or an ActionID")
	}
	return aid
}

func ArgDomainCreateRequest(vm *otto.Otto, v otto.Value) *godo.DomainCreateRequest {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Domain, got a %q", v.Class())
	}
	return &godo.DomainCreateRequest{
		Name:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		IPAddress: ottoutil.String(vm, ottoutil.GetObject(vm, v, "ip_address", true)),
	}
}

func ArgDomainName(vm *otto.Otto, v otto.Value) string {
	var name string
	switch {
	case v.IsString():
		name = ottoutil.String(vm, v)
	case v.IsObject():
		name = ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false))
	default:
		ottoutil.Throw(vm, "argument must be a Domain or a DomainName")
	}
	return name
}

func ArgRecordID(vm *otto.Otto, v otto.Value) int {
	var id int
	switch {
	case v.IsNumber():
		id = ottoutil.Int(vm, v)
	case v.IsObject():
		id = ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false))
	default:
		ottoutil.Throw(vm, "argument must be a Domain or a DomainName")
	}
	return id
}

func ArgDomainRecord(vm *otto.Otto, v otto.Value) *godo.DomainRecordEditRequest {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a DomainRecord, got a %q", v.Class())
	}
	return &godo.DomainRecordEditRequest{
		Type:     ottoutil.String(vm, ottoutil.GetObject(vm, v, "type", true)),
		Name:     ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Data:     ottoutil.String(vm, ottoutil.GetObject(vm, v, "data", true)),
		Priority: ottoutil.Int(vm, ottoutil.GetObject(vm, v, "priority", true)),
		Port:     ottoutil.Int(vm, ottoutil.GetObject(vm, v, "port", true)),
		Weight:   ottoutil.Int(vm, ottoutil.GetObject(vm, v, "weight", true)),
	}
}

func ArgDroplet(vm *otto.Otto, v otto.Value) *godo.Droplet {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Droplet, got a %q", v.Class())
	}
	return &godo.Droplet{
		ID:          ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false)),
		Name:        ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		Memory:      ottoutil.Int(vm, ottoutil.GetObject(vm, v, "memory", false)),
		Vcpus:       ottoutil.Int(vm, ottoutil.GetObject(vm, v, "vcpus", false)),
		Disk:        ottoutil.Int(vm, ottoutil.GetObject(vm, v, "disk", false)),
		Region:      ArgRegion(vm, ottoutil.GetObject(vm, v, "region", false)),
		Image:       ArgImage(vm, ottoutil.GetObject(vm, v, "image", false)),
		Size:        ArgSize(vm, ottoutil.GetObject(vm, v, "size", false)),
		SizeSlug:    ottoutil.String(vm, ottoutil.GetObject(vm, v, "size_slug", false)),
		BackupIDs:   ottoutil.IntSlice(vm, ottoutil.GetObject(vm, v, "backup_ids", false)),
		SnapshotIDs: ottoutil.IntSlice(vm, ottoutil.GetObject(vm, v, "snapshot_ids", false)),
		Locked:      ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "locked", false)),
		Status:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "status", false)),
		Networks:    ArgNetworks(vm, ottoutil.GetObject(vm, v, "networks", false)),
		Kernel:      ArgKernel(vm, ottoutil.GetObject(vm, v, "kernel", false)),
		Tags:        ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "tags", false)),
		VolumeIDs:   ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "volumes", false)),
	}
}

func ArgDropletID(vm *otto.Otto, v otto.Value) int {
	var did int
	switch {
	case v.IsNumber():
		did = ottoutil.Int(vm, v)
	case v.IsObject():
		did = ArgDroplet(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be a Droplet or a DropletID")
	}
	return did
}

func ArgTagCreateRequest(vm *otto.Otto, v otto.Value) *godo.TagCreateRequest {
	req := &godo.TagCreateRequest{
		Name: ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
	}

	return req
}

func ArgTagTagResourcesRequest(vm *otto.Otto, v otto.Value) *godo.TagResourcesRequest {
	req := &godo.TagResourcesRequest{}

	resArgs := ottoutil.GetObject(vm, v, "resources", true)
	ottoutil.LoadArray(vm, resArgs, func(v otto.Value) {
		res := ArgResource(vm, v)
		req.Resources = append(req.Resources, godo.Resource{
			ID:   res.ID,
			Type: res.Type,
		})
	})

	return req
}

func ArgTagUntagResourcesRequest(vm *otto.Otto, v otto.Value) *godo.UntagResourcesRequest {
	req := &godo.UntagResourcesRequest{}

	resArgs := ottoutil.GetObject(vm, v, "resources", true)
	ottoutil.LoadArray(vm, resArgs, func(v otto.Value) {
		res := ArgResource(vm, v)
		req.Resources = append(req.Resources, godo.Resource{
			ID:   res.ID,
			Type: res.Type,
		})
	})

	return req
}

func ArgLoadBalancerCreateRequest(vm *otto.Otto, v otto.Value) *godo.LoadBalancerRequest {
	req := &godo.LoadBalancerRequest{
		Name:                ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Algorithm:           ottoutil.String(vm, ottoutil.GetObject(vm, v, "algorithm", false)),
		Region:              ArgRegionSlug(vm, ottoutil.GetObject(vm, v, "region", true)),
		DropletIDs:          ArgDropletIDs(vm, ottoutil.GetObject(vm, v, "droplet_ids", false)),
		HealthCheck:         ArgHealthCheck(vm, ottoutil.GetObject(vm, v, "health_check", false)),
		StickySessions:      ArgStickySessions(vm, ottoutil.GetObject(vm, v, "sticky_sessions", false)),
		ForwardingRules:     ArgForwardingRules(vm, ottoutil.GetObject(vm, v, "forwarding_rules", true)),
		Tag:                 ottoutil.String(vm, ottoutil.GetObject(vm, v, "tag", false)),
		RedirectHttpToHttps: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "redirect_http_to_https", false)),
	}

	return req
}

func ArgLoadBalancerUpdate(vm *otto.Otto, v otto.Value) *godo.LoadBalancerRequest {
	return &godo.LoadBalancerRequest{
		Name:                ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Algorithm:           ottoutil.String(vm, ottoutil.GetObject(vm, v, "algorithm", false)),
		Region:              ArgRegionSlug(vm, ottoutil.GetObject(vm, v, "region", true)),
		DropletIDs:          ArgDropletIDs(vm, ottoutil.GetObject(vm, v, "droplet_ids", false)),
		HealthCheck:         ArgHealthCheck(vm, ottoutil.GetObject(vm, v, "health_check", false)),
		StickySessions:      ArgStickySessions(vm, ottoutil.GetObject(vm, v, "sticky_sessions", false)),
		ForwardingRules:     ArgForwardingRules(vm, ottoutil.GetObject(vm, v, "forwarding_rules", true)),
		Tag:                 ottoutil.String(vm, ottoutil.GetObject(vm, v, "tag", false)),
		RedirectHttpToHttps: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "redirect_http_to_https", false)),
	}
}

func ArgDropletIDs(vm *otto.Otto, v otto.Value) []int {
	ids := make([]int, 0)
	ottoutil.LoadArray(vm, v, func(v otto.Value) {
		did := ArgDropletID(vm, v)
		ids = append(ids, did)
	})

	return ids
}

func ArgForwardingRules(vm *otto.Otto, v otto.Value) []godo.ForwardingRule {
	rules := make([]godo.ForwardingRule, 0)
	ottoutil.LoadArray(vm, v, func(v otto.Value) {
		rule := ArgForwardingRule(vm, v)
		rules = append(rules, godo.ForwardingRule{
			EntryProtocol:  rule.EntryProtocol,
			EntryPort:      rule.EntryPort,
			TargetProtocol: rule.TargetProtocol,
			TargetPort:     rule.TargetPort,
			CertificateID:  rule.CertificateID,
			TlsPassthrough: rule.TlsPassthrough,
		})
	})

	return rules
}

func ArgLoadBalancer(vm *otto.Otto, v otto.Value) *godo.LoadBalancer {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a LoadBalancer, got a %q", v.Class())
	}

	return &godo.LoadBalancer{
		ID:                  ottoutil.String(vm, ottoutil.GetObject(vm, v, "id", false)),
		Name:                ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		IP:                  ottoutil.String(vm, ottoutil.GetObject(vm, v, "ip", false)),
		Algorithm:           ottoutil.String(vm, ottoutil.GetObject(vm, v, "algorithm", false)),
		Status:              ottoutil.String(vm, ottoutil.GetObject(vm, v, "status", false)),
		ForwardingRules:     ArgForwardingRules(vm, ottoutil.GetObject(vm, v, "forwarding_rules", false)),
		HealthCheck:         ArgHealthCheck(vm, ottoutil.GetObject(vm, v, "health_check", false)),
		StickySessions:      ArgStickySessions(vm, ottoutil.GetObject(vm, v, "sticky_sesions", false)),
		Region:              ArgRegion(vm, ottoutil.GetObject(vm, v, "region", false)),
		Tag:                 ottoutil.String(vm, ottoutil.GetObject(vm, v, "tag", false)),
		DropletIDs:          ArgDropletIDs(vm, ottoutil.GetObject(vm, v, "droplet_ids", false)),
		RedirectHttpToHttps: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "redirect_http_to_https", false)),
	}
}

func ArgLoadBalancerID(vm *otto.Otto, v otto.Value) string {
	var lbID string
	switch {
	case v.IsString():
		lbID = ottoutil.String(vm, v)
	case v.IsObject():
		lbID = ArgLoadBalancer(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be a LoadBalancer or LoadBalancerID")
	}

	return lbID
}

func ArgHealthCheck(vm *otto.Otto, v otto.Value) *godo.HealthCheck {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}

	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a HealthCheck, got a %q", v.Class())
	}

	return &godo.HealthCheck{
		Protocol:               ottoutil.String(vm, ottoutil.GetObject(vm, v, "protocol", true)),
		Port:                   ottoutil.Int(vm, ottoutil.GetObject(vm, v, "port", true)),
		Path:                   ottoutil.String(vm, ottoutil.GetObject(vm, v, "path", false)),
		CheckIntervalSeconds:   ottoutil.Int(vm, ottoutil.GetObject(vm, v, "check_interval_seconds", false)),
		ResponseTimeoutSeconds: ottoutil.Int(vm, ottoutil.GetObject(vm, v, "response_timeout_seconds", false)),
		UnhealthyThreshold:     ottoutil.Int(vm, ottoutil.GetObject(vm, v, "unhealthy_threshold", false)),
		HealthyThreshold:       ottoutil.Int(vm, ottoutil.GetObject(vm, v, "healthy_threshold", false)),
	}
}

func ArgStickySessions(vm *otto.Otto, v otto.Value) *godo.StickySessions {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}

	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a StickySessions, got a %q", v.Class())
	}

	return &godo.StickySessions{
		Type:             ottoutil.String(vm, ottoutil.GetObject(vm, v, "type", false)),
		CookieName:       ottoutil.String(vm, ottoutil.GetObject(vm, v, "cookie_name", false)),
		CookieTtlSeconds: ottoutil.Int(vm, ottoutil.GetObject(vm, v, "cookie_ttl_seconds", false)),
	}
}

func ArgForwardingRule(vm *otto.Otto, v otto.Value) *godo.ForwardingRule {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}

	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a ForwardingRule, got a %q", v.Class())
	}

	return &godo.ForwardingRule{
		EntryProtocol:  ottoutil.String(vm, ottoutil.GetObject(vm, v, "entry_protocol", true)),
		EntryPort:      ottoutil.Int(vm, ottoutil.GetObject(vm, v, "entry_port", true)),
		TargetProtocol: ottoutil.String(vm, ottoutil.GetObject(vm, v, "target_protocol", true)),
		TargetPort:     ottoutil.Int(vm, ottoutil.GetObject(vm, v, "target_port", true)),
		CertificateID:  ottoutil.String(vm, ottoutil.GetObject(vm, v, "certificate_id", false)),
		TlsPassthrough: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "tls_passthrough", false)),
	}
}

func ArgResource(vm *otto.Otto, v otto.Value) *godo.Resource {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}

	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Resource, got a %q", v.Class())
	}

	return &godo.Resource{
		ID:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "id", true)),
		Type: ArgResourceType(vm, ottoutil.GetObject(vm, v, "type", true)),
	}
}

func ArgResourceType(vm *otto.Otto, v otto.Value) godo.ResourceType {
	return godo.ResourceType(ottoutil.String(vm, v))
}

func ArgDropletCreateRequest(vm *otto.Otto, v otto.Value) *godo.DropletCreateRequest {
	image := ArgImage(vm, ottoutil.GetObject(vm, v, "image", true))
	req := &godo.DropletCreateRequest{
		Name:              ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Region:            ottoutil.String(vm, ottoutil.GetObject(vm, v, "region", true)),
		Size:              ottoutil.String(vm, ottoutil.GetObject(vm, v, "size", true)),
		Image:             godo.DropletCreateImage{ID: image.ID, Slug: image.Slug},
		Backups:           ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "backups", false)),
		IPv6:              ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "ipv6", false)),
		PrivateNetworking: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "private_networking", false)),
		UserData:          ottoutil.String(vm, ottoutil.GetObject(vm, v, "user_data", false)),
		Monitoring:        ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "monitoring", false)),
		Tags:              ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "tags", false)),
	}
	sshArgs := ottoutil.GetObject(vm, v, "ssh_keys", false)
	ottoutil.LoadArray(vm, sshArgs, func(v otto.Value) {
		key := ArgKey(vm, v)
		req.SSHKeys = append(req.SSHKeys, godo.DropletCreateSSHKey{
			ID:          key.ID,
			Fingerprint: key.Fingerprint,
		})
	})
	volumeArgs := ottoutil.GetObject(vm, v, "volumes", false)
	ottoutil.LoadArray(vm, volumeArgs, func(v otto.Value) {
		volume := ArgVolume(vm, v)
		req.Volumes = append(req.Volumes, godo.DropletCreateVolume{
			ID:   volume.ID,
			Name: volume.Name,
		})
	})
	return req
}

func ArgImage(vm *otto.Otto, v otto.Value) *godo.Image {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Image, got a %q", v.Class())
	}
	return &godo.Image{
		ID:           ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false)),
		Name:         ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		Type:         ottoutil.String(vm, ottoutil.GetObject(vm, v, "type", false)),
		Distribution: ottoutil.String(vm, ottoutil.GetObject(vm, v, "distribution", false)),
		Slug:         ottoutil.String(vm, ottoutil.GetObject(vm, v, "slug", false)),
		Public:       ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "public", false)),
		Regions:      ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "regions", false)),
		MinDiskSize:  ottoutil.Int(vm, ottoutil.GetObject(vm, v, "min_disk_size", false)),
	}
}

func ArgImageID(vm *otto.Otto, v otto.Value) int {
	var imgID int
	switch {
	case v.IsNumber():
		imgID = ottoutil.Int(vm, v)
	case v.IsObject():
		imgID = ArgImage(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be an Image or a ImageID")
	}
	return imgID
}

func ArgImageSlug(vm *otto.Otto, v otto.Value) string {
	var slug string
	switch {
	case v.IsString():
		slug = ottoutil.String(vm, v)
	case v.IsObject():
		slug = ArgImage(vm, v).Slug
	default:
		ottoutil.Throw(vm, "argument must be an Image or a ImageSlug")
	}
	return slug
}

func ArgImageName(vm *otto.Otto, v otto.Value) string {
	var name string
	switch {
	case v.IsString():
		name = ottoutil.String(vm, v)
	case v.IsObject():
		name = ArgImage(vm, v).Name
	default:
		ottoutil.Throw(vm, "argument must be an Image or an ImageName")
	}
	return name
}

func ArgKernel(vm *otto.Otto, v otto.Value) *godo.Kernel {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Kernel, got a %#v", v)
	}
	return &godo.Kernel{
		ID:      ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false)),
		Name:    ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		Version: ottoutil.String(vm, ottoutil.GetObject(vm, v, "version", false)),
	}
}

func ArgKernelID(vm *otto.Otto, v otto.Value) int {
	var kernID int
	switch {
	case v.IsNumber():
		kernID = ottoutil.Int(vm, v)
	case v.IsObject():
		kernID = ArgKernel(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be a Kernel or a KernelID")
	}
	return kernID
}

func ArgNetworks(vm *otto.Otto, v otto.Value) *godo.Networks {
	net := &godo.Networks{}
	if v4Arg := ottoutil.GetObject(vm, v, "v4", false); v4Arg.IsDefined() {
		ottoutil.LoadArray(vm, v4Arg, func(v otto.Value) {
			net.V4 = append(net.V4, godo.NetworkV4{
				IPAddress: ottoutil.String(vm, ottoutil.GetObject(vm, v, "ip_address", false)),
				Netmask:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "netmask", false)),
				Gateway:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "gateway", false)),
				Type:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "type", false)),
			})
		})
	}

	if v6Arg := ottoutil.GetObject(vm, v, "v6", false); v6Arg.IsDefined() {
		ottoutil.LoadArray(vm, v6Arg, func(v otto.Value) {
			net.V6 = append(net.V6, godo.NetworkV6{
				IPAddress: ottoutil.String(vm, ottoutil.GetObject(vm, v, "ip_address", false)),
				Netmask:   ottoutil.Int(vm, ottoutil.GetObject(vm, v, "netmask", false)),
				Gateway:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "gateway", false)),
				Type:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "type", false)),
			})
		})
	}
	return net
}

func ArgSize(vm *otto.Otto, v otto.Value) *godo.Size {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Size, got a %q", v.Class())
	}
	return &godo.Size{
		Slug:         ottoutil.String(vm, ottoutil.GetObject(vm, v, "slug", false)),
		Memory:       ottoutil.Int(vm, ottoutil.GetObject(vm, v, "memory", false)),
		Vcpus:        ottoutil.Int(vm, ottoutil.GetObject(vm, v, "vcpus", false)),
		Disk:         ottoutil.Int(vm, ottoutil.GetObject(vm, v, "disk", false)),
		PriceMonthly: ottoutil.Float64(vm, ottoutil.GetObject(vm, v, "price_monthly", false)),
		PriceHourly:  ottoutil.Float64(vm, ottoutil.GetObject(vm, v, "price_hourly", false)),
		Regions:      ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "regions", false)),
		Available:    ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "available", false)),
		Transfer:     ottoutil.Float64(vm, ottoutil.GetObject(vm, v, "transfer", false)),
	}
}

func ArgSizeSlug(vm *otto.Otto, v otto.Value) string {
	var slug string
	switch {
	case v.IsString():
		slug = ottoutil.String(vm, v)
	case v.IsObject():
		slug = ArgSize(vm, v).Slug
	default:
		ottoutil.Throw(vm, "argument must be an Size or a SizeSlug")
	}
	return slug
}

func ArgRegion(vm *otto.Otto, v otto.Value) *godo.Region {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Region, got a %q", v.Class())
	}
	return &godo.Region{
		Slug:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "slug", false)),
		Name:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		Sizes:     ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "sizes", false)),
		Available: ottoutil.Bool(vm, ottoutil.GetObject(vm, v, "available", false)),
		Features:  ottoutil.StringSlice(vm, ottoutil.GetObject(vm, v, "features", false)),
	}
}

func ArgRegionSlug(vm *otto.Otto, v otto.Value) string {
	var slug string
	switch {
	case v.IsString():
		slug = ottoutil.String(vm, v)
	case v.IsObject():
		slug = ArgRegion(vm, v).Slug
	default:
		ottoutil.Throw(vm, "argument must be a Region or a RegionSlug")
	}
	return slug
}

func ArgVolume(vm *otto.Otto, v otto.Value) *godo.Volume {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Volume, got a %q", v.Class())
	}
	return &godo.Volume{
		ID:            ottoutil.String(vm, ottoutil.GetObject(vm, v, "id", false)),
		Region:        ArgRegion(vm, ottoutil.GetObject(vm, v, "region", false)),
		Name:          ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		SizeGigaBytes: int64(ottoutil.Int(vm, ottoutil.GetObject(vm, v, "size", false))),
		Description:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "desc", false)),
		DropletIDs:    ottoutil.IntSlice(vm, ottoutil.GetObject(vm, v, "droplet_ids", false)),
	}
}

func ArgVolumeCreateRequest(vm *otto.Otto, v otto.Value) *godo.VolumeCreateRequest {
	if !v.IsDefined() || v.IsNull() {
		ottoutil.Throw(vm, "argument must be a Volume create request, got nothing")
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Volume, got a %q", v.Class())
	}
	return &godo.VolumeCreateRequest{
		Name:          ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Region:        ArgRegionSlug(vm, ottoutil.GetObject(vm, v, "region", true)),
		SizeGigaBytes: int64(ottoutil.Int(vm, ottoutil.GetObject(vm, v, "size", true))),
		Description:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "desc", false)),
	}
}

func ArgVolumeID(vm *otto.Otto, v otto.Value) string {
	var volumeID string
	switch {
	case v.IsString():
		volumeID = ottoutil.String(vm, v)
	case v.IsObject():
		volumeID = ArgVolume(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be an Volume or a VolumeID")
	}
	return volumeID
}

func ArgSnapshotCreateRequest(vm *otto.Otto, v otto.Value) *godo.SnapshotCreateRequest {
	if !v.IsDefined() || v.IsNull() {
		ottoutil.Throw(vm, "argument must be a Snapshot create request, got nothing")
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Snapshot, got a %q", v.Class())
	}
	return &godo.SnapshotCreateRequest{
		VolumeID:    ArgVolumeID(vm, ottoutil.GetObject(vm, v, "volume", true)),
		Name:        ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", true)),
		Description: ottoutil.String(vm, ottoutil.GetObject(vm, v, "desc", false)),
	}
}

func ArgSnapshot(vm *otto.Otto, v otto.Value) *godo.Snapshot {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Snapshot, got a %q", v.Class())
	}
	return &godo.Snapshot{
		ID:         ottoutil.String(vm, ottoutil.GetObject(vm, v, "id", false)),
		ResourceID: ottoutil.String(vm, ottoutil.GetObject(vm, v, "volume_id", false)),
		Regions: []string{
			ArgRegion(vm, ottoutil.GetObject(vm, v, "region", false)).Slug,
		},
		Name:          ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		SizeGigaBytes: float64(ottoutil.Int(vm, ottoutil.GetObject(vm, v, "size", false))),
	}
}

func ArgSnapshotID(vm *otto.Otto, v otto.Value) string {
	var id string
	switch {
	case v.IsString():
		id = ottoutil.String(vm, v)
	case v.IsObject():
		id = ArgSnapshot(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be a Snapshot or a SnapshotID")
	}
	return id
}

func ArgFloatingIPCreateRequest(vm *otto.Otto, v otto.Value) *godo.FloatingIPCreateRequest {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a FloatingIP, got a %q", v.Class())
	}
	req := &godo.FloatingIPCreateRequest{
		Region: ArgRegionSlug(vm, ottoutil.GetObject(vm, v, "region", true)),
	}

	if v := ottoutil.GetObject(vm, v, "droplet", false); v.IsDefined() {
		req.DropletID = ArgDropletID(vm, ottoutil.GetObject(vm, v, "droplet", false))
	}
	return req
}

func ArgFloatingIP(vm *otto.Otto, v otto.Value) *godo.FloatingIP {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a FloatingIP, got a %q", v.Class())
	}
	return &godo.FloatingIP{
		Region:  ArgRegion(vm, ottoutil.GetObject(vm, v, "region", false)),
		Droplet: ArgDroplet(vm, ottoutil.GetObject(vm, v, "droplet", false)),
		IP:      ottoutil.String(vm, ottoutil.GetObject(vm, v, "ip", true)),
	}
}

func ArgFloatingIPActualIP(vm *otto.Otto, v otto.Value) string {
	var ip string
	switch {
	case v.IsString():
		ip = ottoutil.String(vm, v)
	case v.IsObject():
		ip = ArgFloatingIP(vm, v).IP
	default:
		ottoutil.Throw(vm, "argument must be a FloatingIP or an IP")
	}
	return ip
}

func ArgKey(vm *otto.Otto, v otto.Value) *godo.Key {
	if !v.IsDefined() || v.IsNull() {
		return nil
	}
	if !v.IsObject() {
		ottoutil.Throw(vm, "argument must be a Key, got a %q", v.Class())
	}
	return &godo.Key{
		ID:          ottoutil.Int(vm, ottoutil.GetObject(vm, v, "id", false)),
		Name:        ottoutil.String(vm, ottoutil.GetObject(vm, v, "name", false)),
		Fingerprint: ottoutil.String(vm, ottoutil.GetObject(vm, v, "fp", false)),
		PublicKey:   ottoutil.String(vm, ottoutil.GetObject(vm, v, "public_key", false)),
	}
}

func ArgKeyID(vm *otto.Otto, v otto.Value) int {
	var id int
	switch {
	case v.IsNumber():
		id = ottoutil.Int(vm, v)
	case v.IsObject():
		id = ArgKey(vm, v).ID
	default:
		ottoutil.Throw(vm, "argument must be a Key or a KeyID")
	}
	return id
}

func ArgKeyFingerprint(vm *otto.Otto, v otto.Value) string {
	var fp string
	switch {
	case v.IsString():
		fp = ottoutil.String(vm, v)
	case v.IsObject():
		fp = ArgKey(vm, v).Fingerprint
	default:
		ottoutil.Throw(vm, "argument must be a Key or a KeyFingerprint")
	}
	return fp
}

func ArgKeyCreate(vm *otto.Otto, v otto.Value) *godo.KeyCreateRequest {
	k := ArgKey(vm, v)
	return &godo.KeyCreateRequest{
		Name:      k.Name,
		PublicKey: k.PublicKey,
	}
}

func ArgKeyUpdate(vm *otto.Otto, v otto.Value) *godo.KeyUpdateRequest {
	k := ArgKey(vm, v)
	return &godo.KeyUpdateRequest{
		Name: k.Name,
	}
}

// to VM

func AccountToVM(vm *otto.Otto, g *godo.Account) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"droplet_limit":     int64(g.DropletLimit),
		"floating_ip_limit": int64(g.FloatingIPLimit),
		"email":             g.Email,
		"uuid":              g.UUID,
		"email_verified":    g.EmailVerified,
		"status":            g.Status,
		"status_message":    g.StatusMessage,
	})
}

func ActionToVM(vm *otto.Otto, g *godo.Action) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":            int64(g.ID),
		"status":        g.Status,
		"type":          g.Type,
		"started_at":    g.StartedAt.Format(time.RFC3339Nano),
		"completed_at":  g.CompletedAt.Format(time.RFC3339Nano),
		"resource_id":   int64(g.ResourceID),
		"resource_type": g.ResourceType,
		"region_slug":   g.RegionSlug,
	})
}

func DomainToVM(vm *otto.Otto, g *godo.Domain) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"name":      g.Name,
		"ttl":       int64(g.TTL),
		"zone_file": g.ZoneFile,
	})
}

func DomainRecordToVM(vm *otto.Otto, g *godo.DomainRecord) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":       int64(g.ID),
		"type":     g.Type,
		"name":     g.Name,
		"data":     g.Data,
		"priority": int64(g.Priority),
		"port":     int64(g.Port),
		"weight":   int64(g.Weight),
	})
}

func VolumeToVM(vm *otto.Otto, g *godo.Volume) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":          g.ID,
		"name":        g.Name,
		"region":      RegionToVM(vm, g.Region),
		"size":        int64(g.SizeGigaBytes),
		"description": g.Description,
		"droplet_ids": intsToInt64s(g.DropletIDs),
	})
}

func VolumeSnapshotToVM(vm *otto.Otto, g *godo.Snapshot) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":        g.ID,
		"volume_id": g.ResourceID,
		"name":      g.Name,
		"regions":   g.Regions,
		"size":      int64(g.SizeGigaBytes),
	})
}

func DropletToVM(vm *otto.Otto, g *godo.Droplet) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	publicIPv4, _ := g.PublicIPv4()
	publicIPv6, _ := g.PublicIPv6()

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":     int64(g.ID),
		"name":   g.Name,
		"memory": int64(g.Memory),

		"vcpus":        int64(g.Vcpus),
		"disk":         int64(g.Disk),
		"region":       RegionToVM(vm, g.Region),
		"image":        ImageToVM(vm, g.Image),
		"size":         SizeToVM(vm, g.Size),
		"size_slug":    g.SizeSlug,
		"backup_ids":   intsToInt64s(g.BackupIDs),
		"snapshot_ids": intsToInt64s(g.SnapshotIDs),
		"locked":       g.Locked,
		"status":       g.Status,
		"networks":     NetworksToVM(vm, g.Networks),
		"created_at":   g.Created,
		"kernel":       KernelToVM(vm, g.Kernel),
		"tags":         g.Tags,
		"volumes":      g.VolumeIDs,

		// extra

		"public_ipv4": publicIPv4,
		"public_ipv6": publicIPv6,
	})
}

func FloatingIPToVM(vm *otto.Otto, g *godo.FloatingIP) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	fields := map[string]interface{}{
		"region": RegionToVM(vm, g.Region),
		"ip":     g.IP,
	}
	if g.Droplet != nil {
		fields["droplet"] = DropletToVM(vm, g.Droplet)
	}
	return ottoutil.ToPkg(vm, fields)
}

func ImageToVM(vm *otto.Otto, g *godo.Image) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":            int64(g.ID),
		"name":          g.Name,
		"type":          g.Type,
		"distribution":  g.Distribution,
		"slug":          g.Slug,
		"public":        g.Public,
		"regions":       g.Regions,
		"min_disk_size": int64(g.MinDiskSize),
	})
}

func TagToVM(vm *otto.Otto, g *godo.Tag) otto.Value {
	if g == nil {
		return otto.NullValue()
	}

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"name":      g.Name,
		"resources": g.Resources,
	})
}

func LoadBalancerToVM(vm *otto.Otto, g *godo.LoadBalancer) otto.Value {
	if g == nil {
		return otto.NullValue()
	}

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":                     g.ID,
		"name":                   g.Name,
		"ip":                     g.IP,
		"algorithm":              g.Algorithm,
		"status":                 g.Status,
		"created_at":             g.Created,
		"forwarding_rules":       ForwardingRulesToVM(vm, g.ForwardingRules),
		"health_check":           HealthCheckToVM(vm, g.HealthCheck),
		"sticky_sessions":        StickySessionsToVM(vm, g.StickySessions),
		"region":                 RegionToVM(vm, g.Region),
		"tag":                    g.Tag,
		"droplet_ids":            intsToInt64s(g.DropletIDs),
		"redirect_http_to_https": g.RedirectHttpToHttps,
	})
}

func ForwardingRulesToVM(vm *otto.Otto, g []godo.ForwardingRule) otto.Value {
	if g == nil {
		return otto.NullValue()
	}

	var rules = make([]map[string]interface{}, 0)
	for _, rule := range g {
		rules = append(rules, ForwardingRuleToVM(vm, rule))
	}

	v, err := vm.ToValue(rules)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	return v
}

func ForwardingRuleToVM(vm *otto.Otto, g godo.ForwardingRule) map[string]interface{} {
	return map[string]interface{}{
		"entry_protocol":  g.EntryProtocol,
		"entry_port":      int64(g.EntryPort),
		"target_protocol": g.TargetProtocol,
		"target_port":     int64(g.TargetPort),
		"certificate_id":  g.CertificateID,
		"tls_passthrough": g.TlsPassthrough,
	}
}

func HealthCheckToVM(vm *otto.Otto, g *godo.HealthCheck) otto.Value {
	if g == nil {
		return otto.NullValue()
	}

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"protocol": g.Protocol,
		"port":     int64(g.Port),
		"path":     g.Path,
		"check_interval_seconds":   int64(g.CheckIntervalSeconds),
		"response_timeout_seconds": int64(g.ResponseTimeoutSeconds),
		"unhealthy_threshold":      int64(g.UnhealthyThreshold),
		"healthy_threshold":        int64(g.HealthyThreshold),
	})
}

func StickySessionsToVM(vm *otto.Otto, g *godo.StickySessions) otto.Value {
	if g == nil {
		return otto.NullValue()
	}

	return ottoutil.ToPkg(vm, map[string]interface{}{
		"type":               g.Type,
		"cookie_name":        g.CookieName,
		"cookie_ttl_seconds": int64(g.CookieTtlSeconds),
	})
}

func KeyToVM(vm *otto.Otto, g *godo.Key) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":          g.ID,
		"name":        g.Name,
		"fingerprint": g.Fingerprint,
		"public_key":  g.PublicKey,
	})
}

func RegionToVM(vm *otto.Otto, g *godo.Region) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"slug":      g.Slug,
		"name":      g.Name,
		"sizes":     g.Sizes,
		"available": g.Available,
		"features":  g.Features,
	})
}

func SizeToVM(vm *otto.Otto, g *godo.Size) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"slug":          g.Slug,
		"memory":        int64(g.Memory),
		"vcpus":         int64(g.Vcpus),
		"disk":          int64(g.Disk),
		"price_monthly": g.PriceMonthly,
		"price_hourly":  g.PriceHourly,
		"regions":       g.Regions,
		"available":     g.Available,
		"transfer":      g.Transfer,
	})
}

func NetworksToVM(vm *otto.Otto, g *godo.Networks) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	var networkV4 map[string]interface{}
	if len(g.V4) != 0 {
		networkV4 = make(map[string]interface{})
		for i, v4 := range g.V4 {
			key := strconv.Itoa(i)
			networkV4[key] = ottoutil.ToPkg(vm, map[string]interface{}{
				"gateway":    v4.Gateway,
				"ip_address": v4.IPAddress,
				"netmask":    v4.Netmask,
				"type":       v4.Type,
			})
		}
	}
	var networkV6 map[string]interface{}
	if len(g.V6) != 0 {
		networkV6 = make(map[string]interface{})
		for i, v6 := range g.V6 {
			key := strconv.Itoa(i)
			networkV6[key] = ottoutil.ToPkg(vm, map[string]interface{}{
				"gateway":    v6.Gateway,
				"ip_address": v6.IPAddress,
				"netmask":    v6.Netmask,
				"type":       v6.Type,
			})
		}
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"v4": networkV4, "v6": networkV6,
	})
}
func KernelToVM(vm *otto.Otto, g *godo.Kernel) otto.Value {
	if g == nil {
		return otto.NullValue()
	}
	return ottoutil.ToPkg(vm, map[string]interface{}{
		"id":      g.ID,
		"name":    g.Name,
		"version": g.Version,
	})
}

// helpers

func intsToInt64s(in []int) (out []int64) {
	for _, i := range in {
		out = append(out, int64(i))
	}
	return out
}
