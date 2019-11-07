package mapper

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type anObject struct {
	Name  interface{}
	Age   interface{}
	Happy interface{}
}

func (o anObject) Props() map[string]interface{} {
	return map[string]interface{}{"Name": o.Name, "Age": o.Age, "Happy": o.Happy}
}

type ScanDest struct {
	Name  string
	Age   int64
	Happy bool
}

type Item struct {
	Name string
}

var _ = Describe("scanMapToStruct", func() {
	It("scans props of an object into a struct", func() {
		o := map[string]interface{}{"Name": "Rupert", "Age": 29, "Happy": true}
		dest := &ScanDest{}
		err := scanMapToStruct(o, dest)
		Expect(err).NotTo(HaveOccurred())

		Expect(dest.Name).To(Equal("Rupert"))
		Expect(dest.Age).To(Equal(int64(29)))
		Expect(dest.Happy).To(Equal(true))
	})

	Context("When object props has nil value", func() {
		It("scans props of an object into a struct", func() {
			o := map[string]interface{}{"Name": nil, "Age": nil, "Happy": nil}
			dest := &ScanDest{}
			err := scanMapToStruct(o, dest)
			Expect(err).NotTo(HaveOccurred())

			Expect(dest.Name).To(Equal(""))
			Expect(dest.Age).To(Equal(int64(0)))
			Expect(dest.Happy).To(Equal(false))
		})
	})
})

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builder (internal)")
}
