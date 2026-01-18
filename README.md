# Distributed Rate Limiter

This repo contains a distributed rate limiter using redis. 

### Structure
The structure of the repo looks as follows:
```
distributed-rate-limiter/
├── cmd/
│   └── server/
│       └── main.go 
├── internal/
│   ├── limiter/
│   │   ├── bucket.go 
│   │   └── redis.go
│   └── middleware/
│       └── ratelimit.go
├── go.mod
├── go.sum
└── README.md
```

I've chosen to put the main.go in a cmd/server folder, to have the possibility to extend it to have multiple binaries eg a cli monitoring tool. The core logic is put in an internal folder as they should not be imported by other projects. Furthermore, the algorithm (limiter/) is divided from the transport (middleware/) (-> Single Responsibilty Principle).


### Token Bucket Algorithm
Instead of a simple counter, the token bucket algorithm is a rate-controlled system defined by two parameters: Capacity(b) and Refill Rate (r).
This makes the algorithm robust for burstiness as the user can instantly use b tokens after not having made a request for a while but is then throttled to r requests per second.

### bucket.go
The `bucket.go` file contains the core token bucket algorithm logic as well as a manger to handle buckets for multiple users. It features a getbucket method that uses the check-lock-check pattern for high-concurrent environments.
