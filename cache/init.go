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
	"github.com/pkg/errors"

	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/helpers"
)

// Config is a configuration struct for the cache client. It allows the user to specify the type of cache they want to use
// as well as the configuration for that cache type.
type Config struct {
	// CacheProvider is the type of cache to use. These are constants in the cache package.
	CacheProvider constants.CacheType
	// CustomConfiguration is a configuration for a custom cache client. This field is required only if CacheProvider is set to
	// CustomCacheType.
	CustomConfiguration *CustomConfig
	// RedisConfiguration is a configuration for a redis cache client. This field is required only if CacheProvider is set to
	// RedisCacheType.
	RedisConfiguration *RedisConfig
}

// Validate validates the cache configuration.
func (c *Config) Validate() error {
	switch c.CacheProvider {
	case constants.CustomCacheType:
		if c.CustomConfiguration == nil {
			return helpers.TernaryOp(c.CustomConfiguration == nil, errors.Errorf("custom cache config is nil"), nil)
		}
	case constants.RedisCacheType:
		return helpers.TernaryOp(c.RedisConfiguration == nil, errors.Errorf("redis configuration is nil"), nil)
	default:
		return errors.Errorf("cache type is not supported")
	}

	return nil
}

// Freeze freezes the cache configuration and generates the respective cache clients.
func (c *Config) Freeze() (*Client, error) {
	if c == nil {
		return nil, errors.Errorf("config, is nil")
	}
	switch c.CacheProvider {
	case constants.CustomCacheType:
		return newCustom(c.CustomConfiguration)
	case constants.RedisCacheType:
		return newRedis(c.RedisConfiguration)
	default:
		return nil, errors.Errorf("cache type %d is not supported", c.CacheProvider)
	}
}
