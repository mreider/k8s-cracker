# Cracker configuration

## WORKERS

Can be 0, 1, 2, 4, 8, 16, 32 or 64.

* 0 - use "number of CPU cores" cracker goroutines
* 1 - use 1 cracker goroutine
* 2 - use 2 cracker goroutines
* 4 - use 4 cracker goroutines
* etc

## MEMLEAK

* TRUE - 10 MB (megabyte) per second 
* 100 - 100 MB (megabyte) per second 

## CPULEAK

* TRUE - 60 seconds (from 0% to 100%) 
* 30 - 30 seconds (from 0% to 100%) 

# docker compose

```
docker-compose up --build --scale cracker=6
```

# Regenerate *.pb.go files

```
protoc -I rpc/ rpc/*.proto --go_out=plugins=grpc:rpc
```
