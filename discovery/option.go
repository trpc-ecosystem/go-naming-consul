// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package discovery

import "github.com/hashicorp/consul/api"

// Options service discovery configuration.
type Options struct {
	client *api.Client
}

// Option configuration function.
type Option func(options *Options)

// WithClient sets up consul client.
func WithClient(client *api.Client) Option {
	return func(options *Options) {
		options.client = client
	}
}
