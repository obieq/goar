package cloudant_test

import (
	. "github.com/obieq/goar/db/cloudant"
	. "github.com/obieq/goar/tests/models"

	. "github.com/obieq/goar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CloudantAutomobile struct {
	ArCloudant
	Automobile
	SafetyRating int
}

func (m *CloudantAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (model CloudantAutomobile) ToActiveRecord() *CloudantAutomobile {
	return ToAR(&model).(*CloudantAutomobile)
}

var _ = Describe("Cloudant", func() {
	var (
		//DbModel *CloudantAutomobile = CloudantAutomobile{}.ToActiveRecord()
		ModelS *CloudantAutomobile
		MK     *CloudantAutomobile
		Sprite *CloudantAutomobile
	)
	It("should initialize client", func() {
		client := Client()
		Ω(client).ShouldNot(BeNil())
	})

	Context("DB Interactions", func() {
		BeforeEach(func() {
			ModelS = CloudantAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2014, Model: "model s"}}}.ToActiveRecord()
			ModelS.SetId("id1")
			Ω(ModelS.Valid()).Should(BeTrue())

			MK = CloudantAutomobile{SafetyRating: 3, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "3000"}}}.ToActiveRecord()
			MK.SetId("id2")
			Ω(MK.Valid()).Should(BeTrue())

			Sprite = CloudantAutomobile{SafetyRating: 2, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "sprite"}}}.ToActiveRecord()
			Sprite.SetId("id3")
			Ω(Sprite.Valid()).Should(BeTrue())
		})

		Context("Persistance", func() {
			It("should persist a new model and find it by id", func() {
				Ω(ModelS.Save()).Should(BeTrue())

				result, _ := CloudantAutomobile{}.ToActiveRecord().Find(ModelS.Id())
				Ω(result).ShouldNot(BeNil())
				model := result.(*CloudantAutomobile)
				Ω(model.Id()).Should(Equal(ModelS.Id()))
			})
		})
	})
})
