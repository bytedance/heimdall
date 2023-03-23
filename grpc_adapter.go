/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package heimdall

import (
	"context"
	"time"

	"github.com/bytedance/heimdall/helpers"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// GRPCCall wraps a grpc call method with Heimdall. It uses the global default set hard and soft TTLs.
func GRPCCall[request, response any](grpcFunc func(ctx context.Context, req *request, opts ...grpc.CallOption) (*response, error), ctx context.Context, req *request, opts ...grpc.CallOption) (*response, error) {
	return GRPCCallWithTTL(grpcFunc, ctx, req, defaultSoftTTL, defaultHardTTL)
}

// GRPCCallWithTTL wraps a grpc call method with Heimdall. It uses user defined hard and soft TTLs.
func GRPCCallWithTTL[request, response any](grpcFunc func(ctx context.Context, req *request, opts ...grpc.CallOption) (*response, error), ctx context.Context, req *request, softTTL, hardTTL time.Duration, opts ...grpc.CallOption) (*response, error) {
	if grpcFunc == nil {
		return nil, errors.Errorf("grpcFunc is nil")
	}

	rpcCallName := helpers.GetFunctionName(grpcFunc)

	cacheKey, err := helpers.GenerateCacheKey(req, rpcCallName, softTTL, hardTTL, version)
	if err != nil {
		return nil, err
	}

	return getData(ctx, wrapGRPCCallFunc(grpcFunc, ctx, req, opts...), rpcCallName, cacheKey, softTTL,
		hardTTL, func() bool { return true }, func(resp *response) bool { return true })
}

func wrapGRPCCallFunc[request, response any](grpcFunc func(ctx context.Context, req *request, opts ...grpc.CallOption) (*response, error), ctx context.Context, req *request, opts ...grpc.CallOption) func() (*response, error) {
	return func() (*response, error) {
		return grpcFunc(ctx, req, opts...)
	}
}
