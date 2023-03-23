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

package metrics

import (
	"github.com/pkg/errors"

	"github.com/bytedance/heimdall/collections/set"
	"github.com/bytedance/heimdall/helpers"
)

// MetricsType is the type of metrics provider that Heimdall will use to emit cache hit, soft hit, miss metrics.
type MetricsType int32

const (
	// CustomMetricsType allows the user to BYOM (Bring your own metrics) as long as the user's metrics client fullfils the interface.
	CustomMetricsType MetricsType = iota + 1
)

var supportedMetricsType = set.New[MetricsType]().
	Add(CustomMetricsType)

// Config is a configuration struct for the metrics client.
type Config struct {
	// MetricsProvider is the type of metrics to use. These are constants in the metrics package.
	MetricsProvider MetricsType
	// CustomConfiguration is a configuration for a custom metrics client. This field is required only if MetricsProvider is set to
	// CustomMetricsType.
	CustomConfiguration *CustomConfig
}

// Validate validates the metrics configuration.
func (c *Config) Validate() error {
	switch c.MetricsProvider {
	case CustomMetricsType:
		if c.CustomConfiguration == nil {
			return helpers.TernaryOp(c.CustomConfiguration == nil, errors.Errorf("custom metrics config is nil"), nil)
		}
	default:
		return errors.Errorf("metrics type is not supported")
	}

	return nil
}

// Freeze freezes the metrics configuration and generates the respective metrics client.
func (c *Config) Freeze() (*Client, error) {
	if c == nil {
		return nil, errors.Errorf("config, is nil")
	}
	switch c.MetricsProvider {
	case CustomMetricsType:
		return newCustom(c.CustomConfiguration)
	default:
		return nil, errors.Errorf("metric type %d is not supported", c.MetricsProvider)
	}
}
