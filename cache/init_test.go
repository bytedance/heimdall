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

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/heimdall/constants"
	"github.com/stretchr/testify/assert"
)

var (
	sampleInvalidCustomCacheConfig = &Config{
		CacheProvider: constants.CustomCacheType,
	}
	sampleValidCustomCacheConfig = &Config{
		CacheProvider:       constants.CustomCacheType,
		CustomConfiguration: &CustomConfig{Client: &testCustomCache{}},
	}
)

func TestCacheInit(t *testing.T) {
	c, _ := newCustom(sampleValidCustomCacheConfig.CustomConfiguration)
	tests := []struct {
		name          string
		config        *Config
		validateError bool
		client        *Client
		freezeError   bool
	}{
		{
			name:          "invalid custom cache config",
			config:        sampleInvalidCustomCacheConfig,
			validateError: true,
			client:        nil,
			freezeError:   true,
		}, {
			name:          "valid custom cache config",
			config:        sampleValidCustomCacheConfig,
			validateError: false,
			client:        c,
			freezeError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.Equal(t, tt.validateError, err != nil)

			client, err := tt.config.Freeze()
			assert.Equal(t, tt.freezeError, err != nil)
			assert.Equal(t, tt.client, client)
		})
	}
}

type testCustomCache struct{}

func (c *testCustomCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
func (c *testCustomCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	return nil
}
