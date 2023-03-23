package heimdall

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/stretchr/testify/assert"
)

func TestToggleCache(t *testing.T) {
	isCacheEnabled := true
	ToggleCache(isCacheEnabled)
	assert.Equal(t, isCacheEnabled, skipCache)

	isCacheEnabled = false
	ToggleCache(isCacheEnabled)
	assert.Equal(t, isCacheEnabled, skipCache)
}

func TestGetData(t *testing.T) {
	ctx := context.Background()
	cacheKey := "cacheKey"
	softTTL := 1 * time.Second
	hardTTL := 2 * time.Second

	rpcCall := func() (*string, error) {
		resp := "response"
		return &resp, nil
	}
	rpcCallName := "rpcCallName"
	readFromCache := func() bool { return false }
	writeToCache := func(resp *string) bool { return true }

	res, err := getData(ctx, rpcCall, rpcCallName, cacheKey, softTTL, hardTTL, readFromCache, writeToCache)
	assert.NoError(t, err)
	assert.Equal(t, "response", *res)
}

func TestGenerateResp(t *testing.T) {
	cacheVal := &CacheValue{
		UpdatedTS: 123,
		SoftTTL:   1 * time.Second,
		Data:      `{"response": "response"}`,
	}

	actual := map[string]string{
		"response": "response",
	}

	res, err := generateResp[map[string]string](cacheVal)
	assert.NoError(t, err)
	assert.Equal(t, actual, *res)
}

func TestFetchFromCache(t *testing.T) {
	ctx := context.Background()
	key := "cacheKey"

	cacheVal := &CacheValue{
		UpdatedTS: 123,
		SoftTTL:   1 * time.Second,
		Data:      `{"response": "response"}`,
	}

	mockClient := &MockedCache{}
	cacheProvider = &cache.Client{
		GetAPI: mockClient,
		SetAPI: mockClient,
	}
	result, err := fetchFromCache(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, cacheVal, result)
}

func TestMakeCacheValue(t *testing.T) {
	type response struct {
		TestData string
	}

	testCases := []struct {
		desc             string
		resp             *response
		ttl              time.Duration
		expectedCacheVal *CacheValue
	}{
		{
			desc: "Test case 1: Should return CacheValue struct with correct Data and SoftTTL",
			resp: &response{TestData: "test"},
			ttl:  1 * time.Hour,
			expectedCacheVal: &CacheValue{
				Data:    `{"TestData":"test"}`,
				SoftTTL: 1 * time.Hour,
			},
		},
		{
			desc: "Test case 2: Should return CacheValue struct with correct Data and SoftTTL",
			resp: &response{TestData: "another test"},
			ttl:  2 * time.Hour,
			expectedCacheVal: &CacheValue{
				Data:    `{"TestData":"another test"}`,
				SoftTTL: 2 * time.Hour,
			},
		},
	}

	for _, tc := range testCases {
		cacheVal, err := makeCacheValue(tc.resp, tc.ttl)
		assert.Equal(t, cacheVal.Data, tc.expectedCacheVal.Data)
		assert.Equal(t, cacheVal.SoftTTL, tc.expectedCacheVal.SoftTTL)
		assert.Nil(t, err)
	}
}

func TestMakeCacheValue_ForInt64Response(t *testing.T) {
	data := int64(10)
	resp := &data
	softTTL := time.Second * 10

	cacheVal, _ := makeCacheValue(resp, softTTL)

	assert.Equal(t, softTTL, cacheVal.SoftTTL)

	var result int64
	err := json.Unmarshal([]byte(cacheVal.Data), &result)
	assert.Nil(t, err)
	assert.Equal(t, data, result)
}

func TestMakeCacheValue_ForFloat64Response(t *testing.T) {
	data := float64(10.1)
	resp := &data
	softTTL := time.Second * 10

	cacheVal, _ := makeCacheValue(resp, softTTL)

	assert.Equal(t, softTTL, cacheVal.SoftTTL)

	var result float64
	err := json.Unmarshal([]byte(cacheVal.Data), &result)
	assert.Nil(t, err)
	assert.Equal(t, data, result)
}

func TestHandleCacheHit(t *testing.T) {
	cacheKey := "cacheKey"
	softTTL := 1 * time.Second
	hardTTL := 2 * time.Second

	rpcCall := func() (*string, error) {
		resp := "first response"
		return &resp, nil
	}

	rpcCallName := "rpcCallName"
	readFromCache := func() bool { return true }
	writeToCache := func(resp *string) bool { return true }
	ctx := context.Background()

	client := &MockedCache{}
	cacheProvider = &cache.Client{
		GetAPI: client,
		SetAPI: client,
	}

	getData(ctx, rpcCall, rpcCallName, cacheKey, softTTL, hardTTL, readFromCache, writeToCache)
	time.Sleep(time.Second)
	getData(ctx, rpcCall, rpcCallName, cacheKey, softTTL, hardTTL, readFromCache, writeToCache)

	expectedGet := 2
	expectedSet := 1

	assert.Equal(t, expectedGet, client.NumGet)
	assert.Equal(t, expectedSet, client.NumSet)
}

type MockedCache struct {
	NumGet int
	NumSet int
}

func (m *MockedCache) Get(ctx context.Context, key string) ([]byte, error) {
	cacheVal := &CacheValue{
		UpdatedTS: 123,
		SoftTTL:   1 * time.Second,
		Data:      `{"response": "response"}`,
	}
	compressedVal, err := CompressStruct(ctx, cacheVal, constants.NoCompressionType)
	if err != nil {
		return nil, err
	}
	m.NumGet++
	return compressedVal, nil
}

func (m *MockedCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	m.NumSet++
	return nil
}
