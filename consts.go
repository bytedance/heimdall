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
	"time"

	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/metrics"
)

var (
	defaultSoftTTL time.Duration
	defaultHardTTL time.Duration

	cacheProvider   *cache.Client
	metricsProvider *metrics.Client

	skipCache bool

	compressionLibrary constants.CompressionLibraryType

	version string
)

func InjectCacheProvider(c *cache.Client) {
	cacheProvider = c
}

func InjectMetricsProvider(m *metrics.Client) {
	metricsProvider = m
}

func InjectSkipCache(b bool) {
	skipCache = b
}

func InjectSoftTTL(softTTL time.Duration) {
	defaultSoftTTL = softTTL
}
func InjectHardTTL(hardTTL time.Duration) {
	defaultHardTTL = hardTTL
}
func InjectCompressionLibrary(t constants.CompressionLibraryType) {
	compressionLibrary = t
}
func InjectVersion(v string) {
	version = v
}
