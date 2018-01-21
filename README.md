marathon
=========

# 说明
**marathon**是用**Go**编写的一个**HTTP的RPC调用框架**，提供了**HttpClient**，通过**配置**的方式，
可以让**HttpClient**带上**服务发现**，**软负载均衡**，**健康检查**，**故障摘除**，**重试**，**限流**，**监控统计上报**，**日志打印**等功能。
    
目前流行的微服务架构中，经常会有很多HTTP请求调用。看似简单的HTTP的调用，其实也有很多技术的细节需要考虑，
例如:

1）如何获取请求服务的IP和Port？
    
2）当获取到一堆IP和Port的时候，如何选取哪台机器访问？

3）如何确认被访问的机器都是健康可用的？
 
4）当被访问的机器出现故障时，如何自动从可选机器列表中摘除，当机器恢复时，又自动添加回可选机器列表中？

5）如何保证当自身流量突增的时候，不会把被请求的服务打挂？或者在大批量请求下游服务的时候，怎么平滑请求而不把下游打挂？

6）请求时，需要埋点上报请求服务的状况，方便监控、报警及接口优化。

7）每个接口的SLA都不一样，如何通过配置的方式，对每个接口问做一些差异化的访问（如重试次数、访问超时时间等）？

等等。
    
因此，需要一个统一的框架来来解决这些共性的问题。这就是我创立marathon项目的初衷。
    
项目名字**marathon**，即中文的“**马拉松**”。之所以将项目名字取名为**marathon**，是两方面的考量：
第一，本人非常喜欢跑步，marathon是距离最长的跑步赛事，也是奥运会最后的比赛项目；
第二，本项目都是我个人独自利用业余时间设计和编写，前前后后经历4个多月时间，这个项目的开发对我来说也是一项marathon式的长跑。 
   
# 特点

1. 服务发现。

    marathon框架没有真正实现服务发现的逻辑，只是提供抽象的interface方便和服务发现配合使用。
用户只需要将服务发现的逻辑实现在server.List的GetInitialListOfServers和GetUpdatedListOfServers
这两个方法中，marathon的HttpClient就具有服务发现的功能。

-----------------

2. 软负载均衡。
    
    当服务发现或者配置获取到一堆ip和port时，需要有合适的策略选取访问的机器。marathon提供软负载均衡，提供Random（随机）、RoundRobin（轮询）、
LeastConnection（最少连接数）、LeastResponseTime（最少响应时间）、Hash（哈希）、WeightedResponseTime（加权的最小响应时间）六中常用的负载均衡算法来选取机器。marathon提供软负载均衡的框架和负载均衡算法的抽象loadbalancer.Rule，
用户可以很方便的开发自己的负载均衡算法。

-----------------

3. 健康检查。
    
    为了保证被访问的下游机器列表都是可用的，需要对下游机器做周期性的健康检查。当发现下游机器不可用的时候，从可选机器列表中移除；当下游机器又可用的时候，
重新添加到可选机器列表中。marathon提供周期性检查机器是否可用的框架，并对健康检查提供抽象loadbalancer.Ping，用户只需要实现这个interface，就可以实现
自己的健康检查逻辑。

-----------------

4. 故障摘除。
    
    当访问某台机器时，某类错误连续出现多次时(例如http_status是502/503/504或者连接拒绝)，很有可能是机器出现故障，需要临时摘除，等休眠一段时间后再访问。
真正从可选列表中摘除是健康检查模块来做。marathon集成了故障临时自动摘除的逻辑。用户可以配置连续出错的阈值，自定义哪些出错的类型是被认为是需要摘除的错误。    
    
-----------------

5. 重试。

    当访问失败的时候，根据业务的需要来决定是否重试及重试次数。marathon提供一个优雅的重试机制，可以针对接口级别设置重试方案。
    
-----------------

6. 限流。
    
    面对突发流量，如果被请求的服务容量有限且没有做限流保护等措施，可能导致下游服务被打挂，从而影响整个服务。因此服务的调用方也是有义务来保护服务的提供方。
