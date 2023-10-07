// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package error

import "errors"

var (
	// ServerNotAvailableError service unavailable, no nodes available.
	ServerNotAvailableError = errors.New("server can not available")
	// BalancerNotExistError there is no corresponding load balancing strategy.
	BalancerNotExistError = errors.New("load balancer is not exist")
)
