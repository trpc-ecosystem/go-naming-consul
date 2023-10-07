// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package discovery

import (
	"github.com/hashicorp/consul/api"
	"golang.org/x/sync/singleflight"
	"trpc.group/trpc-go/trpc-go/log"
	tdiscovery "trpc.group/trpc-go/trpc-go/naming/discovery"
	"trpc.group/trpc-go/trpc-go/naming/registry"
	consul_error "trpc.group/trpc-go/trpc-naming-consul/error"
)

// DefaultDiscovery instantiated objects by Discovery structure.
var DefaultDiscovery *Discovery

// Discovery structure.
type Discovery struct {
	opts  *Options
	cache *cache
	sg    singleflight.Group
}

type serviceNodes struct {
	HealthyNodes   []*registry.Node
	UnhealthyNodes []*registry.Node
}

// New instantiates discovery.
func New(options ...Option) (*Discovery, error) {
	d := &Discovery{
		opts: &Options{},
	}
	for _, o := range options {
		o(d.opts)
	}
	var err error
	d.cache, err = newCache(options...)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// List gets available service nodes, including only the healthy ones.
func (d *Discovery) List(serviceName string, opts ...tdiscovery.Option) ([]*registry.Node, error) {
	nodes, _, err := d.ListAll(serviceName, opts...)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, consul_error.ServerNotAvailableError
	}
	return nodes, nil
}

// ListAll gets all service nodes, including healthy and unhealthy ones.
func (d *Discovery) ListAll(serviceName string, opts ...tdiscovery.Option) (healthyNodes []*registry.Node,
	unhealthyNodes []*registry.Node, err error) {
	var nodes *serviceNodes
	nodes, err = d.cache.List(serviceName)
	if nodes != nil {
		return nodes.HealthyNodes, nodes.UnhealthyNodes, err
	}

	// The cache is not found, go to consul to get it
	val, err, _ := d.sg.Do(serviceName, func() (interface{}, error) {
		nodes, err = d.cache.List(serviceName)
		if err != nil || healthyNodes != nil {
			return &serviceNodes{HealthyNodes: healthyNodes, UnhealthyNodes: unhealthyNodes}, err
		}

		o := &tdiscovery.Options{}
		for _, opt := range opts {
			opt(o)
		}
		queryOpts := &api.QueryOptions{}
		queryOpts = queryOpts.WithContext(o.Ctx)
		serviceEntries, queryMeta, err := d.opts.client.Health().Service(serviceName, "", true, queryOpts)
		if err != nil {
			return nil, err
		}
		var healthEntries, unhealthEntries []*api.ServiceEntry
		for _, service := range serviceEntries {
			if service.Checks.AggregatedStatus() == api.HealthPassing {
				healthEntries = append(healthEntries, service)
			} else {
				unhealthEntries = append(unhealthEntries, service)
			}
		}
		nodes := &serviceNodes{
			HealthyNodes:   convertNodes(healthEntries),
			UnhealthyNodes: convertNodes(unhealthEntries),
		}
		_ = d.cache.cache(serviceName, queryMeta.LastIndex, nodes)
		return nodes, nil
	})
	if err != nil {
		log.Errorf("Discovery::ListAll failed to get serviceName:%s nodes from consul , err: %s", serviceName, err)
		return nil, nil, err
	}
	result, _ := val.(*serviceNodes)
	if result != nil {
		return result.HealthyNodes, result.UnhealthyNodes, nil
	}
	return nil, nil, nil
}
