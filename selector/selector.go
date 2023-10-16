//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Package selector is a package for selector.
package selector

import (
	"time"

	tdiscovery "trpc.group/trpc-go/trpc-go/naming/discovery"
	"trpc.group/trpc-go/trpc-go/naming/loadbalance"
	tregistry "trpc.group/trpc-go/trpc-go/naming/registry"
	tselector "trpc.group/trpc-go/trpc-go/naming/selector"
	"trpc.group/trpc-go/trpc-naming-consul/discovery"
	consul_error "trpc.group/trpc-go/trpc-naming-consul/error"
)

// Selector structure.
type Selector struct {
	Opts *Options // configuration
}

// DefaultSelector instantiated objects by Selector structure.
var DefaultSelector *Selector

// New instantiates selector.
func New(call ...Option) *Selector {
	s := &Selector{
		Opts: &Options{
			LoadBalancer: "random",
		},
	}

	for _, o := range call {
		o(s.Opts)
	}
	return s

}

// Select node.
func (s *Selector) Select(serviceName string, opts ...tselector.Option) (node *tregistry.Node, err error) {
	o := &tselector.Options{}
	for _, opt := range opts {
		opt(o)
	}
	d := discovery.DefaultDiscovery
	nodes, err := d.List(serviceName, tdiscovery.WithContext(o.Ctx))
	if err != nil {
		return nil, err
	}

	var loadBalanceType string
	if o.LoadBalanceType != "" {
		loadBalanceType = o.LoadBalanceType
	} else {
		loadBalanceType = s.Opts.LoadBalancer
	}

	load := loadbalance.Get(loadBalanceType)
	if load == nil {
		return nil, consul_error.BalancerNotExistError
	}

	loadBalanceOpts := []loadbalance.Option{
		loadbalance.WithContext(o.Ctx),
		loadbalance.WithLoadBalanceType(loadBalanceType),
		loadbalance.WithKey(o.Key),
		loadbalance.WithNamespace(o.Namespace),
	}
	return load.Select(serviceName, nodes, loadBalanceOpts...)
}

// Report for reporting the information.
func (s *Selector) Report(node *tregistry.Node, cost time.Duration, err error) error {
	return nil
}
