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

2. Run the following command to trigger 50 workflow executions.
```
for wid in {1..50}
do
  go run query-bench/starter/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>" \
    -i "${wid}"
done
```
This command will take a while to run since each execution is creating 500 activities

3. Run the following command to query the workflow based on the wid, e.g. 1
```
wid=1
go run query-bench/query/main.go \
    -target-host "<namespace>.tmprl.cloud:7233" \
    -namespace "<namespace>" \
    -client-cert "<cert path>" \
    -client-key "<key path>" \
    -i "${wid}"
```
By default, the query runs in 20 concurrent goroutines with each making 10 queries. RPS is shown at the end.
