/*
 *
 * Copyright 2015 gRPC authors.
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
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	heimdall "github.com/bytedance/heimdall"
	"github.com/bytedance/heimdall/cache"
	"github.com/bytedance/heimdall/constants"
	"github.com/bytedance/heimdall/metrics"
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

// InitHeimdallSingular is an example of initializing a Heimdall with a singular Redis.
func InitHeimdallSingular() {
	metricsClient := &MyMetricsClient{}
	heimdallConfig := &heimdall.Config{ // Configuration used to initialise Heimdall
		DefaultSoftTTL: time.Second * 10, // Global softTTL if one is not provided at the time of function call. (Refer to concepts and notes for more info)
		DefaultHardTTL: time.Second * 40, // Global hardTTL if one is not provided at the time of functon call. (Refer to concepts and notes for more info)
		CacheConfig: cache.Config{ // Cache configuration, this must be provided.
			CacheProvider: constants.RedisCacheType, // Choose cache provider as Redis
			RedisConfiguration: &cache.RedisConfig{ // Redis configuration is required if CacheProvider is set to Redis
				RedisServerType: constants.SingularRedisType, // Server type chosen as SingularRedisType
				SingularConfig: &redis.Options{ // go-redis Options (https://redis.uptrace.dev/guide/)
					Addr:     "localhost:6379",
					Password: "", // no password set
					DB:       0,  // use default DB
				},
			},
		},
		EnableMetricsEmission: true,  // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
		SkipCache:             false, // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
		MetricsConfig: &metrics.Config{
			MetricsProvider: metrics.CustomMetricsType,
			CustomConfiguration: &metrics.CustomConfig{
				Client: metricsClient,
			},
		},
		CompressionLibrary: constants.SnappyCompressionType, // Toggling compression library as GzipCompression;
		Version:            "v1.0.0",                        // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
	}

	if err := heimdall.Init(heimdallConfig); err != nil {
		panic(err) // panic if heimdall fails to initialise
	}
	fmt.Println("heimdall initialised successfully")
}

// InitHeimdallCluster is an example of initializing a Heimdall with a cluster of Redis.
func InitHeimdallCluster() {
	metricsClient := &MyMetricsClient{}
	heimdallConfig := &heimdall.Config{ // Configuration used to initialise Heimdall
		DefaultSoftTTL: time.Second * 10, // Global softTTL if one is not provided at the time of function call. (Refer to concepts and notes for more info)
		DefaultHardTTL: time.Second * 40, // Global hardTTL if one is not provided at the time of functon call. (Refer to concepts and notes for more info)
		CacheConfig: cache.Config{ // Cache configuration, this must be provided.
			CacheProvider: constants.RedisCacheType, // Choose cache provider as Redis
			RedisConfiguration: &cache.RedisConfig{ // Redis configuration is required if CacheProvider is set to Redis
				RedisServerType: constants.ClusterRedisType, // Server type chosen as ClusterRedisType
				ClusterConfig: &redis.ClusterOptions{ // go-redis Options (https://redis.uptrace.dev/guide/)
					Addrs:    []string{"localhost:7000", "localhost:7001", "localhost:7002"},
					Password: "", // no password set
				},
			},
		},
		EnableMetricsEmission: true, // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
		MetricsConfig: &metrics.Config{
			MetricsProvider: metrics.CustomMetricsType,
			CustomConfiguration: &metrics.CustomConfig{
				Client: metricsClient,
			},
		},
		SkipCache:          false,                           // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
		CompressionLibrary: constants.SnappyCompressionType, //Toggling compression library as GzipCompression;
		Version:            "v1.0.1",                        // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
	}

	if err := heimdall.Init(heimdallConfig); err != nil {
		panic(err) // panic if heimdall fails to initialise
	}
	fmt.Println("heimdall initialised successfully")
}

// CustomClient structs and examples of implementations to fulfil the CustomClient interface
type CustomClient struct {
	RedisClient *redis.Client
}

func (c *CustomClient) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (c *CustomClient) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	err := c.RedisClient.Set(ctx, key, val, 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		return err
	}
	return nil
}

// InitHeimdallCustom is an example of initializing a Heimdall with a custom cache, in this case it is still Redis.
func InitHeimdallCustom() {
	CustomClient := &CustomClient{}
	CustomClient.RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	metricsClient := &MyMetricsClient{}

	heimdallConfig := &heimdall.Config{ // Configuration used to initialise Heimdall
		DefaultSoftTTL: time.Second * 10, // Global softTTL if one is not provided at the time of function call. (Refer to concepts and notes for more info)
		DefaultHardTTL: time.Second * 40, // Global hardTTL if one is not provided at the time of functon call. (Refer to concepts and notes for more info)
		CacheConfig: cache.Config{ // Cache configuration, this must be provided.
			CacheProvider: constants.CustomCacheType, // Choose cache provider as CustomCacheType
			CustomConfiguration: &cache.CustomConfig{ // Custom configuration is required if CacheProvider is set to CustomConfiguration
				Client: CustomClient, // Custom client that implements Get and Set
			},
		},
		EnableMetricsEmission: true, // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
		MetricsConfig: &metrics.Config{
			MetricsProvider: metrics.CustomMetricsType,
			CustomConfiguration: &metrics.CustomConfig{
				Client: metricsClient,
			},
		},
		SkipCache: false,    // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
		Version:   "v1.0.2", // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
	}

	if err := heimdall.Init(heimdallConfig); err != nil {
		panic(err) // panic if heimdall fails to initialise
	}
	fmt.Println("heimdall initialised successfully")
}

// Examples of fulfilling the Metrics interface
type IIncreaseMetric interface {
	IncreaseCacheHitMetric(ctx context.Context, metricName string)
	IncreaseCacheMissMetric(ctx context.Context, metricName string)
	IncreaseCacheSoftHitMetric(ctx context.Context, metricName string)
}

type MyMetricsClient struct {
	// any fields needed for the client can be added here
}

func (c *MyMetricsClient) IncreaseCacheHitMetric(ctx context.Context, metricName string) {
	// Implementation for increasing cache hit metric
	fmt.Printf("Cache hit: %s\n", metricName)
}

func (c *MyMetricsClient) IncreaseCacheMissMetric(ctx context.Context, metricName string) {
	// Implementation for increasing cache miss metric
	fmt.Printf("Cache miss: %s\n", metricName)
}

func (c *MyMetricsClient) IncreaseCacheSoftHitMetric(ctx context.Context, metricName string) {
	// Implementation for increasing cache soft hit metric
	fmt.Printf("Cache soft hit: %s\n", metricName)
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Comment and uncomment the following lines according to the example you wish to test,

	InitHeimdallSingular()
	// InitHeimdallCustom()
	// InitHeimdallCluster()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	for {
		r, err := heimdall.GRPCCall(c.SayHello, ctx, &pb.HelloRequest{Name: *name})

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Message Received: %s", r.GetMessage())
		time.Sleep(2 * time.Second)
	}
}