marathon内置限流模块，提供MaxConcurrency/MaxRequest(最大并发/最大请求数)、TokenBucket(令牌桶)和LeakyBucket(漏桶)三种限流算法，可以通过配置的方式
定制针对某个接口的限流策略，非常自由灵活。同时marathon提供框架，方便用户开发自己的限流算法。
    
-----------------

7. 监控统计上报。
    
    为了保证程序的稳定性，添加监控、统计上报和报警是必不可少。marathon有对访问过程中进行埋点，也提供抽象的方法metric.Collector，只要实现Collector的方法就能够实现数据上报。

-----------------
 
8. 配置。 

    marathon提供对配置的抽象config.ClientConfig。所有功能都有默认的配置，因此即使不传任何配置，marathon都能够正常运行，也允许用户自定义配置来覆盖默认值。
同时，所有这些功能都可以通过配置来选择打开或者关闭。配置的粒度可以细化到针对某个接口做个性化配置项，而达到对某个接口进行差异化的访问。
   
-----------------

9. HttpClient。

    marathon对官方的HttpClient做了简单的封装，让原生的HttpClient具有了服务发现、软负载均衡、健康检查、故障临时摘除、限流、统计上报、配置、日志打印等功能；
保留了和原生HttpClient基本上相同的API，方便使用。

-----------------

10. 可扩展性。
    
    marathon在设计时就考虑到程序的可扩展性。很多功能都是采用可配置的插件设计，用户自定义的功能只需要实现插件的接口，然后注册到marathon，就能够得到执行。

-----------------

# 使用

1. 配置

    marathon提供配置的抽象config.ClientConfig。提供了默认实现config.DefaultClientConfig，默认实现采用.properties格式的配置文件。
例如，我们的ClientName是demo，默认的请求超时时间是500ms，我们要修改请求的超时时间为300ms。我们这样配置：

``` .properties
    demo.RequestTimeout = 300ms
```
    
-----------------

2. 负载均衡。

    创建负载均衡器的示例代码：
    
``` go
    //Step 1:
    //负载均衡的算法采用随机选择的算法。
    rule := loadbalancer.NewRandomRule()
    
    //Step 2:
    //设置健康检查的方法。默认提供URLPing的方法
    //也可以实现ping.Ping接口的方法，实现自己的健康检查方法。
    pingAction := ping.NewURLPing("/health/check", "SUCCESS")
    
    //Step 3:
    //设置健康检查的执行策略策略
    //提供ParallelStrategy（并发执行健康检查，默认）和SerialStrategy（串行执行健康检查）
    //也可以实现ping.Strategy接口的方法，实现自己的健康检查执行策略。
    pingStrategy := ping.NewParallelStrategy()
    
    //Step 4:
    //读取配置
    clientConfig := config.NewDefaultClientConfig("demo", props)
    
    //Step 5:
    //创建固定servers列表的loadbalancer
    lb := loadbalancer.NewBaseLoadBalancer(clientConfig, rule, pingAction, pingStrategy)
    
    //Step 6:
    //创建serverList对象，解析clientConfig中的服务器列表
    //服务器列表采用这种格式：http://127.0.0.1:8080@cluster1,http://localhost:8080@cluster2
    serverList := server.NewConfigurationBasedServerList(clientConfig)
    
    //Step 7:
    //获取得到服务器列表
    servers := serverList.GetInitialListOfServers()
    
    //Step 8:
    //将服务器列表添加入loadbalancer中。
    lb.AddServers(servers)
```

-----------------

3. 服务发现。
    
    marathon没有直接提供服务发现的功能，不过提供接口来融入服务发现的功能。
示例代码如下：    
    
