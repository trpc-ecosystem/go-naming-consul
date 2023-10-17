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

func Test_newConsulWatcher(t *testing.T) {
	Convey("新建consul watcher", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		watcher, err := newConsulWatcher(WithClient(client))
		So(err, ShouldBeNil)
		So(watcher, ShouldNotBeNil)
	})
}

func Test_consulWatcher_watchService(t *testing.T) {
	Convey("关注服务", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		watcher, err := newConsulWatcher(WithClient(client))
		So(err, ShouldBeNil)
		So(watcher, ShouldNotBeNil)
		watcher.watchService("test")

		So(watcher.serviceWatcher["test"], ShouldNotBeNil)
	})
}

func Test_consulWatcher_stop(t *testing.T) {
	Convey("停止关注consul服务变更", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		watcher, err := newConsulWatcher(WithClient(client))
		So(err, ShouldBeNil)
		So(watcher, ShouldNotBeNil)

		watcher.stop()
		watcher.stop()
	})
}

func Test_consulWatcher_watch(t *testing.T) {
	Convey("测试watch变更chan", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		watcher, err := newConsulWatcher(WithClient(client))
		So(err, ShouldBeNil)
		So(watcher, ShouldNotBeNil)
		watcher.resultChan = make(chan *watchResult, 1)
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10
		watcher.resultChan <- &watchResult{
			serviceName:    "test",
			Version:        1,
			healthyEntries: []*api.ServiceEntry{tmp},
		}
		result := <-watcher.watch()
		So(result, ShouldNotBeNil)
	})
}

func Test_serviceWatcher_serviceHandler(t *testing.T) {
	Convey("测试server watch变更", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		watcher := newServiceWatcher("test", make(chan *watchResult, 1))
		So(watcher, ShouldNotBeNil)
		tmp := &api.ServiceEntry{Service: &api.AgentService{}}
		tmp.Service.Meta = make(map[string]string)
		tmp.Service.Meta["key"] = "value"
		tmp.Service.ID = "1"
		tmp.Service.Address = "8.8.8.8"
		tmp.Service.Port = 1000
		tmp.Service.Weights.Passing = 10
		watcher.serviceHandler(1, []*api.ServiceEntry{tmp})

		result := <-watcher.resultChan
		So(result, ShouldNotBeNil)

		watcher.serviceHandler(2, []*api.ServiceEntry{})

		result = <-watcher.resultChan
		So(result, ShouldNotBeNil)
		watcher.serviceHandler(2, nil)
	})

}
