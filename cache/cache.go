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
	"time"

	"github.com/pkg/errors"

	"github.com/bytedance/heimdall/constants"
)

// Client is a generic client structure for caches. All caches must fullfil the APIs that Client provides.
// Client is structured in this way to allow for easy extension of the cache package as well as ease of mocking.
type Client struct {
	// GetAPI is any cache client that can get items from the cache.
	GetAPI IGet
	// SetAPI is any cache client that can set items
	SetAPI ISet
}

var CompressionLibrary constants.CompressionLibraryType

// IGet is an interface for all cache clients that support Get operations.
type IGet interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

// ISet is an interface for all cache clients that support Set operations.
type ISet interface {
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
}

// Get simply gets an item from the cache based on the API provided by the cache client.
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	compressedData, err := c.GetAPI.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pull from cache")
	}

	return compressedData, nil
}

// Set simply sets an item in the cache based on the API provided by the cache client.
func (c *Client) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	err := c.SetAPI.Set(ctx, key, val, ttl)
	if err != nil {
		return errors.Wrap(err, "unable to set into cache")
	}
	return nil
}
