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

package registry

import "github.com/hashicorp/consul/api"

// ServiceOptions a struct for service registry configuration.
type ServiceOptions struct {
	Interval                       string            // The time period between two health checks.
	Timeout                        string            // Timeout.
	Path                           string            // Reporting method of http.
	TLSSkipVerify                  *bool             // Whether to verify the https certificate.
	Tags                           []string          // Tag.
	Meta                           map[string]string // Metadata.
	Weight                         int               // Weights.
	DeregisterCriticalServiceAfter string            // Log out of the critical service.
}

// Options is for registering the configuration class.
type Options struct {
	DefaultServiceOptions *ServiceOptions
	ServicesOptions       map[string]*ServiceOptions
	client                *api.Client
}

// Option function for setting options.
type Option func(*Options)

// WithTags sets up tag.
func WithTags(tags []string) Option {
	return func(options *Options) {
		if len(tags) > 0 {
			options.DefaultServiceOptions.Tags = tags
		}
	}
}

// WithWeight sets permissions.
func WithWeight(weight int) Option {
	return func(options *Options) {
		if weight == 0 {
			weight = 1
		}
		options.DefaultServiceOptions.Weight = weight
	}
}

// WithMeta sets metadata.
func WithMeta(meta map[string]string) Option {
	return func(options *Options) {
		if len(meta) > 0 {
			options.DefaultServiceOptions.Meta = meta
		}
	}
}

// WithInterval sets the time period between two health checks.
func WithInterval(interval string) Option {
	return func(options *Options) {
		if interval != "" {
			options.DefaultServiceOptions.Interval = interval
		}
	}
}

// WithTimeout sets timeout.
func WithTimeout(timeout string) Option {
	return func(options *Options) {
		if timeout != "" {
			options.DefaultServiceOptions.Timeout = timeout
		}
	}
}

// WithPath sets to use http method.
func WithPath(path string) Option {
	return func(options *Options) {
		if path != "" {
			options.DefaultServiceOptions.Path = path
		}
	}
}

// WithTLSSkipVerify is to decide whether to verify tls for https mode.
func WithTLSSkipVerify(tlsSkipVerify *bool) Option {
	return func(options *Options) {
		options.DefaultServiceOptions.TLSSkipVerify = tlsSkipVerify
	}
}

// WithDeRegisterCriticalServiceAfter service critical for the auto logout.
func WithDeRegisterCriticalServiceAfter(t string) Option {
	return func(options *Options) {
		if t != "" {
			options.DefaultServiceOptions.DeregisterCriticalServiceAfter = t
		}
	}
}

// WithClient sets a consul client.
func WithClient(client *api.Client) Option {
	return func(options *Options) {
		options.client = client
	}
}

// WithServicesOptions sets service configuration.
func WithServicesOptions(servicesOptions map[string]*ServiceOptions) Option {
	return func(options *Options) {
		options.ServicesOptions = servicesOptions
	}
}
