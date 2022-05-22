# Go "functional options" are slow

A small microbenchmark showing that Go's functional arguments are slow. See my [blog post for details](https://www.evanjones.ca/go-functional-options-slow.html).


## Running

```
go test . -bench=.
```

### Checking inlining

```
go test -gcflags='-m' .
```


## Examples of APIs using functional arguments

* gRPC DialContext https://pkg.go.dev/google.golang.org/grpc@v1.45.0#DialContext and NewServer https://pkg.go.dev/google.golang.org/grpc#NewServer
* dd-trace-go StartSpanFromContext https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer#StartSpanFromContext
* OpenTelemetry Tracer.Start: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Tracer.Start
* OpenCensus StartSpan https://pkg.go.dev/go.opencensus.io/trace#StartSpan
* Zap logging New: https://pkg.go.dev/go.uber.org/zap#NewProduction
* AWS SDK V2 dynamodb GetItem: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Client.GetItem ; This API is a weird combination of both explicit config structs, and functional options. This API is also generated, and full of unfortunate details that come from generated APIs, such as "required" arguments being passed as pointers to strings (e.g. dynamodb.GetItemInput.TableName). However, the ... arguments don't leak according to the compiler!


## Example of APS using configuration structs

* http.Server and http.Client https://devdocs.io/go/net/http/index#Server https://devdocs.io/go/net/http/index#Client
* tls.Config
* gocloud, for example Bucket.NewWriter: https://pkg.go.dev/gocloud.dev@v0.25.0/blob#Bucket.NewWriter
