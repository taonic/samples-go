### Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run the following command to start the worker
```
go run retry-nde/worker/main.go
```
3) Run the following command to start the example
```
go run retry-nde/starter/main.go
```
4) Follow the "Restart worker now" output to restart worker
5) There is 50% chance the workflow would hit NDE then retry
