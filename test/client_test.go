package mapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sagittaros/neo4j-go-mapper/mapper"
)

type Params map[string]interface{}

type Category struct {
	ID   string
	Name string
}

type TagsWith struct {
	Agent string
}

type Document struct {
	ID string
}

type classifyRelation struct {
	Category
	TagsWith
	Document
}

var _ = Describe("Client", func() {
	BeforeEach(func() {
		mustExec(client.Exec("MATCH (n) DETACH DELETE n", nil))
	})

	Describe("Ping", func() {
		Specify("Return error for wrong connection credentials", func() {
			c, err := mapper.NewClient("bolt://localhost:7687", "neo4j", "WRONGPASS")
			Expect(err).NotTo(HaveOccurred())
			mustNotExec(c.Ping())
		})
	})

	Specify("Create a unique constraint", func() {
		mustExec(client.Exec("CREATE CONSTRAINT ON (l:Category) ASSERT l.ID IS UNIQUE", nil))
		mustExec(client.Exec("CREATE (n:Category {ID: $id})", Params{"id": "collision"}))
		mustNotExec(client.Exec("CREATE (n:Category {ID: $id})", Params{"id": "collision"}))
	})

	Specify("Create a node", func() {
		mustExec(client.Exec("CREATE (n:Category {ID: $id, Name: $name})", Params{"id": "AnID", "name": "VanGogh"}))
	})

	Describe("Create a relationship", func() {
		Specify("with CREATE", func() {
			cypher := `
				CREATE (n:Document {ID: $documentID})-[:TAGS_WITH {Agent: $agent}]->(m:Category {ID: $categoryID})`
			mustExec(client.Exec(cypher, Params{"documentID": "documentID", "agent": "user", "categoryID": "categoryID"}))
		})

		Specify("with MERGE", func() {
			cypher := `
				MERGE (n:Document {ID: $documentID})
				MERGE (m:Category {ID: $categoryID})
				MERGE (n)-[:TAGS_WITH {Agent: $agent}]->(m)`
			mustExec(client.Exec(cypher, Params{"documentID": "documentID", "agent": "user", "categoryID": "categoryID"}))
		})
	})

})
