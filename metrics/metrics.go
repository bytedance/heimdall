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
)

// Client is a metrics client.
type Client struct {
	IncreaseMetricAPI IIncreaseMetric
}

// IIncreaseMetric is an interface for a basic increase metrics client.
type IIncreaseMetric interface {
	IncreaseCacheHitMetric(ctx context.Context, metricName string)
	IncreaseCacheMissMetric(ctx context.Context, metricName string)
	IncreaseCacheSoftHitMetric(ctx context.Context, metricName string)
}

// IncreaseCacheHitMetric increases the cache hit metric.
func (c *Client) IncreaseCacheHitMetric(ctx context.Context, metricName string) {
	c.IncreaseMetricAPI.IncreaseCacheHitMetric(ctx, metricName)
}

// IncreaseCacheMissMetric increases the cache miss metric.
func (c *Client) IncreaseCacheMissMetric(ctx context.Context, metricName string) {
	c.IncreaseMetricAPI.IncreaseCacheMissMetric(ctx, metricName)
}

// IncreaseCacheSoftHitMetric increases the cache soft hit metric.
func (c *Client) IncreaseCacheSoftHitMetric(ctx context.Context, metricName string) {
	c.IncreaseMetricAPI.IncreaseCacheSoftHitMetric(ctx, metricName)
}
