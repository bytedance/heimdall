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
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/metrics"
)

var (
	initOnce sync.Once
)

// Config is the configuration for Heimdall. The following is a simple sample configuration:
/*
Sample config with Redis Singular cluster set

	&Config{
		DefaultSoftTTL: time.Second * 10,
		DefaultHardTTL: time.Second * 40,
		CacheConfig: cache.Config{
			CacheProvider: cache.RedisCacheType,
			RedisConfiguration: &cache.RedisConfig{
				RedisServerType: cache.SingularRedisType,
				SingularConfig: &redis.Options{
					Addr:     "localhost:6379",
					Password: "", // no password set
					DB:       0,  // use default DB
				},
			},
		},
		EnableMetricsEmission: false,
		SkipCache:             false,
	}
*/
type Config struct {
	// DefaultSoftTTL refers to the soft TTL for all cache entries if one is not specified.
	// DefaultSoftTTL must be less than or equal to DefaultHardTTL. All data that is fetched from cache
	// before the initial time of cache entry + softTTL will be returned. The data will then be asynchronously
	// refreshed in the background.
	DefaultSoftTTL time.Duration `json:"default_soft_ttl,omitempty" yaml:"default_soft_ttl,omitempty" xml:"default_soft_ttl,omitempty"`
	// DefaultHardTTL refers to the hard TTL for all cache entries if one is not specified. This TTL will
	// be the cache's own stale expiration time if the cache supports eviction based on TTL.
	DefaultHardTTL time.Duration `json:"default_hard_ttl,omitempty" yaml:"default_hard_ttl,omitempty" xml:"default_hard_ttl,omitempty"`
	// CacheConfig refers to the cache configuration. Based on different cache providers the user can choose,
	// the user might need to fill in different information to initialize the cache. This field is required.
	CacheConfig cache.Config `json:"cache_config,omitempty" yaml:"cache_config,omitempty" xml:"cache_config,omitempty"`

	// EnableMetricsEmission refers to whether to enable metrics emission. If this is set to true, the user must
	// pass in a valid metrics config.
	EnableMetricsEmission bool `json:"enable_metrics_emission,omitempty" yaml:"enable_metrics_emission,omitempty" xml:"enable_metrics_emission,omitempty"`
	// MetricsConfig is the configuration for metrics emission. This field is required if user wishes to enable
	// metrics emission
	MetricsConfig *metrics.Config `json:"metrics_config,omitempty" yaml:"metrics_config,omitempty" xml:"metrics_config,omitempty"`

	// SkipCache is a global toggle to enable or disable the cache. Useful for multiple different environments.
	SkipCache bool `json:"skip_cache,omitempty" yaml:"skip_cache,omitempty" xml:"skip_cache,omitempty"`

	// CompressLibrary is a toggle to enable or disable library compression of choice.
	// Currently we support GZIP and Snappy compression, both of which should be used in different use cases.
	CompressionLibrary constants.CompressionLibraryType `json:"compress_library,omitempty" yaml:"compress_library,omitempty" xml:"compress_library,omitempty"`

	// Version is the version of cache you wish to use. This will be appended to the key name.
	// If there are any upgrades, this prevents breaking changes as old keys will not be re-used
	Version string `json:"version,omitempty" yaml:"version,omitempty" xml:"version,omitempty"`
}

func (c *Config) freeze() error {
	if err := c.validate(); err != nil {
		return err
	}

	cacheProv, err := c.CacheConfig.Freeze()
	if err != nil {
		return err
	}

	var metricsProv *metrics.Client
	if c.EnableMetricsEmission && c.MetricsConfig != nil {
		if metricsProv, err = c.MetricsConfig.Freeze(); err != nil {
			return err
		}
	}

	InjectSoftTTL(c.DefaultSoftTTL)
	InjectHardTTL(c.DefaultHardTTL)
	InjectCacheProvider(cacheProv)
	InjectMetricsProvider(metricsProv)
	InjectSkipCache(c.SkipCache)
	InjectCompressionLibrary(c.CompressionLibrary)
	InjectVersion(c.Version)

	return nil
}

func (c *Config) validate() error {
	if c.SkipCache {
		return nil
	}

	if c.DefaultHardTTL == 0 {
		return errors.Errorf("hard ttl is not set, if you want to disable cache, set SkipCache to true")
	}

	if c.DefaultHardTTL < c.DefaultSoftTTL {
		return errors.Errorf("hard ttl is less than soft ttl, if you want to disable soft TTL behavior, set DefaultSoftTTL to HardTTL")
	}

	if err := c.CacheConfig.Validate(); err != nil {
		return err
	}

	// 0 - No compression
	// 1 - Gzip compression
	// 2 - Snappy compression

	if c.CompressionLibrary < 0 || c.CompressionLibrary > 3 {
		return errors.Errorf("invalid compression library type specified.")
	}
	return nil
}

// Init initializes Heimdall. This function must be called before any other functions, preferably during initialisation of your application.
func Init(cfg *Config) error {
	var err error
	initOnce.Do(func() {
		err = cfg.freeze()
	})
	return err
}
