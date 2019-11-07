# Neo4J go mapper

Data mapper library for [Neo4J go driver](https://github.com/neo4j/neo4j-go-driver)
Refer [test](./test) for usage examples.

## How it works
Package mapper makes heavy use of `reflect` to construct values out of specified types (an empty, initialized value of a type).
This makes it easy to read arbitrary values from neo4j client. Examples [here](./test/reader_test.go)
```
Usage:
- pass empty, initialized type(s) as the last argument(s) of `ReadSingleRow` and `ReadSingleRow`
- get back values and cast them back into the required types
```


## Installation
```
go get github.com/sagittaros/neo4j-go-mapper/mapper
```