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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleValidCustomMetricConfig = &Config{
		MetricsProvider: CustomMetricsType,
		CustomConfiguration: &CustomConfig{
			Client: &testCustomMetrics{},
		},
	}
	sampleInvalidCustomMetricConfig = &Config{
		MetricsProvider: MetricsType(-1),
	}
)

func TestMetricInit(t *testing.T) {
	c, _ := newCustom(sampleValidCustomMetricConfig.CustomConfiguration)
	tests := []struct {
		name          string
		config        *Config
		validateError bool
		client        *Client
		freezeError   bool
	}{
		{
			name:          "invalid custom metrics config",
			config:        sampleInvalidCustomMetricConfig,
			validateError: true,
			client:        nil,
			freezeError:   true,
		}, {
			name:          "valid custom cache config",
			config:        sampleValidCustomMetricConfig,
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

type testCustomMetrics struct{}

// IncreaseCacheHitMetric increases the cache hit metric.
func (c *testCustomMetrics) IncreaseCacheHitMetric(ctx context.Context, metricName string) {
}

// IncreaseCacheMissMetric increases the cache miss metric.
func (c *testCustomMetrics) IncreaseCacheMissMetric(ctx context.Context, metricName string) {
}

// IncreaseCacheSoftHitMetric increases the cache soft hit metric.
func (c *testCustomMetrics) IncreaseCacheSoftHitMetric(ctx context.Context, metricName string) {
}
