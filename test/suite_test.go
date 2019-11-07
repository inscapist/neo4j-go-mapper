package mapper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/sagittaros/neo4j-go-mapper/mapper"
)

var (
	client *mapper.Client
	err    error
)

func TestGraphDB(t *testing.T) {
	RegisterFailHandler(Fail)

	if client, err = mapper.NewClient("bolt://localhost:7687", "neo4j", "password"); err != nil {
		t.Fatal(err)
	}
	if err = client.Ping(); err != nil {
		t.Fatal(err)
	}
	RunSpecs(t, "GraphDB Suite")
	if err = client.Close(); err != nil {
		t.Fatal("unable to perform database tear down")
	}
}

var _ = AfterSuite(func() {
	mustExec(client.Exec("MATCH (n) DETACH DELETE n", nil))
})

func mustExec(err error) {
	Expect(err).NotTo(HaveOccurred())
}

func mustNotExec(err error) {
	Expect(err).To(HaveOccurred())
}
