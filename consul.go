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

// Package consul is a package for consul.
package consul

import (
	"net/http"
	"runtime"

	"github.com/hashicorp/consul/api"
	tregistry "trpc.group/trpc-go/trpc-go/naming/registry"
	tselector "trpc.group/trpc-go/trpc-go/naming/selector"
	"trpc.group/trpc-go/trpc-go/plugin"
	"trpc.group/trpc-go/trpc-naming-consul/discovery"
	"trpc.group/trpc-go/trpc-naming-consul/registry"
	"trpc.group/trpc-go/trpc-naming-consul/selector"
)

func init() {
	plugin.Register(pluginName, &Plugin{})
	s := &selector.Selector{}
	tselector.Register("consul", s)
}

const (
	pluginType = "naming"
	pluginName = "consul"
)

// Plugin structure.
type Plugin struct{}

// Type for plugin type.
func (p *Plugin) Type() string {
	return pluginType
}

// Setup for Setting up.
func (p *Plugin) Setup(name string, decoder plugin.Decoder) error {
	cfg := Config{}
	err := decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	clientConfig := api.DefaultNonPooledConfig()
	clientConfig.Address = cfg.Address
	runtime.SetFinalizer(clientConfig.Transport, func(tr *http.Transport) {
		tr.CloseIdleConnections()
	})

	c, err := api.NewClient(clientConfig)
	if err != nil {
		return err
	}
	servicesOptions := make(map[string]*registry.ServiceOptions, 0)
	for _, register := range cfg.ServicesRegister {
		servicesOptions[register.Service] = convertServiceRegister2ServiceOptions(&cfg, register)
	}

	// Set registry.
	opts := []registry.Option{
		registry.WithTimeout(cfg.Register.Timeout),
		registry.WithInterval(cfg.Register.Interval),
		registry.WithTLSSkipVerify(cfg.Register.TLSSkipVerify),
		registry.WithPath(cfg.Register.Path),
		registry.WithClient(c),
		registry.WithMeta(cfg.Register.Meta),
		registry.WithWeight(cfg.Register.Weight),
		registry.WithTags(cfg.Register.Tags),
		registry.WithDeRegisterCriticalServiceAfter(cfg.Register.DeregisterCriticalServiceAfter),
		registry.WithServicesOptions(servicesOptions),
	}
	registry.DefaultRegistry = registry.New(opts...)
	// Each service is registered separately.
	for _, service := range cfg.Services {
		tregistry.Register(service, registry.DefaultRegistry)
	}
	for _, register := range cfg.ServicesRegister {
		tregistry.Register(register.Service, registry.DefaultRegistry)
	}

	// Set select.
	opt := []selector.Option{
		selector.WithLoadBalancer(cfg.Selector.LoadBalancer),
	}
	selector.DefaultSelector = selector.New(opt...)
	tselector.Register(pluginName, selector.DefaultSelector)

	// Set discovery.
	adopts := []discovery.Option{
		discovery.WithClient(c),
	}
	discovery.DefaultDiscovery, err = discovery.New(adopts...)
	if err != nil {
		return err
	}
	return nil
}

// convertServiceRegister2ServiceOptions converts ServiceRegister to ServiceOptions
// and use the global configuration to overwrite the configuration that does not exist locally.
func convertServiceRegister2ServiceOptions(cfg *Config, serviceRegister *ServiceRegister) *registry.ServiceOptions {
	if serviceRegister.Interval == "" && cfg.Register.Interval != "" {
		serviceRegister.Interval = cfg.Register.Interval
	}
	if serviceRegister.Timeout == "" && cfg.Register.Timeout != "" {
		serviceRegister.Timeout = cfg.Register.Timeout
	}
	if serviceRegister.Path == "" && cfg.Register.Path != "" {
		serviceRegister.Path = cfg.Register.Path
	}
	if serviceRegister.TLSSkipVerify == nil && cfg.Register.TLSSkipVerify != nil {
		serviceRegister.TLSSkipVerify = cfg.Register.TLSSkipVerify
	}
	if len(serviceRegister.Tags) == 0 && len(cfg.Register.Tags) > 0 {
		serviceRegister.Tags = cfg.Register.Tags
	}
	if len(serviceRegister.Meta) == 0 && len(cfg.Register.Meta) > 0 {
		serviceRegister.Meta = cfg.Register.Meta
	}
	if serviceRegister.Weight == 0 && cfg.Register.Weight != 0 {
		serviceRegister.Weight = cfg.Register.Weight
	}
	if serviceRegister.DeregisterCriticalServiceAfter == "" && cfg.Register.DeregisterCriticalServiceAfter != "" {
		serviceRegister.DeregisterCriticalServiceAfter = cfg.Register.DeregisterCriticalServiceAfter
	}
	return &registry.ServiceOptions{
		Interval:                       serviceRegister.Interval,
		Timeout:                        serviceRegister.Timeout,
		Path:                           serviceRegister.Path,
		TLSSkipVerify:                  serviceRegister.TLSSkipVerify,
		Tags:                           serviceRegister.Tags,
		Meta:                           serviceRegister.Meta,
		Weight:                         serviceRegister.Weight,
		DeregisterCriticalServiceAfter: serviceRegister.DeregisterCriticalServiceAfter,
	}
}
