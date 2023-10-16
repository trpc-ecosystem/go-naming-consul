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

// Package registry 提供consul注册服务
package registry

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/hashicorp/consul/api"
	"trpc.group/trpc-go/trpc-go/naming/registry"
)

// Registry consul is for registering service implementation.
type Registry struct {
	opts *Options
}

// DefaultRegistry instantiated objects by Registry structure.
var DefaultRegistry *Registry

var serviceIDMap sync.Map

// New Register via http method.
func New(opts ...Option) *Registry {
	r := &Registry{
		opts: &Options{
			DefaultServiceOptions: &ServiceOptions{
				Timeout:  "10s",
				Interval: "20s",
			},
		},
	}
	// Configure.
	for _, o := range opts {
		o(r.opts)
	}
	return r
}

// Register for registering the service to the Consul instance.
func (r *Registry) Register(service string, opts ...registry.Option) error {

	options := &registry.Options{}
	for _, opt := range opts {
		opt(options)
	}

	address := options.Address
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	pt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	serviceOptions := r.opts.DefaultServiceOptions
	if existServiceOptions, ok := r.opts.ServicesOptions[service]; ok {
		serviceOptions = existServiceOptions
	}
	var TLSSkipVerify bool
	if serviceOptions.TLSSkipVerify != nil {
		TLSSkipVerify = *serviceOptions.TLSSkipVerify
	}
	// todo: http and grpc methods for registration.
	check := &api.AgentServiceCheck{
		Interval:                       serviceOptions.Interval,
		Timeout:                        serviceOptions.Timeout,
		TCP:                            fmt.Sprintf("%s:%s", host, port),
		TLSSkipVerify:                  TLSSkipVerify,
		DeregisterCriticalServiceAfter: serviceOptions.DeregisterCriticalServiceAfter,
	}
	return r.opts.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Kind:    api.ServiceKindTypical,
		ID:      genAgentServiceID(service, host, port),
		Name:    service,
		Port:    pt,
		Address: host,
		Check:   check,
		Tags:    serviceOptions.Tags,
		Meta:    serviceOptions.Meta,
		Weights: &api.AgentWeights{
			Passing: serviceOptions.Weight,
			Warning: serviceOptions.Weight,
		},
	})
}

// Deregister for unregistering service.
func (r *Registry) Deregister(service string) error {
	serviceID, ok := serviceIDMap.Load(service)
	if !ok {
		return nil
	}
	err := r.opts.client.Agent().ServiceDeregister(serviceID.(string))
	if err != nil {
		return err
	}
	serviceIDMap.Delete(service)
	return nil
}

// genAgentServiceID for constructing and generating a service instance name to prevent duplicate names.
func genAgentServiceID(service string, host string, port string) string {
	serviceID := service + "-" + host + "-" + port
	serviceIDMap.Store(service, serviceID)
	return serviceID
}
