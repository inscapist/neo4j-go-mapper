package mapper

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// This is the interface implemented by Client
type Mapper interface {
	// Ensure Neo4J connection
	Ping() error

	// Closes bolt driver
	Close() error

	// Execute a cypher statement
	Exec(cypher string, params map[string]interface{}) error

	// Query all results/rows from Neo4J
	Query(
		cypher string,
		params map[string]interface{},
		transform func(record neo4j.Record) interface{},
	) ([]interface{}, error)

	// Query a single Row from Neo4J, for example when the result is a `count`
	QuerySingle(
		cypher string,
		params map[string]interface{},
		transform func(record neo4j.Record) interface{},
	) (interface{}, error)

	// The following 2 functions are Reader utilities for convenience.
	// Pass in initiated empty values in the ordering that corresponds to result elements, cast it back such as `val.(MyType)`
	ReadSingleRow(cypher string, params map[string]interface{}, blankTypes ...interface{}) ([]interface{}, error)
	ReadRows(cypher string, params map[string]interface{}, blankTypes ...interface{}) ([][]interface{}, error)

	// Use this to run `CREATE INDEX/CONSTRAINTS`
	Bootstrap(cypherStmts []string) error
}
