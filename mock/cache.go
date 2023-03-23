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

package mock

import (
	"context"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/pkg/errors"

	heimdall "github.com/bytedance/heimdall"
	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/helpers"
)

var jsonAPI = json.Config{UseNumber: true}.Froze()

// Cache is a simple mocked cache where user provides the map of key value pairs and it will automatically wrap it with Heimdall's metadata.
// This mocked cache supports errors where if the value to a key is an error, it will return that error.
// This mocked cache also supports string or any other structure types in the value field.
func Cache(data map[string]any) {
	processedData := make(map[string]any, len(data))
	for k, v := range data {
		JSON, err := jsonAPI.MarshalToString(v)
		if err != nil {
			panic(err)
		}
		cacheVal := &heimdall.CacheValue{
			UpdatedTS: time.Now().Unix(),
			SoftTTL:   time.Hour * 24,
			Data:      JSON,
		}
		compressedDat, _ := heimdall.CompressStruct(context.Background(), cacheVal, constants.GzipCompressionType)
		processedData[k] = compressedDat
	}

	newMockedClient := &mockedCache{mockedData: processedData}
	heimdall.InjectCacheProvider(
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

	return mockedVal.([]byte), nil
}

func (m *mockedCache) Set(_ context.Context, key string, val any, _ time.Duration) error {
	m.mockedData[key] = val
	return nil
}
