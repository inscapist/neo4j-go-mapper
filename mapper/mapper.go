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

type Client struct {
	driver neo4j.Driver
}

func NewClient(uri, user, password string) (Mapper, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, err
	}
	return &Client{
		driver,
	}, nil
}

func (c *Client) Ping() error {
	return c.Exec("MATCH (n) RETURN n LIMIT 1", nil)
}

func (c *Client) Exec(cypher string, params map[string]interface{}) error {
	var (
		err     error
		session neo4j.Session
	)
	if session, err = c.driver.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}
	defer session.Close()

	if _, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		return result.Consume()
	}); err != nil {
		return err
	}
	return nil
}

// Query is a low level function to read row results. Use the readers from `reader.go` whenever possible
func (c *Client) Query(
	cypher string,
	params map[string]interface{},
	transform func(record neo4j.Record) interface{},
) ([]interface{}, error) {
	var (
		items   []interface{}
		err     error
		session neo4j.Session
		result  neo4j.Result
	)

	if session, err = c.driver.Session(neo4j.AccessModeWrite); err != nil {
		return nil, err
	}
	defer session.Close()

	if result, err = session.Run(cypher, params); err != nil {
		return nil, err
	}

	for result.Next() {
		if transform != nil {
			items = append(items, transform(result.Record()))
		} else {
			items = append(items, result.Record())
		}
	}

	if err = result.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// Query is a low level function to read a single row. Use the readers from `reader.go` whenever possible
func (c *Client) QuerySingle(
	cypher string,
	params map[string]interface{},
	transform func(record neo4j.Record) interface{},
) (interface{}, error) {
	var (
		item    interface{}
		err     error
		session neo4j.Session
		result  neo4j.Result
	)

	if session, err = c.driver.Session(neo4j.AccessModeWrite); err != nil {
		return nil, err
	}
	defer session.Close()

	if result, err = session.Run(cypher, params); err != nil {
		return nil, err
	}

	result.Next()
	if transform != nil {
		item = transform(result.Record())
	} else {
		item = result.Record()
	}

	if err = result.Err(); err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Client) Bootstrap(cypherStmts []string) error {
	for _, stmt := range cypherStmts {
		if err := g.Exec(stmt, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Close() error {
	return c.driver.Close()
}
