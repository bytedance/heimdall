# Heimdall - A simple RPC Caching SDK

**Description**: 
Heimdall is a powerful caching middleware to cache downstream RPC calls. It is designed to adopt a simple cache-aside strategy and employs a soft TTL and hard TTL approach to enhance overall fault tolerance.

By wrapping your RPC call with the SDK's function, Heimdall can act as a middleware and cache your calls, thus improving the performance and reliability of your application.

## Why use Heimdall
Front-end requests may require multiple RPC calls to different services. When the same RPC calls are frequently reused or the data fetched from the RPC calls does not go stale quickly, using Heimdall as a wrapper for the RPC calls to cache downstream calls can significantly reduce latency.

By adopting a simple cache-aside strategy and employing a soft TTL and hard TTL approach, Heimdall provides fault tolerance and high availability, ensuring that cached data is always up-to-date.

Moreover, Heimdall has been battle-tested in ByteDance's production environment, where it delivered remarkable improvements in performance and reliability. 

## Heimdall usage in ByteDance
Heimdall is a caching middleware for RPC calls at ByteDance, providing remarkable performance and reliability enhancements. Our empirical data shows impressive results: we have achieved a 34.4% reduction in latency on all APIs, while P50 has seen a latency reduction of 50%. Additionally, the average latency for all APIs has decreased by 44.4%. Furthermore, we have been able to reduce peak QPS for an internal service, resulting in a QPS reduction of 70.9%.

These results demonstrate the effectiveness of Heimdall in improving the performance and reliability of downstream calls, making it an essential component of our tech stack.
## When to use Heimdall
Heimdall can very simply be used to to cache downstream calls, particularly the ones with data that does not change very frequently (e.g static metadata).

