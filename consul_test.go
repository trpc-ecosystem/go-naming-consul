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

package consul

import (
	"testing"

	. "github.com/glycerine/goconvey/convey"
	"github.com/stretchr/testify/require"
	trpc "trpc.group/trpc-go/trpc-go"
	// register http codec to avoid panic when calling trpc.NewServer() without stub code
	_ "trpc.group/trpc-go/trpc-go/http"
)

func TestPlugin_Setup(t *testing.T) {
	require.NotNil(t, trpc.NewServer())
}

func Test_convertServiceRegister2ServiceOptions(t *testing.T) {
	Convey("测试配置转换", t, func() {
		// Use custom configuration.
		verify := true
		options := convertServiceRegister2ServiceOptions(&Config{
			Address:  "127.0.0.1:8080",
			Services: []string{"test"},
			Register: Register{
				Interval:                       "1s",
				Timeout:                        "1s",
				Path:                           "test",
				TLSSkipVerify:                  &verify,
				Tags:                           []string{"test"},
				Meta:                           map[string]string{"key": "value"},
				Weight:                         10,
				DeregisterCriticalServiceAfter: "10m",
			},
		}, &ServiceRegister{
			Service: "real",
			Register: Register{
				Interval:                       "2s",
				Timeout:                        "2s",
				Path:                           "test",
				Tags:                           []string{"test2", "test2"},
				Meta:                           map[string]string{"key1": "value1", "key2": "value2"},
				Weight:                         100,
				DeregisterCriticalServiceAfter: "1m",
			},
		})
		So(options.Interval, ShouldEqual, "2s")
		So(options.Timeout, ShouldEqual, "2s")
		So(len(options.Tags), ShouldEqual, 2)
		So(len(options.Meta), ShouldEqual, 2)
		So(options.Weight, ShouldEqual, 100)
		So(options.DeregisterCriticalServiceAfter, ShouldEqual, "1m")

		// Use global configuration.
		options = convertServiceRegister2ServiceOptions(&Config{
			Address:  "127.0.0.1:8080",
			Services: []string{"test"},
			Register: Register{
				Interval:                       "1s",
				Timeout:                        "1s",
				Path:                           "test",
				TLSSkipVerify:                  &verify,
				Tags:                           []string{"test"},
				Meta:                           map[string]string{"key": "value"},
				Weight:                         10,
				DeregisterCriticalServiceAfter: "10m",
			},
		}, &ServiceRegister{
			Service:  "real",
			Register: Register{},
		})
		So(options.Interval, ShouldEqual, "1s")
		So(options.Timeout, ShouldEqual, "1s")
		So(len(options.Tags), ShouldEqual, 1)
		So(len(options.Meta), ShouldEqual, 1)
		So(options.Weight, ShouldEqual, 10)
		So(options.DeregisterCriticalServiceAfter, ShouldEqual, "10m")
	})
}
