/**
 * Tencent is pleased to support the open source community by making CL5 available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	polaris "github.com/polarismesh/grpc-go-polaris"
	"github.com/polarismesh/grpc-go-polaris/examples/common/pb"
	"github.com/polarismesh/polaris-go/pkg/config"
)

const (
	listenPort = 16011
)

func main() {
	// grpc客户端连接获取
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := &config.ConfigurationImpl{
		Global: &config.GlobalConfigImpl{
			System: &config.SystemConfigImpl{
				DiscoverCluster: &config.ServerClusterConfigImpl{
					Namespace: "Polaris-test",
					Service:   "polaris.discover-test",
				},
			},
			ServerConnector: &config.ServerConnectorConfigImpl{
				Addresses: []string{"127.0.0.1:8081", "127.0.0.1:8081"},
			},
		},
		Consumer: &config.ConsumerConfigImpl{
			LocalCache: &config.LocalCacheConfigImpl{
				PersistDir: "/var/polaris",
			},
		},
	}
	targetAddress, err := polaris.BuildTarget("HelloWorld",
		polaris.WithClientNamespace("Test"),
		polaris.WithConfig(cfg))

	conn, err := grpc.DialContext(ctx, targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	echoClient := pb.NewEchoServerClient(conn)

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if nil != err {
			log.Printf("fail to parse request form: %v\n", err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		values := r.Form["value"]
		log.Printf("receive value is %s\n", values)
		var value string
		if len(values) > 0 {
			value = values[0]
		}
		resp, err := echoClient.Echo(ctx, &pb.EchoRequest{Value: value})
		log.Printf("send message, resp (%v), err(%v)", resp, err)
		if nil != err {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(resp.GetValue()))
	}
	http.HandleFunc("/echo", indexHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil); nil != err {
		log.Fatal(err)
	}

}
