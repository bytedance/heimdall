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
	"fmt"
	"testing"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/helpers"
)

var (
	testReq = &TestRPCRequest{
		UserID: "2139219754375423",
	}

	testResp = &TestRPCResponse{
		UserName: "John Doe",
	}

	testConfig = &Config{
		DefaultSoftTTL: time.Second * 10,
		DefaultHardTTL: time.Second * 40,
		CacheConfig: cache.Config{
			CacheProvider: constants.CustomCacheType,
			CustomConfiguration: &cache.CustomConfig{
				Client: &mockedCache{},
			},
		},
		EnableMetricsEmission: false,
		SkipCache:             false,
	}
)

func TestGRPCCall(t *testing.T) {
	c := &TestRPCClient{}
	tests := []struct {
		name                  string
		grpcFunc              func(ctx context.Context, req *TestRPCRequest, opts ...grpc.CallOption) (*TestRPCResponse, error)
		cacheHit              bool
		afterSoftTTLThreshold bool
		req                   *TestRPCRequest
		resp                  *TestRPCResponse
		err                   bool
		version               string
	}{
		{
			name:                  "cache hit",
			grpcFunc:              c.TestRPCCall,
			cacheHit:              true,
			afterSoftTTLThreshold: false,
			req:                   testReq,
			resp:                  testResp,
			err:                   false,
			version:               "v1.0.0",
		}, {
			name:                  "cache soft hit",
			grpcFunc:              c.TestRPCCall,
			cacheHit:              true,
			afterSoftTTLThreshold: true,
			req:                   testReq,
			resp:                  testResp,
			err:                   false,
			version:               "v1.0.0",
		}, {
			name:                  "cache miss",
			grpcFunc:              c.TestRPCCall,
			cacheHit:              false,
			afterSoftTTLThreshold: false,
			req:                   testReq,
			resp:                  testResp,
			err:                   false,
			version:               "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(testConfig)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			assert.Equal(t, tt.err, err != nil)

			cacheVal := map[string]any{}
			if tt.cacheHit {
				cacheKey, _ := helpers.GenerateCacheKey(tt.req, helpers.GetFunctionName(tt.grpcFunc), defaultSoftTTL, defaultHardTTL, version)
				cacheVal[cacheKey] = testMakeCacheValue(tt.resp, tt.afterSoftTTLThreshold)
			}

			mockCache(cacheVal)
			got, err := GRPCCall(tt.grpcFunc, context.Background(), tt.req)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.resp, got)
		})
	}
}

type TestRPCRequest struct {
	UserID string
}

type TestRPCResponse struct {
	UserName string
}

type TestRPCClient struct{}

func (c *TestRPCClient) TestRPCCall(ctx context.Context, req *TestRPCRequest, opts ...grpc.CallOption) (*TestRPCResponse, error) {
	return testResp, nil
}

func testMakeCacheValue(resp any, afterSoftTTLThreshold bool) []byte {
	respStr, _ := json.MarshalString(resp)
	ctx := context.Background()
	if afterSoftTTLThreshold {
		cacheVal := &CacheValue{
			UpdatedTS: time.Now().Add(-10 * time.Second).Unix(),
			SoftTTL:   5 * time.Second,
			Data:      respStr,
		}
		compressedCacheVal, _ := CompressStruct(ctx, cacheVal, constants.GzipCompressionType)
		return compressedCacheVal
	}
	cacheVal := &CacheValue{
		UpdatedTS: time.Now().Unix(),
		SoftTTL:   5 * time.Second,
		Data:      respStr,
	}
	compressedCacheVal, _ := CompressStruct(ctx, cacheVal, constants.GzipCompressionType)
	return compressedCacheVal
}

func mockCache(data map[string]any) {
	newMockedClient := &mockedCache{mockedData: data}
	InjectCacheProvider(
		&cache.Client{
			GetAPI: newMockedClient,
			SetAPI: newMockedClient,
		},
	)
}

type mockedCache struct {
	mockedData map[string]any
}

func (m *mockedCache) Get(_ context.Context, key string) ([]byte, error) {
	mockedVal, ok := m.mockedData[key]
	if !ok {
		return nil, errors.Errorf("cannot find key in mocked cache: %s, cache: %s", key, helpers.DumpJSON(m.mockedData))
	}

	switch mockedVal := mockedVal.(type) {
	case error:
		return nil, mockedVal
	}

	return mockedVal.([]byte), nil
}

func (m *mockedCache) Set(_ context.Context, key string, val any, _ time.Duration) error {
	m.mockedData[key] = val
	return nil
}
