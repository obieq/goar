package couchbase_test

import (
	. "github.com/obieq/goar"
	. "github.com/obieq/goar/db/couchbase/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/obieq/goar/db/couchbase/Godeps/_workspace/src/github.com/onsi/gomega"
	. "github.com/obieq/goar/tests/models"
)

var _ = Describe("Couchbase", func() {
	var (
		ModelS, MK, Sprite, Panamera, Evoque, Bugatti, Out CouchbaseAutomobile
		ar                                                 *CouchbaseAutomobile
	)

	BeforeEach(func() {
		ar = CouchbaseAutomobile{}.ToActiveRecord()
		// delete instances from prior test
		ids := []string{"id1", "id2", "id3"}
		for idx := range ids {
			ca := CouchbaseAutomobile{}.ToActiveRecord()
			ca.SetKey(ids[idx])
			ca.Delete()
		}
	})

	Context("DB Interactions", func() {
		BeforeEach(func() {
			//ModelS = CouchbaseAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2009, Model: "model s"}}}.ToActiveRecord()
			ModelS = CouchbaseAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2009, Model: "model s"}}}
			ToAR(&ModelS)
			ModelS.SetKey("id1")
			Ω(ModelS.Valid()).Should(BeTrue())

			MK = CouchbaseAutomobile{SafetyRating: 3, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "3000"}}}
			ToAR(&MK)
			MK.SetKey("id2")
			Ω(MK.Valid()).Should(BeTrue())

			Sprite = CouchbaseAutomobile{SafetyRating: 2, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "sprite"}}}
			ToAR(&Sprite)
			Sprite.SetKey("id3")
			Ω(Sprite.Valid()).Should(BeTrue())

			Panamera = CouchbaseAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "porsche", Year: 2010, Model: "panamera"}}}
			ToAR(&Panamera)
			Panamera.SetKey("id4")
			Ω(Panamera.Valid()).Should(BeTrue())

			Evoque = CouchbaseAutomobile{SafetyRating: 1, Automobile: Automobile{Vehicle: Vehicle{Make: "land rover", Year: 2013, Model: "evoque"}}}
			ToAR(&Evoque)
			Evoque.SetKey("id5")
			Ω(Evoque.Valid()).Should(BeTrue())

			Bugatti = CouchbaseAutomobile{SafetyRating: 4, Automobile: Automobile{Vehicle: Vehicle{Make: "bugatti", Year: 2013, Model: "veyron"}}}
			ToAR(&Bugatti)
			Bugatti.SetKey("id6")
			Ω(Bugatti.Valid()).Should(BeTrue())
		})

		Context("Persistance", func() {
			It("should create a model and find it by id", func() {
				success, err := ModelS.Save()
				Ω(success).Should(BeTrue())

				model := Out
				err = CouchbaseAutomobile{}.ToActiveRecord().Find(ModelS.ID, &model)
				Ω(err).NotTo(HaveOccurred())
				Ω(model.ID).Should(Equal(ModelS.ID))
			})

			It("should create a model with an auto-generated id", func() {
				ModelS.ID = ""
				success, err := ModelS.Save()
				Ω(success).Should(BeTrue())

				model := Out
				err = CouchbaseAutomobile{}.ToActiveRecord().Find(ModelS.ID, &model)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(model.ID)).Should(Equal(36))
			})

			It("should not create a model using an existing id", func() {
				Sprite.Delete()
				Ω(Sprite.Save()).Should(BeTrue())

				// reset CreatedAt
				Sprite.CreatedAt = nil
				success, err := Sprite.Save() // id is still the same, so save should fail
				Ω(err).To(HaveOccurred())
				Ω(success).Should(BeFalse())
			})

			It("should update an existing model", func() {
				Sprite.Delete()
				Ω(Sprite.Save()).Should(BeTrue())
				year := Sprite.Year
				modelName := Sprite.Model

				// create
				result := Out
				err := ar.Find(Sprite.ID, &result)
				Ω(err).NotTo(HaveOccurred())
				Ω(result.ID).Should(Equal(Sprite.ID))
				Ω(result.CreatedAt).ShouldNot(BeNil())
				Ω(result.UpdatedAt).Should(BeNil())

				// update
				dbModel := result.ToActiveRecord()
				dbModel.Year++
				dbModel.Model += " updated"
				Ω(dbModel.Save()).Should(BeTrue())

				// verify updates
				result = Out
				err = ar.Find(Sprite.ID, &result)
				Ω(err).NotTo(HaveOccurred())
				Ω(result.Year).Should(Equal(year + 1))
				Ω(result.Model).Should(Equal(modelName + " updated"))
				Ω(result.CreatedAt).ShouldNot(BeNil())
				Ω(result.UpdatedAt).ShouldNot(BeNil())
			})

			It("should delete an existing model", func() {
				// create and verify
				Ω(MK.Save()).Should(BeTrue())
				result := Out
				err := ar.Find(MK.ID, &result)
				Ω(err).NotTo(HaveOccurred())

				// delete
				err = MK.Delete()
				Ω(err).NotTo(HaveOccurred())

				// verify delete
				result = Out
				err = ar.Find(MK.ID, &result)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(Equal("Key not found."))
			})
		})
	})
})
