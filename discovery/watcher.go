// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package discovery

import (
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

var (
	emptyServiceEntry = make([]*api.ServiceEntry, 0)
)

// watchResult packages consul change notification.
type watchResult struct {
	serviceName      string
	Version          uint64
	healthyEntries   []*api.ServiceEntry
	unhealthyEntries []*api.ServiceEntry
}

// serviceWatcher watches service changes.
type serviceWatcher struct {
	serviceName string
	plan        *watch.Plan
	resultChan  chan *watchResult
}

// newServiceWatcher watches new service.
func newServiceWatcher(serviceName string, resultChan chan *watchResult) *serviceWatcher {
	return &serviceWatcher{
		serviceName: serviceName,
		resultChan:  resultChan,
	}
}

// serviceHandler handles consul service changes.
func (sw *serviceWatcher) serviceHandler(idx uint64, data interface{}) {
	entries, ok := data.([]*api.ServiceEntry)
	if !ok {
		return
	}
	if len(entries) == 0 {
		sw.resultChan <- &watchResult{
			serviceName:      sw.serviceName,
			Version:          idx,
			healthyEntries:   emptyServiceEntry,
			unhealthyEntries: emptyServiceEntry,
		}
		return
	}
	var healthEntries, unhealthyEntries []*api.ServiceEntry
	for _, service := range entries {
		if service.Checks.AggregatedStatus() == api.HealthPassing {
			healthEntries = append(healthEntries, service)
		} else {
			unhealthyEntries = append(unhealthyEntries, service)
		}
	}
	if len(healthEntries) == 0 {
		healthEntries = emptyServiceEntry
	}
	if len(unhealthyEntries) == 0 {
		unhealthyEntries = emptyServiceEntry
	}
	sw.resultChan <- &watchResult{
		serviceName:      sw.serviceName,
		Version:          idx,
		healthyEntries:   healthEntries,
		unhealthyEntries: unhealthyEntries,
	}
}

// consulWatcher watches consul changes.
type consulWatcher struct {
	opts           *Options
	serviceWatcher map[string]*serviceWatcher
	resultChan     chan *watchResult
	exit           chan bool
}

// newConsulWatcher for creating a new consul.
func newConsulWatcher(options ...Option) (*consulWatcher, error) {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	if opts.client == nil {
		return nil, errors.New("consul client can not be nil")
	}
	cw := &consulWatcher{
		opts:           opts,
		exit:           make(chan bool),
		serviceWatcher: make(map[string]*serviceWatcher),
		resultChan:     make(chan *watchResult),
	}
	return cw, nil
}

// watchService watches service changes.
func (cw *consulWatcher) watchService(serviceName string) {
	sw := newServiceWatcher(serviceName, cw.resultChan)
	wp, _ := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": serviceName,
	})
	wp.Handler = sw.serviceHandler
	go wp.RunWithClientAndHclog(cw.opts.client, nil)
	sw.plan = wp
	cw.serviceWatcher[serviceName] = sw
}

// watch returns consul changes.
func (cw *consulWatcher) watch() <-chan *watchResult {
	return cw.resultChan
}

// Stop listening to consul changes.
func (cw *consulWatcher) stop() {
	select {
	case <-cw.exit:
		return
	default:
		close(cw.exit)
		for _, watcher := range cw.serviceWatcher {
			if watcher.plan == nil {
				continue
			}
			watcher.plan.Stop()
		}
	}
}
