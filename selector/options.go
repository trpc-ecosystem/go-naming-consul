// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package selector

// Options selector configuration
type Options struct {
	LoadBalancer string //load balancing strategy
}

// Option function for setting options.
type Option func(*Options)

// WithLoadBalancer sets load balancing policy
func WithLoadBalancer(loadBalancer string) Option {
	return func(options *Options) {
		if loadBalancer != "" {
			options.LoadBalancer = loadBalancer
		}

	}
}
