# ROUTE MANAGEMENT SERVICE

## Description
- This repo provides a scalable and modular project code 
and architecture for finding the optimal route for a delivery agent.
- The algorithm to find the opitaml root currently uses naive Travelling 
Salesman Problem algorithm with time complexity `O(n!)` . n being the
number of places the delivery agent has to visit.
- Due to the extensible design of the codebase, the routing algorithm is pluggable
as it's usage is implemented using [strategy design pattern](https://en.wikipedia.org/wiki/Strategy_pattern).

## Capabilities
- Get the optimal route for a delivery agent given a list of order (restaurant, consumer).
- Implementation of reusable graph data structure,

## Prerequisites
```
- GO 1.22
```

## Assumptions
- User can't order from multiple restaurants in single order

## Get started
### Build
```shell
make build
```

### Run
- Note need to build before running.
```shell
make run
```

### Test with coverage report
```shell
make test
```