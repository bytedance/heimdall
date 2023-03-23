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

	json "github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

var (
	jsonAPI = json.Config{UseNumber: true}.Froze() // used to maintain int64 and float64 precision when unmarshalling into an interface{}
)

// ToggleCache is a utility function that can be called to toggle caching on / off.
// Event listeners can be configured to toggle the cache on / off with a simple configuration change.
func ToggleCache(isCacheEnabled bool) {
	skipCache = isCacheEnabled
}

type CacheValue struct {
	UpdatedTS int64
	SoftTTL   time.Duration
	Data      string
}

func getData[response any](
	ctx context.Context,
	rpcCall func() (*response, error),
	rpcCallName string,
	cacheKey string,
	softTTL time.Duration,
	hardTTL time.Duration,
	readFromCache func() bool,
	writeToCache func(*response) bool,
) (
	res *response,
	err error,
) {
	if isSkipCache() {
		return rpcCall()
	}

	var result *CacheValue
	if !readFromCache() {
		result, err = handleCacheMiss(ctx, cacheKey, rpcCall, softTTL, hardTTL, rpcCallName, writeToCache)
		if err != nil {
			return nil, err
		}
		return generateResp[response](result)
	}

	result, err = fetchFromCache(ctx, cacheKey)
	if err != nil {
		result, err = handleCacheMiss(ctx, cacheKey, rpcCall, softTTL, hardTTL, rpcCallName, writeToCache)
		if err != nil {
			return nil, err
		}
	}

	if isPastSoftTTLThreshhold(result) {
		handleCacheSoftHit(ctx, cacheKey, rpcCall, softTTL, hardTTL, rpcCallName, writeToCache)
	}

	handleCacheHit(ctx, rpcCallName)

	return generateResp[response](result)
}

func generateResp[response any](cacheVal *CacheValue) (res *response, err error) {
	resp := new(response)
	err = generateResponseStructFromCacheVal(cacheVal, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func fetchFromCache(ctx context.Context, key string) (*CacheValue, error) {
	cacheVal := &CacheValue{}
	val, err := cacheProvider.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	err = DecompressStruct(ctx, val, cacheVal, compressionLibrary)
	return cacheVal, err
}

func handleCacheMiss[response any](ctx context.Context, key string, rpcCall func() (*response, error), softTTL,
	hardTTL time.Duration, rpcCallName string, writeToCache func(*response) bool) (*CacheValue, error) {
	if !isSkipMetrics() {
		metricsProvider.IncreaseCacheMissMetric(ctx, rpcCallName)
	}

	resp, err := rpcCall()
	if err != nil {
		return nil, errors.Wrap(err, "rpc call failed")
	}

	go updateCache(ctx, key, resp, softTTL, hardTTL, writeToCache)
	return makeCacheValue(resp, softTTL)
}

func handleCacheSoftHit[response any](ctx context.Context, key string, rpcCall func() (*response, error), softTTL,
	hardTTL time.Duration, rpcCallName string, writeToCache func(*response) bool) {
	if !isSkipMetrics() {
		metricsProvider.IncreaseCacheSoftHitMetric(ctx, rpcCallName)
	}

	go func() {
		resp, err := rpcCall()
		if err != nil {
			return // don't write to cache on error
		}

		updateCache(ctx, key, resp, softTTL, hardTTL, writeToCache)
	}()
}

func handleCacheHit(ctx context.Context, rpcCallName string) {
	if !isSkipMetrics() {
		metricsProvider.IncreaseCacheHitMetric(ctx, rpcCallName)
	}
}

func generateResponseStructFromCacheVal[response any](cacheVal *CacheValue, res *response) error {
	err := jsonAPI.UnmarshalFromString(cacheVal.Data, res)
	if err != nil {
		return errors.Wrap(err, "unable to marshal cache value to rpc response")
	}
	return nil
}

func updateCache[response any](ctx context.Context, key string, rpcCallResp *response,
	softTTL, hardTTL time.Duration, writeToCache func(*response) bool) {
	if rpcCallResp == nil || !writeToCache(rpcCallResp) {
		return
	}
	cacheVal, err := makeCacheValue(rpcCallResp, softTTL)
	if err != nil {
		return
	}
	compressedData, err := CompressStruct(ctx, cacheVal, compressionLibrary)
	if err != nil {
		return
	}
	err = cacheProvider.Set(ctx, key, compressedData, hardTTL)
	if err != nil {
		return
	}
}

func isPastSoftTTLThreshhold(cacheVal *CacheValue) bool {
	return cacheVal.UpdatedTS+int64(cacheVal.SoftTTL.Seconds()) < time.Now().Unix()
}

func isSkipCache() bool {
	return skipCache
}

func isSkipMetrics() bool {
	return metricsProvider == nil
}

func makeCacheValue(val any, softTTL time.Duration) (*CacheValue, error) {
	data, err := jsonAPI.MarshalToString(val)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal rpc response")
	}
	return &CacheValue{
		UpdatedTS: time.Now().Unix(),
		Data:      data,
		SoftTTL:   softTTL,
	}, nil
}
