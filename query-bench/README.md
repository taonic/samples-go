This is a basic example to benchmark workflow queries against Temporal Cloud.

Steps to run this sample:

1. Run the following command to start a worker
```
go run query-bench/worker/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>" \
    -task-queue "aws-query-1"
```

2. Run the following command to start workflow followed by query requests
```
go run query-bench/starter/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>"
    -iterations 200 \
    -payload-size 32768 \
    -task-queue "aws-query-1" \
    -query-interval-ms 250 \
    -num-workflows 10
```

3. Review the test output, e.g.
```
✅ All tests completed
Test spec: {Iterations:1000 PayloadSize:32768 TaskQueue:aws-query-1 NumberOfWorkflows:10 QueryIntervalMs:250}
Test executed in between: 2023-07-04T11:02:48Z - 2023-07-04T11:07:17Z
Time elapsed: 4m29.142188624s
Total queries: 10000
RPS: 37.16
Query latency:
P(50): 16 ms
P(90): 22 ms
P(99): 39 ms
```
