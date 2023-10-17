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

package discovery

import (
	"testing"

	. "github.com/glycerine/goconvey/convey"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
)

func Test_cache_newCache(t *testing.T) {
	Convey("新建缓存", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		c, err := newCache(WithClient(client))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
	})
}

func Test_cache_List(t *testing.T) {
	Convey("通过缓存获取服务", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10

		c, err := newCache(WithClient(client))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		// Obtain it once first, which means paying attention to this service, and subsequent updates of this service will take effect.
		nodes, err := c.List("test")
		// Cache an empty cache first.
		err = c.cache("test", 1, &serviceNodes{HealthyNodes: emptyNodes, UnhealthyNodes: emptyNodes})
		So(err, ShouldBeNil)

		// Since the node is 0, an error should be reported at this time
		nodes, err = c.List("test")
		So(err, ShouldBeNil)
		So(len(nodes.HealthyNodes), ShouldEqual, 0)
		// Manually update the cache
		c.update(nil)
		nodes, err = c.List("test")
		So(err, ShouldBeNil)
		So(len(nodes.HealthyNodes), ShouldEqual, 0)

		c.update(&watchResult{
			serviceName:    "test",
			Version:        2,
			healthyEntries: []*api.ServiceEntry{tmp},
		})
		// At this time, the cache can be obtained.
		nodes, err = c.List("test")
		So(err, ShouldBeNil)
		So(len(nodes.HealthyNodes), ShouldNotEqual, 0)
		c.update(&watchResult{
			serviceName:    "test",
			Version:        1,
			healthyEntries: []*api.ServiceEntry{tmp},
		})
		// At this time, the cache can be obtained.
		nodes, err = c.List("test")
		So(err, ShouldBeNil)
		So(len(nodes.HealthyNodes), ShouldNotEqual, 0)
	})
}

func Test_cache_stop(t *testing.T) {
	Convey("关闭缓存", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		c, err := newCache(WithClient(client))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		c.stop()
		c.stop()
	})
}

func Test_convertNodes(t *testing.T) {
	Convey("节点转换", t, func() {
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10
		nodes := convertNodes([]*api.ServiceEntry{tmp})
		So(len(nodes), ShouldEqual, 1)
		nodes = convertNodes([]*api.ServiceEntry{})
		So(len(nodes), ShouldEqual, 0)
	})
}

func Test_cache_watch(t *testing.T) {
	Convey("watch", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		c, err := newCache(WithClient(client))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		c.watch()
	})
}

func Test_cache_cache(t *testing.T) {
	Convey("cache", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10

		c, err := newCache(WithClient(client))
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		// Empty cache.
		err = c.cache("test", 2, &serviceNodes{HealthyNodes: nil, UnhealthyNodes: nil})
		So(err, ShouldBeNil)
		// The cache has not been obtained and does not take effect.
		_ = c.cache("test", 2, &serviceNodes{
			HealthyNodes:   convertNodes([]*api.ServiceEntry{tmp}),
			UnhealthyNodes: convertNodes([]*api.ServiceEntry{tmp}),
		})
		nodes, err := c.List("test")
		So(nodes, ShouldBeNil)
		So(err, ShouldBeNil)

		// Obtain it once first, which means paying attention to this service, and subsequent updates of this service will take effect.
		_, _ = c.List("test")
		// Cache an empty cache first.
		err = c.cache("test", 2, &serviceNodes{
			HealthyNodes:   convertNodes([]*api.ServiceEntry{tmp}),
			UnhealthyNodes: convertNodes([]*api.ServiceEntry{tmp}),
		})
		So(err, ShouldBeNil)
		nodes, _ = c.List("test")
		So(len(nodes.HealthyNodes), ShouldEqual, 1)
	})
}
