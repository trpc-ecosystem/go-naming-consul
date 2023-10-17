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
)

func TestDNSDiscovery_List(t *testing.T) {
	Convey("http发现", t, func() {
		d := &DNSDiscovery{}
		nodes, _ := d.List("test")
		So(len(nodes), ShouldEqual, 0)
	})
}
