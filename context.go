/**
 * Tencent is pleased to support the open source community by making Polaris available.
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

package grpcpolaris

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	polarisRequestLbHashKey = "polaris.balancer.request.hashKey"
	polarisRequestLbPolicy  = "polaris.balancer.request.lbPolicy"
)

// SetLbHashKey set request scope LbHashKey
func SetLbHashKey(ctx context.Context, key string) context.Context {
	_, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			polarisRequestLbHashKey: key,
		}))

		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, polarisRequestLbHashKey, key)
}

// SetLbPolicy set request scope LbPolicy
func SetLbPolicy(ctx context.Context, policy string) context.Context {
	_, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			polarisRequestLbPolicy: policy,
		}))

		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, polarisRequestLbPolicy, policy)
}
