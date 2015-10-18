package dynamodb

import (
	. "github.com/obieq/goar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dynamodb", func() {
	var (
		ModelS, MK, Sprite, Panamera, Evoque, Bugatti DynamodbAutomobile
		ar                                            *DynamodbAutomobile
	)

	BeforeEach(func() {
		ar = DynamodbAutomobile{}.ToActiveRecord()
	})

	Context("DB Interactions", func() {
		BeforeEach(func() {
			ModelS = DynamodbAutomobile{SafetyRating: 5, Make: "tesla", Year: 2009, Model: "model s"}
			ToAR(&ModelS)
			ModelS.SetKey("id1")
			Ω(ModelS.Valid()).Should(BeTrue())

			MK = DynamodbAutomobile{SafetyRating: 3, Make: "austin healey", Year: 1960, Model: "3000"}
			ToAR(&MK)
			MK.SetKey("id2")
			Ω(MK.Valid()).Should(BeTrue())

			Sprite = DynamodbAutomobile{SafetyRating: 2, Make: "austin healey", Year: 1960, Model: "sprite"}
			ToAR(&Sprite)
			Sprite.SetKey("id3")
			Ω(Sprite.Valid()).Should(BeTrue())

			Panamera = DynamodbAutomobile{SafetyRating: 5, Make: "porsche", Year: 2010, Model: "panamera"}
			ToAR(&Panamera)
			Panamera.SetKey("id4")
			Ω(Panamera.Valid()).Should(BeTrue())

			Evoque = DynamodbAutomobile{SafetyRating: 1, Make: "land rover", Year: 2013, Model: "evoque"}
			ToAR(&Evoque)
			Evoque.SetKey("id5")
			Ω(Evoque.Valid()).Should(BeTrue())

			Bugatti = DynamodbAutomobile{SafetyRating: 4, Make: "bugatti", Year: 2013, Model: "veyron"}
			ToAR(&Bugatti)
			Bugatti.SetKey("id6")
			Ω(Bugatti.Valid()).Should(BeTrue())
		})

		Context("Error Handling", func() {
			It("should return an error when the Truncate() method is called", func() {
				auto := DynamodbAutomobile{}.ToActiveRecord()
				_, err := auto.Truncate()
				Ω(err).ShouldNot(BeNil())
			})

			It("should return an error when the All() method is called", func() {
				auto := DynamodbAutomobile{}.ToActiveRecord()
				err := auto.All(auto, nil)
				Ω(err).ShouldNot(BeNil())
			})

			It("should return an error when the Search() method is called", func() {
				auto := DynamodbAutomobile{}.ToActiveRecord()
				err := auto.DbSearch(auto)
				Ω(err).ShouldNot(BeNil())
			})

			It("should return an error when trying to find an ID that doesn't exist", func() {
				var auto DynamodbAutomobile
				err := DynamodbAutomobile{}.ToActiveRecord().Find("does not exist", &auto)
				Expect(err).To(HaveOccurred())
			})

			// It("should return an error when trying to patch an ID that doesn't exist", func() {
			// 	auto := DynamodbAutomobile{}.ToActiveRecord()
			// 	auto.SetKey("does not exist")
			// 	success, err := auto.Patch()
			// 	Expect(err).To(HaveOccurred())
			// 	Ω(success).Should(BeFalse())
			// })
		})

		Context("Persistance", func() {
			It("should create a model and find it by id", func() {
				autos := map[string]interface{}{"makes": []string{"honda", "lexus", "nissan"}, "height": 52.62}
				attributes := map[string]interface{}{"age": 10, "name": "ichabod", "autos": autos}
				ModelS.Junk = attributes
				success, err := ModelS.Save()

				Ω(ModelS.ModelName()).Should(Equal("DynamodbAutomobiles"))
				Ω(err).Should(BeNil())
				Ω(success).Should(BeTrue())

				auto := DynamodbAutomobile{}
				err = DynamodbAutomobile{}.ToActiveRecord().Find(ModelS.ID, &auto)
				Ω(err).Should(BeNil())
				Ω(auto).ShouldNot(BeNil())
				Ω(auto.ID).Should(Equal(ModelS.ID))
				Ω(auto.SafetyRating).Should(Equal(ModelS.SafetyRating))
				Ω(auto.Year).Should(Equal(ModelS.Year))
				Ω(auto.Make).Should(Equal(ModelS.Make))
				Ω(auto.Model).Should(Equal(ModelS.Model))
				// Ω(auto.CreatedAt).ShouldNot(BeNil())
			})

			It("should update an existing model", func() {
				Sprite.Delete()
				Ω(Sprite.Save()).Should(BeTrue())
				year := Sprite.Year
				modelName := Sprite.Model

				// create
				var auto DynamodbAutomobile
				err := ar.Find(Sprite.ID, &auto)
				Expect(err).NotTo(HaveOccurred())
				Ω(auto).ShouldNot(BeNil())
				Ω(auto.ID).Should(Equal(Sprite.ID))
				// Ω(auto.CreatedAt).ShouldNot(BeNil())
				Ω(auto.UpdatedAt).Should(BeNil())

				// update
				auto.Year++
				auto.Model += " updated"

				success, err := auto.ToActiveRecord().Save()
				Ω(err).Should(BeNil())
				Ω(success).Should(BeTrue())

				// verify updates
				var updatedAuto DynamodbAutomobile
				err = ar.Find(Sprite.ID, &updatedAuto)
				Expect(err).NotTo(HaveOccurred())
				Ω(updatedAuto).ShouldNot(BeNil())
				Ω(updatedAuto.Year).Should(Equal(year + 1))
				Ω(updatedAuto.Model).Should(Equal(modelName + " updated"))
				// Ω(updatedAuto.CreatedAt).ShouldNot(BeNil())
				// Ω(updatedAuto.UpdatedAt).ShouldNot(BeNil())
			})

			It("should perform partial (patch) updates", func() {
				Sprite.Delete()

				// create
				Ω(Sprite.Save()).Should(BeTrue())

				// verify
				var auto DynamodbAutomobile
				err := ar.Find(Sprite.ID, &auto)
				Expect(err).NotTo(HaveOccurred())
				Ω(auto).ShouldNot(BeNil())
				Ω(auto.ID).Should(Equal(Sprite.ID))
				// Ω(auto.CreatedAt).ShouldNot(BeNil())
				Ω(auto.UpdatedAt).Should(BeNil())

				// // partial update
				// s2 := DynamodbAutomobile{SafetyRating: safetyRating + 1}.ToActiveRecord()
				// s2.SetKey(Sprite.ID)
				// success, err := s2.Patch()
				// Ω(err).Should(BeNil())
				// Ω(s2.Validation.Errors).Should(BeNil())
				// Ω(success).Should(BeTrue())
				//
				// // verify updates
				// result, err = ar.Find(Sprite.ID)
				// Expect(err).NotTo(HaveOccurred())
				// Ω(result).ShouldNot(BeNil())
				// dbModel = result.(*DynamodbAutomobile).ToActiveRecord()
				// Ω(dbModel.Year).Should(Equal(year))                     // should be no change
				// Ω(dbModel.Model).Should(Equal(modelName))               // should be no change
				// Ω(dbModel.SafetyRating).Should(Equal(safetyRating + 1)) // should have incremented by one
				// Ω(dbModel.CreatedAt).ShouldNot(BeNil())                 // should be no change
				// Ω(dbModel.UpdatedAt).ShouldNot(BeNil())                 // should have been set via active record framework
			})

			It("should delete an existing model", func() {
				var auto DynamodbAutomobile

				// create and verify
				Ω(MK.Save()).Should(BeTrue())
				err := ar.Find(MK.ID, &auto)
				Expect(err).NotTo(HaveOccurred())
				Ω(auto).ShouldNot(BeNil())
				Ω(MK.ID).Should(Equal(auto.ID))

				// delete
				err = MK.Delete()
				Ω(err).NotTo(HaveOccurred())

				// verify delete
				err = ar.Find(MK.ID, &auto)
				Ω(err).To(HaveOccurred())
			})
		})
	})
})
