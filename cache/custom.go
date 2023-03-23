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
)

// ClientAPIs consolidates the APIs that a custom cache client must implement.
type ClientAPIs interface {
	// IGet is an interface for all cache clients that support Get operations.
	IGet
	// ISet is an interface for all cache clients that support Set operations.
	ISet
}

// CustomConfig is a configuration struct for a custom cache client.
type CustomConfig struct {
	Client ClientAPIs
}

func newCustom(cfg *CustomConfig) (*Client, error) {
	if cfg == nil {
		return nil, errors.Errorf("nil ptr passed in for custom cache config")
	}

	return &Client{
		GetAPI: cfg.Client,
		SetAPI: cfg.Client,
	}, nil
}