``` go
    //Step 1: 
    //ServiceDiscoveryList 定义用服务发现获取机器列表的类
    type ServiceDiscoveryList struct {}
    //GetInitialListOfServers 获取初始化的机器列表
    func (l *ServiceDiscoveryList)GetInitialListOfServers() []*server.Server {
        //TODO: Add your code
    }

    //GetUpdatedListOfServers 获取更新的机器列表
    func (l *ServerDiscoveryList)GetUpdatedListOfServers() []*server.Server {
        //TODO: Add your code
    }
    
    //Step 2:
    //读取配置
    clientConfig := config.NewDefaultClientConfig("demo", props)
    
    //Step 3:
    //负载均衡的算法采用随机选择的算法。
    rule := loadbalancer.NewRandomRule()
        
    //Step 4:
    //创建动态servers列表的loadbalancer
    //注意：动态列表的loadbalancer默认没有开启健康检查，是因为通过服务发现都能够动态获取机器列表，
    //就没有必要检查机器的健康状态，让服务发现来保证每次获取最新健康的机器列表。
    lb := loadbalancer.NewDynamicServerListLoadBalancer(clientConfig, rule, &ServiceDiscoveryList{})
```

-----------------

4. 健康检查。 
    
    marathon提供主动的健康检查。默认实现是URLPing的策略，可以实现自己的健康检查策略。
健康检查的执行策略，marathon提供SerialStrategy和ParallelStrategy。也可以实现自己的健康检查执行策略。
示例：

``` go
    //Ping Interface that defines how we "ping" a server to check if its alive
    type Ping interface {
	    //IsAlive Checks whether the given Server is "alive"
	    IsAlive(*server.Server) bool
    }
    
    //MyPing 自定义的健康检查方法。
    type MyPing struct {}
    
    //IsAlive ....
    func (p *Myping)IsAlive(*server.Server) bool {
        //TODO: Add your code ...
    }
    
    //Strategy defines the strategy, used to ping all servers, registered in BaseLoadBalancer.
    //You would typically create custom implementation of this interface, if you
    //want your servers to be pinged in parallel.
    type Strategy interface {
    	//PingServers ...
    	PingServers(ping Ping, servers []*server.Server) []bool
    }
    
    //MyStrategy 自定义的健康检查执行策略。
    type MyStrategy struct {}
    
    //PingServers 自定义健康检查执行策略。
    func (s *MyStrategy)PingServers(ping Ping, servers []*server.Server) []bool {
        //TODO: Add your code ...
    }
```
        
-----------------

5. 故障摘除。

    当连续多次连不上服务器或者连读多次HTTP_StatusCode都是502/503/504，会自动摘除该机器一段时间。连续出错的次数和机器摘除的时间可以通过配置自定义。
示例：

``` go
    clientConfig := config.NewDefaultClientConfig("demo", prop)
    //某类错误连续出现五次，则该机器会自动摘除一段时间。
    clientConfig.Set("ConnectionFailureThreshold", 5)
    //摘除的最大时间。
    clientConfig.Set("CircuitTripMaxTimeout", 60 * time.Second)
```

-----------------

6. 重试。

    设置重试的示例代码：
    
``` go
    //如果需要对request进行差异化的控制，构建request级别的Config.
    requestConfig := config.NewDefaultConfig("example", nil)
    //出错后，在第一次选取的机器上重试一次
    requestConfig.Set("MaxAutoRetries", 1)
    //如果在出错，重新选取机器后再请求一次
    requestConfig.Set("MaxAutoRetriesNextServer", 1)
```    
    
-----------------

7. 限流。

    marathon提供接口级别的限流。提供限流的抽象，并提供提供MaxConcurrency/MaxRequest(最大并发/最大请求数)、
TokenBucket(令牌桶)和LeakyBucket(漏桶)三种限流算法。使用示例如下：

``` go
    //创建request级别的配置。
    //例如对/ratelimit接口，我们的限流配置如下
    requestConfig := config.NewDefaultClientConfig("ratelimit", nil)
    //打开最大并发数限流控制
    requestConfig.Set("ConcurrencyRateLimitSwitch", true)
    //设置最大并发数50
    requestConfig.Set("MaxConnectionsPerHost", 50)
    //打开令牌桶限流控制
    requestConfig.Set("TokenBucketRateLimitSwitch", true)
    //设置令牌桶的最大容量
    requestConfig.Set("TokenBucketCapacity", 50)
    //设置令牌放置的周期
    requestConfig.Set("TokenBucketFillInterval", 100 * time.Millisecond)
    //设置令牌放置的数量
    requestConfig.Set("TokenBucketFillCount", 2)
    //打开漏桶限流控制
    requestConfig.Set("LeakyBucketRateLimitSwitch", true)
    //设置漏桶的最大容量
    requestConfig.Set("LeakyBucketCapacity", 50)
    //设置漏桶的漏桶周期
    requestConfig.Set("LeakyBucketInterval", 50 * time.Millisecond)
```

