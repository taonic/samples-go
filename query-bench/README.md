This is a basic example to benchmark workflow queries against Temporal Cloud.

Steps to run this sample:

1. Run the following command to start a worker
```
go run query-bench/worker/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>"
```

2. Run the following command to start workflow followed by query requests
```
go run query-bench/starter/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>"
    -payload-size 8192 \
    -iterations 2000 \
    -concurrency 1
```

3. Review the test output, e.g.
```
✅ All tests completed
Test spec: {Concurrency:1 Iterations:2000 PayloadSize:15360 TaskQueue:query4}
Test executed in between: 2023-07-03T12:39:06Z - 2023-07-03T12:52:09Z
Time elapsed: 13m2.2542525s
Total queries: 2000
RPS: 2.56
Query latency:
P(50): 328 ms
P(90): 669 ms
P(99): 1096 ms
```
