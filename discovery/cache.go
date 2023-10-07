// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package discovery

import (
	"net"
	"strconv"
	"sync"

	"github.com/hashicorp/consul/api"
	tregistry "trpc.group/trpc-go/trpc-go/naming/registry"
)

var (
	emptyNodes = make([]*tregistry.Node, 0)
)

// The cache service caches the consul service registration information
// to prevent consul from being overly pressured by each request to consul.
type cache struct {
	opts *Options
	sync.RWMutex

	// Cache service.
	nodesCache map[string]*serviceNodes

	// Observe whether watched has been called, and modify the value when called.
	watched map[string]bool
	watcher *consulWatcher
	// The current cached consul data version.
	version uint64
	// Exit.
	exit chan bool
}

// setLocked function sets up the service nodes, must guarded by write lock and then operate.
func (c *cache) setLocked(serviceName string, nodes *serviceNodes) {
	if nodes == nil {
		return
	}
	c.nodesCache[serviceName] = nodes
}

// cache service nodes.
func (c *cache) cache(serviceName string, version uint64, nodes *serviceNodes) error {
	if nodes == nil {
		return nil
	}
	if nodes.HealthyNodes == nil {
		nodes.HealthyNodes = emptyNodes
	}
	if nodes.UnhealthyNodes == nil {
		nodes.UnhealthyNodes = emptyNodes
	}

	c.Lock()
	defer c.Unlock()
	if _, ok := c.watched[serviceName]; !ok {
		return nil
	}
	if version > c.version {
		c.version = version
	}
	c.setLocked(serviceName, nodes)
	return nil
}

// update updates the cache according to consul changes.
func (c *cache) update(result *watchResult) {
	if result == nil || result.healthyEntries == nil {
		return
	}

	c.Lock()
	defer c.Unlock()
	serviceName := result.serviceName

	// Nodes that are not concerned return directly.
	if _, ok := c.watched[serviceName]; !ok {
		return
	}
	_, ok := c.nodesCache[serviceName]
	if !ok {
		// Incremental quantity updates only start after getting more than full data.
		return
	}
	nodes := &serviceNodes{
		HealthyNodes:   convertNodes(result.healthyEntries),
		UnhealthyNodes: convertNodes(result.unhealthyEntries),
	}
	c.setLocked(serviceName, nodes)
}

// List gets service nodes from cache, including healthy and unhealthy ones.
func (c *cache) List(serviceName string) (*serviceNodes, error) {
	// Obtain from cache first.
	c.RLock()
	nodes, isExisting := c.nodesCache[serviceName]
	if isExisting {
		c.RUnlock()
		return nodes, nil
	}

	// Set up services that need attention.
	_, ok := c.watched[serviceName]
	c.RUnlock()
	if !ok {
		c.Lock()
		c.watched[serviceName] = true
		c.watcher.watchService(serviceName)
		c.Unlock()
	}
	return nil, nil
}

// watch for modifying the cache follows consul changes.
func (c *cache) watch() {
	go func() {
		watchResults := c.watcher.watch()
		for result := range watchResults {
			select {
			case <-c.exit:
				return
			default:
				c.update(result)
			}
		}
	}()
}

// stop stops the cache and stops the consul watch at the same time.
func (c *cache) stop() {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.exit:
		return
	default:
		close(c.exit)
	}
	c.watcher.stop()
}

// convertNodes converts consul node to trpc node.
func convertNodes(entries []*api.ServiceEntry) []*tregistry.Node {
	nodes := make([]*tregistry.Node, 0, len(entries))
	for _, s := range entries {
		meta := make(map[string]interface{})
		for k, v := range s.Service.Meta {
			meta[k] = v
		}
		node := &tregistry.Node{
			ServiceName: s.Service.ID,
			Address:     net.JoinHostPort(s.Service.Address, strconv.Itoa(s.Service.Port)),
			Metadata:    meta,
			Weight:      s.Service.Weights.Passing,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// newCache creates a new cache.
func newCache(options ...Option) (*cache, error) {
	watcher, err := newConsulWatcher(options...)
	if err != nil {
		return nil, err
	}
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	c := &cache{
		opts:       opts,
		watched:    make(map[string]bool),
		nodesCache: make(map[string]*serviceNodes),
		exit:       make(chan bool),
		watcher:    watcher,
	}
	c.watch()
	return c, nil
}
