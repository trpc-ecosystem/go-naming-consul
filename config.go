// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package consul

// Register configuration.
type Register struct {
	Interval                       string            `json:"interval,omitempty" yaml:"interval,omitempty"`                                                   // The time period between two health checks.
	Timeout                        string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`                                                     // Timeout.
	Path                           string            `json:"http,omitempty" yaml:"http,omitempty"`                                                           // Path.
	TLSSkipVerify                  *bool             `json:"tls_skip,omitempty" yaml:"tls_skip,omitempty"`                                                   // Whether to verify the https certificate.
	Tags                           []string          `json:"tags,omitempty" yaml:"tags,omitempty"`                                                           // Tag.
	Meta                           map[string]string `json:"meta,omitempty" yaml:"meta,omitempty"`                                                           // Metadata.
	Weight                         int               `json:"weight,omitempty" yaml:"weight,omitempty"`                                                       // Weights.
	DeregisterCriticalServiceAfter string            `json:"deregister_critical_service_after,omitempty" yaml:"deregister_critical_service_after,omitempty"` // How long does it take to cancel registration after the service hangs up.Register configuration.Register configuration.
}

// Config component support.
type Config struct {
	Address          string             `json:"address,omitempty" yaml:"address,omitempty"`                     // Consul address, compatible with the old one.
	Services         []string           `json:"services,omitempty" yaml:"services,omitempty"`                   // Registration service required.
	Register         Register           `json:"register,omitempty" yaml:"register,omitempty"`                   // Global registration configuration.
	ServicesRegister []*ServiceRegister `json:"services_register,omitempty" yaml:"services_register,omitempty"` // ServiceRegister enables different configurations for different services.
	// Selector configuration.
	Selector struct {
		LoadBalancer string `json:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"` //load balancing strategy
	}
}

// ServiceRegister enables different configurations for different services.
type ServiceRegister struct {
	Service string `json:"service,omitempty" yaml:"service,omitempty"` // registration service required
	Register
}