An important metric to measure to ascertain whether or not to cache a downstream request is if there is a good [cache hit ratio](https://www.cloudflare.com/en-gb/learning/cdn/what-is-a-cache-hit-ratio/) across requests. That is, results of calls to the dependency can be used across multiple requests or operations. If each request typically requires a unique query to the dependent service with unique-per-request results, then a cache would have a negligible hit rate and the cache does no good. A second consideration is how tolerant a team’s service and its clients are to eventual consistency. Cached data necessarily grows inconsistent with the source over time, so caching can only be successful if both the service and its clients compensate accordingly. The rate of change of the source data, as well as the cache policy for refreshing data, will determine how inconsistent the data tends to be. In our experience, we've found setting a shorter softTTL and hardTTL for cases whereby data goes stale quicker helpful.

We need to ensure that the underlying service is resilient in the face of cache non-availability, which includes a variety of circumstances that lead to the inability to serve requests using cached data. These include cold starts, caching fleet outages, changes in traffic patterns, or extended downstream outages. In many cases, this could mean trading some of your availability to ensure that your servers and your dependent services don’t brown out (for example by shedding load, capping requests to dependent services, or serving stale data). Run load tests with caches disabled to validate this.

## Concepts and Notes

### Cache 
Heimdall allows user to setup a Redis Instance or Redis Cluster as their preferred cache provider or their own custom cache implementation. Redis client is supported through the use of [go-redis](https://github.com/go-redis/redis)

### Metrics
Heimdall supports emission of metrics. However the user must provide their own metrics implementation.

### Compression
Heimdall supports the use of no compression, GZIP compression and Snappy compression under the hood to reduce space used for cache storage. All of these can be set under the CompressionLibrary attribute when initialising Heimdall. It is important to choose the appropriate compression library for your application. If your data is accessed frequently, it is better to use Snappy compression that has a smaller compression ratio but with faster performance. If your data is accessed less frequently and the space used is large, it might be better to use the GZip compression library with higher compression ratio. Otherwise, it is also wise to not use any form of compression to reduce overhead if memory usage is not a concern.

### Time To Live (TTL)
Heimdall supports two types of TTLs by default, known as SoftTTL and HardTTL.

* HardTTL specifies the time after which the data becomes stale and is invalidated. This is calculated as the sum of the initial data cache time and the HardTTL, which determines the maximum allowable staleness of the data.

* SoftTTL, on the other hand, is a value that falls between the initial data cache time and the HardTTL. Data that is accessed within the SoftTTL period will be returned from the cache, while data accessed between the SoftTTL and HardTTL periods will be served from the cache, but asynchronously updated in the background to refresh the data that has become tolerably stale. It's worth noting that SoftTTL can be disabled by setting its value equal to that of HardTTL.

With Heimdall's TTL-based caching strategy, you can ensure that your application always serves fresh and up-to-date data to your users, while still delivering optimal performance and reducing the load on your backend services.

### Key Eviction Policy
Heimdall does not automatically handle key eviction. Please configure redis with your own key eviction policy, we've used `allkeys-lru` and found that it worked pretty well.

### int64 and float64 data types
Marshalling and unmarshalling of interface{} objects that represent int64 and float64 data types can incur a loss of precision. Please enforce the types in the request and response structs with the specific data types.

## Data redundancy and Fault tolerance

### Data redundancy
In the event the Redis instance starts to hit its memory limit, Redis will automatically evict keys with the `allkeys-lru` strategy. 
[Overview of Redis key eviction policies](https://redis.io/docs/reference/eviction/#:~:text=allkeys-lru%3A%20Keeps%20most%20recently,expire%20field%20set%20to%20true%20)

### Fault Tolerance
* Heimdall provides exceptional fault tolerance. In the event of a Redis cluster outage, our caching solution ensures that there is no significant impact on the overall service, apart from a slight increase in API latency. Our services continue to function normally, ensuring smooth operations for our users.

* Similarly, if there is a temporary downtime in the RPC service, our caching solution ensures that users will not immediately experience any issues, as data is still available from the cache. This ensures uninterrupted service and enhanced user experience.

## Dependencies
* Go1.18+
* A Redis instance (Redis 7) (https://redis.io/) or a custom cache implementation.
* A metrics provider. (Optional)
* Linux/MacOS/Windows
* x86/ARM
* gRPC, more coming soon!

## Installation
Simply run:
```bash
go get github.com/bytedance/heimdall@latest
```

## Usage
Heimdall supports an advanced form of usage as well as simple form of usage. Both usages are provided below with examples.
Initialize Heimdall before the application starts running:

### Initialising Heimdall With Singular Redis Instance
Simple sample configuration with metrics disabled and singular redis
```go
import (
  "github.com/bytedance/heimdall"
  "github.com/bytedance/heimdall/cache"
  "github.com/go-redis/redis/v8"

)
func init() {
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
    EnableMetricsEmission: false, // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
    SkipCache:             false, // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
    CompressionLibrary :   constants.GzipCompressionType,     // Toggling compression library as GzipCompression;
    Version:               "v1.0.0", // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
  }

  if err := heimdall.Init(heimdallConfig); err != nil { 
    panic(err) // panic if heimdall fails to initialise
  }
}
```

### Initialising Heimdall With Cluster Redis Instance
Simple sample configuration with metrics disabled and cluter redis
```go
import (
  "github.com/bytedance/heimdall"
  "github.com/bytedance/heimdall/cache"
  "github.com/go-redis/redis/v8"

)
func init() {
  heimdallConfig := &heimdall.Config{ // Configuration used to initialise Heimdall
    DefaultSoftTTL: time.Second * 10, // Global softTTL if one is not provided at the time of function call. (Refer to concepts and notes for more info)
    DefaultHardTTL: time.Second * 40, // Global hardTTL if one is not provided at the time of functon call. (Refer to concepts and notes for more info)
    CacheConfig: cache.Config{ // Cache configuration, this must be provided.
      CacheProvider: constants.RedisCacheType, // Choose cache provider as Redis 
      RedisConfiguration: &cache.RedisConfig{ // Redis configuration is required if CacheProvider is set to Redis
        RedisServerType: constants.ClusterRedisType, // Server type chosen as ClusterRedisType
        ClusterConfig: &redis.ClusterOptions{ // go-redis Options (https://redis.uptrace.dev/guide/)
          Addrs:     []string{"localhost:7000", "localhost:7001", "localhost:7002"},
          Password: "", // no password set
          DB:       0,  // use default DB
        },
      },
    },
    EnableMetricsEmission: false, // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
    SkipCache:             false, // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
    CompressionLibrary :   constants.GzipCompressionType,     //Toggling compression library as GzipCompression;
    Version:               "v1.0.0", // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
  }

  if err := heimdall.Init(heimdallConfig); err != nil { 
    panic(err) // panic if heimdall fails to initialise
  }
}
```

### Using Heimdall in a gRPC call
It is pretty easy to support existing gRPC calls with a simple wrapper.

```go
import (
  "context"
	"time"

  "github.com/bytedance/heimdall"
)

// Original
resp, err := client.GetSomeFunctionCall(context.Background(), &pb.SomeFunctionCallRequest{Hello: "World"})
if err != nil {
  return err
}

// Enhanced with Heimdall, using default SoftTTL and HardTTL defined during init
resp, err := heimdall.GRPCCall(client.GetSomeFunctionCall, context.Background(), &pb.SomeFunctionCallRequest{Hello: "World"})
if err != nil {
  return err
}

// Enhanced with Heimdall with Custom-Defined TTL
resp, err := heimdall.GRPCCallWithTTL(client.GetSomeFunctionCall, context.Background(), &pb.SomeFunctionCallRequest{Hello: "World"}, 30 * time.Second, 60 * time.Second)
if err != nil {
  return err
}
```

## Advanced Usage
If one needs to use a custom defined client for Redis, your client needs to fulfill the following interfaces (if it doesn't you must wrap it): 
```go
// IGet is an interface for all cache clients that support Get operations.
type IGet interface {
    Get(ctx context.Context, key string) (string, error)
}

// ISet is an interface for all cache clients that support Set operations.
type ISet interface {
  Set(ctx context.Context, key string, val any, ttl time.Duration) error
}
```


### Initialising Heimdall With Custom Cache Instance
```go
import (
  "github.com/bytedance/heimdall"
  "github.com/bytedance/heimdall/cache"
  "github.com/go-redis/redis/v8"

)

type CustomClient struct {
	// Your implementation of CustomClient
}

func (c *CustomClient) Get(ctx context.Context, key string) ([]byte, error) {
	// Get Implementation
}

func (c *CustomClient) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	// Set Implementation
}


func init() {

  CustomClient := &CustomClient{}

  heimdallConfig := &heimdall.Config{ // Configuration used to initialise Heimdall
    DefaultSoftTTL: time.Second * 10, // Global softTTL if one is not provided at the time of function call. (Refer to concepts and notes for more info)
    DefaultHardTTL: time.Second * 40, // Global hardTTL if one is not provided at the time of functon call. (Refer to concepts and notes for more info)
    CacheConfig: cache.Config{ // Cache configuration, this must be provided.
      CacheProvider: constants.CustomCacheType, // Choose cache provider as CustomCacheType
      CustomConfiguration: &cache.CustomConfig{ // Custom configuration is required if CacheProvider is set to CustomConfiguration
        Client: CustomClient, // Custom client that implements Get and Set
      },
    },
    EnableMetricsEmission: false, // Whether or not to emit metrics. If this is true, metricsconfig must be provided.
    SkipCache:             false, // Global toggle on whether or not to skip heimdall. Useful for toggling between environments.
    CompressionLibrary :   constants.GzipCompressionType,    //Toggling compression library as GzipCompression;
    Version:               "v1.0.0", // Version number set. If there are any upgrades, this prevents breaking changes as old keys will not be re-used.
  }

  if err := heimdall.Init(heimdallConfig); err != nil { 
    panic(err) // panic if heimdall fails to initialise
  }
}
```


Heimdall can allow you to keep track of the cache hit rate, cache miss rate and cache soft hit rate. To keep track of metrics, your metrics client needs to fulfill the following interface (if it doesn't you must wrap it):
```go
// IIncreaseMetric is an interface for a basic increase metrics client.
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


```
Do note that the metric name emitted in this case will be the RPC call function name/method name.


## Examples
An example of a client and server is provided in the "examples" directory. This example is built on top of the "Hello World" project implemented in gRPC, as documented in the [grpc-go repository](https://github.com/grpc/grpc-go). The provided example includes three options for initializing Heimdall.  

### Run the example
From the root directory:
1. Ensure you have a Redis instance running locally on port 6379.
2. Compile and execute server code:
```
$ go run example/greeter_server/main.go
2023/03/10 14:06:03 server listening at [::]:50051
```

3. From another terminal, compile and execute client code:
```
$ go run example/greeter_client/main.go
heimdall initialised successfully
2023/03/10 14:06:30 Message Received: Hello world

```
The client code contains an infinite loop that sends an RPC call every 2 seconds. In the terminal running the server code, you will initially see:

```
2023/03/10 14:06:30 Received: world
```
Even though the client continues to send RPC calls, the server does not receive any, indicating that the data has been cached and the client is retrieving it from the cache instead of making a server call. After several calls, the server finally receives an RPC call, which is due to the SoftTTL being set to 10 seconds. Once the SoftTTL has expired, the RPC call to the server is invoked which updates the data in the cache asynchronously.

### Cluster configuration
To test the example with Redis Cluster configuration. Please refer to the instructions below for details.

1. To configure Redis Cluster, there is a minimum requirement of at least three nodes. In the example directory, there are three directories, namely `./7000`, `./7001`, and `./7002`. To initialise the Redis Cluster, run the following command in three different terminals from the root directory:

```
$ cd ./example/7000
$ redis-server ./redis.conf
```
```
$ cd ./example/7001
$ redis-server ./redis.conf
```

```
$ cd ./example/7002
$ redis-server ./redis.conf
```
2. In the last terminal, run the following command: 

```
$ redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002
```
- This command will create a Redis Cluster with the three nodes you started in the previous step.

3. Once all the nodes have joined the cluster, you can run the following command to visualize the effects of caching:

```
$ redis-cli -p 7000
127.0.0.1:7000> keys *
```
- This command will show all the keys in the Redis Cluster running on port 7000. You should see the cached keys from the example.


## How to test the software
Simply spin up a local Redis instance with default settings and run
``` go test -v ./...```

## Usage warning
Users of Heimdall should be mindful to only cache data that does not change very frequently.

If you have questions, concerns, bug reports, etc, please file an issue in this repository's Issue Tracker.

## Getting involved
We welcome contributions from the community to help improve Heimdall even further. Whether you have suggestions, bug reports, feature requests, or pull requests, we encourage you to submit them via our Github repository. Our team will be more than happy to review them and work with you to merge them into the codebase.

General instructions on _how_ to contribute should be stated with a link to [CONTRIBUTING](CONTRIBUTING.md).

https://redis.io/
----

## Open source licensing info
1. [TERMS](TERMS.md)
2. [LICENSE](LICENSE)


----

## Credits and references

1. [go-redis](https://github.com/go-redis/redis)
