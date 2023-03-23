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

package constants

import "github.com/bytedance/heimdall/collections/set"

// CacheType is the type of cache to use for Heimdall
type CacheType int32

type CompressionLibraryType int32

const (
	// CustomCacheType allows the user to BYOC (Bring your own cache) as long as the user's cache client
	// implements the ClientAPIs interface
	CustomCacheType CacheType = iota + 1
	// RedisCacheType uses redis as a cache. It currently supports a single redis instance or a cluster.
	// Under the hood, it uses the go-redis library. The user has to pass in required attributes to connect
	// to redis itself. Supports redis 7.
	RedisCacheType
)

// RedisType is the type of redis server configuration the user is using.
type RedisType int32

const (
	// SingularRedisType is a default redis server.
	SingularRedisType RedisType = iota + 1
	// ClusterRedisType is a redis server that is set up in cluster mode.
	ClusterRedisType
)

var supportedCacheTypes = set.New[CacheType]().
	Add(RedisCacheType).
	Add(CustomCacheType)

const (
	// NoCompression will disable compression and uncompressed values are stored in the cache.
	NoCompressionType CompressionLibraryType = iota
	// GzipCompression will enable compression and values are compressed with the GZIP library and stored in the cache.
	GzipCompressionType
	// SnappyCompression will enable compression and values are compressed with the Snappy library and stored in the cache.
	SnappyCompressionType
)
