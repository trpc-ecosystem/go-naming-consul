// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// Package registry 提供consul注册服务
package registry

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"trpc.group/trpc-go/trpc-go/naming/registry"

	. "github.com/glycerine/goconvey/convey"
)

func TestRegistry_Register(t *testing.T) {
	Convey("注册", t, func() {
		c, _ := api.NewClient(&api.Config{})
		r := New(WithTags([]string{"test"}),
			WithWeight(100),
			WithMeta(map[string]string{"dyeing": "false"}),
			WithInterval("10s"),
			WithTimeout("1s"),
			WithClient(c),
		)
		err := r.Register("testService", registry.WithAddress("test"))
		So(err, ShouldNotBeNil)
		err = r.Register("testService", registry.WithAddress("8.8.8.8:test"))
		So(err, ShouldNotBeNil)

		_ = r.Register("testService", registry.WithAddress("8.8.8.8:1000"))

		_ = r.Deregister("testService")

		_ = r.Deregister("testService1")
	})
}

func Test_genAgentServiceID(t *testing.T) {
	Convey("注册", t, func() {
		serviceID := genAgentServiceID("test", "127.0.0.1", "8080")
		So(serviceID, ShouldEqual, "test-127.0.0.1-8080")
	})
}
