go build -buildmode=plugin -o BenchmarkOneSimpleAspect.so BenchmarkOneSimpleAspect/BenchmarkOneSimpleAspect.go
go build -buildmode=plugin -o Benchmark50SimpleAspect.so Benchmark50SimpleAspect/Benchmark50SimpleAspect.go
go build -buildmode=plugin -o Benchmark50SimpleAspectsAsyncModel.so Benchmark50SimpleAspectsAsyncModel/Benchmark50SimpleAspectsAsyncModel.go

