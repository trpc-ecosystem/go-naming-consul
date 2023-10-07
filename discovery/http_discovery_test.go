// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package discovery

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"

	. "github.com/agiledragon/gomonkey"
	. "github.com/glycerine/goconvey/convey"
)

var (
	client = getConsulClient()
)

// getConsulClient gets the consul client.
func getConsulClient() *api.Client {
	c, _ := api.NewClient(&api.Config{})
	health := c.Health()
	ApplyMethod(reflect.TypeOf(c), "Health", func(c *api.Client) *api.Health {
		return health
	}).ApplyMethod(reflect.TypeOf(health), "Service", func(h *api.Health, service, tag string,
		passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10

		tmp2 := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp2.Service.Meta = make(map[string]string)
		tmp2.Service.Meta["key"] = "value"
		tmp2.Service.Service = "test"
		tmp2.Service.ID = "2"
		tmp2.Service.Address = "8.8.8.8"
		tmp2.Service.Port = 1000
		tmp2.Service.Weights.Passing = 10
		tmp2.Checks = append(tmp2.Checks, &api.HealthCheck{
			Status:    api.HealthWarning,
			ServiceID: "2",
		})
		return []*api.ServiceEntry{tmp, tmp2}, &api.QueryMeta{LastIndex: 1}, nil
	}).ApplyMethod(reflect.TypeOf(health), "ServiceMultipleTags", func(h *api.Health, service string, tags []string,
		passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10
		return []*api.ServiceEntry{tmp}, &api.QueryMeta{LastIndex: 1}, nil
	})
	return c
}

func TestDiscovery_List(t *testing.T) {
	Convey("http发现", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		d, err := New(WithClient(client))
		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		nodes, err := d.List("test")
		So(err, ShouldBeNil)
		So(len(nodes), ShouldNotEqual, 0)
	})
}

func TestDiscovery_ListAll(t *testing.T) {
	Convey("http发现", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		d, err := New(WithClient(client))
		So(err, ShouldBeNil)
		So(d, ShouldNotBeNil)
		nodes, nodes2, err := d.ListAll("test")
		So(err, ShouldBeNil)
		So(len(nodes), ShouldNotEqual, 0)
		So(len(nodes2), ShouldNotEqual, 0)
	})
}