-----------------

8. 监控统计上报。
    
    marathon提供了监控上报的抽象类metric.Collector。
使用示例如下：
``` go
    //Step 1:
    //MyCollector 定义自己的Collector
    type MyCollector struct {}
    
    //RPC 实现自己的上报逻辑...
    func (c *MyCollector)RPC(ctx context.Context, req client.Request, resp client.Response, err error, t time.Duration) {
        //TODO: Add your code ...
    }
    
    //Step 2:
    //将自己定义的Collector注册进marathon。
    metric.RegisterCollectors(&MyCollector{})
```    
    
-----------------

9. HttpClient。

    HttpClient使用示例：

``` go
    //Step 1:
    //定义带loadbalancer 功能的httpClient
    httpClient := httpclient.NewHTTPLoadBalancerClient(clientConfig, lb)
    //自定义给所有的请求统一加上某些Header
    httpClient.RegisterBeforeHook(func(ctx context.Context, req *HTTPRequest){
        req.Header.Set("Marathon-Extension", "marathon")
    })
    
    //Step 2:
    //构造请求。
    body := bytes.NewBuffer([]byte(`{"name":"nienie","hobby":"marathon"}`))
    request, err := httpclient.NewHTTPRequest(http.MethodPost, "/example", body, nil)
    if err != nil {
        //TODO: Add your code
    }
    
    //Step 3:
    //如果需要对request进行差异化的控制，构建request级别的Config.
    requestConfig := config.NewDefaultConfig("example", nil)
    //出错后，在第一次选取的机器上重试一次
    requestConfig.Set("MaxAutoRetries", 1)
    //如果在出错，重新选取机器后再请求一次
    requestConfig.Set("MaxAutoRetriesNextServer", 1)
    
    //Step 4:
    //请求
    response, err := httpClient.Do(ctx, reqeust, requestConfig)
```

10. 日志打印。
    
    marathon提供默认的Logger来打印日志，默认的Logger直接将日志输出到标准输出。用户也可以使用自定义的Logger组件，具体步骤如下：
    
``` go
    //Step 1: 实现Logger中的方法
    //Logger ...
    type Logger interface {
    
        //Debugf ...
        Debugf(ctx context.Context, format string, args ...interface{})
    
        //Infof ...
        Infof(ctx context.Context, format string, args ...interface{})
    
        //Warnf ...
        Warnf(ctx context.Context, format string, args ...interface{})
    
        //Errorf ...
        Errorf(ctx context.Context, format string, args ...interface{})
    
        //SetLevel ...
        SetLevel(Level)
    }
    
    //MyLogger ...
    type MyLogger struct {
        //TODO: Add your code
    }
    
    //Debugf ...
    func (l *MyLogger)Debugf(ctx context.Context, format string, args ...interface{}) {
        //TODO: Add your code
    }
    
    //Infof ...
    func (l *MyLogger)Infof(ctx context.Context, format string, args ...interface{}) {
        //TODO: Add your code
    }

    //Warnf ...
    func (l *MyLogger)Warnf(ctx context.Context, format string, args ...interface{}) {
        //TODO: Add your code
    }    

    //Errorf ...
    func (l *MyLogger)Errorf(ctx context.Context, format string, args ...interface{}) {
        //TODO: Add your code
    }        
    
    //SetLevel ...
    func (l *MyLogger)SetLevel(Level) {
        //TODO: Add your code
    }
    
    //Step 2: 将自定义的Logger注册到marathon中。
    logger.SetLogger(&MyLogger{})
```

-----------------

# TODO 
由于本人时间有限，所以文档、测试用例和示例都不全。未来有时间再补充。
1）文档
2）测试用例。
3）示例。

# 联系
QQ: 525999199@qq.com