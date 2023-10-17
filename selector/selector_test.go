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

package selector

import (
	"reflect"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/glycerine/goconvey/convey"
	"github.com/golang/mock/gomock"
	tdiscovery "trpc.group/trpc-go/trpc-go/naming/discovery"
	"trpc.group/trpc-go/trpc-go/naming/loadbalance"
	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-naming-consul/discovery"
)

func TestSelector_Select(t *testing.T) {
	Convey("Select寻址", t, func() {
		s := New(WithLoadBalancer("random"))
		So(s, ShouldNotBeNil)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		d, err := discovery.New()
		So(err, ShouldNotBeNil)
		discovery.DefaultDiscovery = d

		ApplyMethod(reflect.TypeOf(d), "List", func(d *discovery.Discovery,
			service string, opt ...tdiscovery.Option) (nodes []*registry.Node, err error) {
			return nil, nil
		}).ApplyFunc(loadbalance.Get, func(name string) loadbalance.LoadBalancer {
			return &loadbalance.Random{}
		})

		node, err := s.Select("service")
		_ = s.Report(node, time.Second, err)
	})
}
