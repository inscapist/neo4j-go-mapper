package mapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Result Readers", func() {
	BeforeEach(func() {
		mustExec(client.Exec("MATCH (n) DETACH DELETE n", nil))
	})

	JustBeforeEach(func() {
		client.Exec(`CREATE (n:Document {ID: $documentID})`, Params{"documentID": "doc"})
		client.Exec(`
				MATCH (n:Document {ID: $documentID})
				MERGE (m:Category {ID: $categoryID})
				MERGE (n)-[:TAGS_WITH {Agent: $agent}]->(m)`,
			Params{"documentID": "doc", "agent": "user", "categoryID": "catA"})
		client.Exec(`
				MATCH (n:Document {ID: $documentID})
				MERGE (m:Category {ID: $categoryID})
				MERGE (n)-[:TAGS_WITH {Agent: $agent}]->(m)`,
			Params{"documentID": "doc", "agent": "user", "categoryID": "catB"})
		client.Exec(`
				MATCH (n:Document {ID: $documentID})
				MERGE (m:Category {ID: $categoryID})
				MERGE (n)-[:TAGS_WITH {Agent: $agent}]->(m)`,
			Params{"documentID": "doc", "agent": "user", "categoryID": "catC"})
	})

	Describe("ReadSingleRow", func() {
		It("reads arbitrary types and returns a slice", func() {
			row, err := client.ReadSingleRow("MATCH (n:Document)-[p:TAGS_WITH]->(m:Category) RETURN n, m, p, 300 ORDER BY m.ID asc", nil, Document{}, Category{}, TagsWith{}, int64(0))
			relationship := classifyRelation{
				Category: row[1].(Category),
				TagsWith: row[2].(TagsWith),
				Document: row[0].(Document),
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(relationship).To(Equal(classifyRelation{
				Category: Category{ID: "catA"},
				TagsWith: TagsWith{Agent: "user"},
				Document: Document{ID: "doc"},
			}))
			// shows that builtin types work too
			Expect(row[3].(int64)).To(Equal(int64(300)))
		})

		It("returns nil if pattern does not match anything", func() {
			row, err := client.ReadSingleRow(
				"MATCH (n:NonExistent)-[p:TAGS_WITH]->(m:Category) RETURN n, m, p, 300 ORDER BY m.ID asc",
				nil, Document{}, Category{}, TagsWith{}, int64(0))
			Expect(err).NotTo(HaveOccurred())
			Expect(row).To(BeNil())
		})
	})

	Describe("ReadRows", func() {
		It("reads arbitrary types and returns a slice of slices", func() {
			rows, err := client.ReadRows(`
				MATCH (n:Document)-[p:TAGS_WITH]->(m:Category) 
				RETURN n, m, p, 300, [1,2,3], ["a", "b"] 
				ORDER BY m.ID asc`, nil,
				Document{}, Category{}, TagsWith{}, int64(0), []int64{}, []string{})
			Expect(err).NotTo(HaveOccurred())
			for i, catName := range []string{"catA", "catB", "catC"} {
				row := rows[i]
				relationship := classifyRelation{
					Category: row[1].(Category),
					TagsWith: row[2].(TagsWith),
					Document: row[0].(Document),
				}
				test := classifyRelation{
					Category: Category{ID: catName},
					TagsWith: TagsWith{Agent: "user"},
					Document: Document{ID: "doc"},
				}
				Expect(relationship).To(Equal(test))

				// shows that builtin types work too
				Expect(row[3].(int64)).To(Equal(int64(300)))
				Expect(row[4].([]int64)).To(Equal([]int64{1, 2, 3}))
				Expect(row[5].([]string)).To(Equal([]string{"a", "b"}))
			}
		})

		It("returns nil if pattern does not match anything", func() {
			rows, err := client.ReadRows(`
				MATCH (n:NonExistent)-[p:TAGS_WITH]->(m:Category) 
				RETURN n, m, p, 300, [1,2,3], ["a", "b"] 
				ORDER BY m.ID asc`, nil,
				Document{}, Category{}, TagsWith{}, int64(0), []int64{}, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(BeNil())
		})
	})

})
