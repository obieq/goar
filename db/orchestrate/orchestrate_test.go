package orchestrate_test

import (
	"fmt"
	"time"

	. "github.com/obieq/goar"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Orchestrate", func() {
	var (
		ModelS, MK, Sprite, Panamera, Evoque, Bugatti, Out IntegrationTestAutomobile
		ar                                                 *IntegrationTestAutomobile
	)

	BeforeEach(func() {
		ar = IntegrationTestAutomobile{}.ToActiveRecord()
	})

	Context("DB Interactions", func() {
		BeforeEach(func() {
			ModelS = IntegrationTestAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2009, Model: "model s"}}}
			ToAR(&ModelS)
			ModelS.SetKey("id1")
			Ω(ModelS.Valid()).Should(BeTrue())

			MK = IntegrationTestAutomobile{SafetyRating: 3, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "3000"}}}
			ToAR(&MK)
			MK.SetKey("id2")
			Ω(MK.Valid()).Should(BeTrue())

			Sprite = IntegrationTestAutomobile{SafetyRating: 2, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "sprite"}}}
			ToAR(&Sprite)
			Sprite.SetKey("id3")
			Ω(Sprite.Valid()).Should(BeTrue())

			Panamera = IntegrationTestAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "porsche", Year: 2010, Model: "panamera"}}}
			ToAR(&Panamera)
			Panamera.SetKey("id4")
			Ω(Panamera.Valid()).Should(BeTrue())

			Evoque = IntegrationTestAutomobile{SafetyRating: 1, Automobile: Automobile{Vehicle: Vehicle{Make: "land rover", Year: 2013, Model: "evoque"}}}
			ToAR(&Evoque)
			Evoque.SetKey("id5")
			Ω(Evoque.Valid()).Should(BeTrue())

			Bugatti = IntegrationTestAutomobile{SafetyRating: 4, Automobile: Automobile{Vehicle: Vehicle{Make: "bugatti", Year: 2013, Model: "veyron"}}}
			ToAR(&Bugatti)
			Bugatti.SetKey("id6")
			Ω(Bugatti.Valid()).Should(BeTrue())
		})

		Context("Persistance", func() {
			It("should create a model and find it by id", func() {
				children := map[string]interface{}{"makes": []string{"honda", "lexus", "nissan"}, "height": 52.62}
				attributes := map[string]interface{}{"age": 10, "name": "ichabod", "children": children}
				ModelS.Junk = attributes
				Ω(ModelS.Save()).Should(BeTrue())

				model := Out
				err := IntegrationTestAutomobile{}.ToActiveRecord().Find(ModelS.ID, &model)
				Ω(err).NotTo(HaveOccurred())
				Ω(model.ID).Should(Equal(ModelS.ID))
				Ω(model.Year).Should(Equal(ModelS.Year))
				Ω(model.Make).Should(Equal(ModelS.Make))
				Ω(model.Model).Should(Equal(ModelS.Model))
				Ω(model.SafetyRating).Should(Equal(ModelS.SafetyRating))
				Ω(model.Junk).ShouldNot(BeNil())
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
				Expect(err).NotTo(HaveOccurred())
				Ω(result.Year).Should(Equal(year + 1))
				Ω(result.Model).Should(Equal(modelName + " updated"))
				Ω(result.CreatedAt).ShouldNot(BeNil())
				Ω(result.UpdatedAt).ShouldNot(BeNil())
			})

			It("should delete an existing model", func() {
				// create and verify
				Ω(MK.Save()).Should(BeTrue())
				model := Out
				err := ar.Find(MK.ID, &model)
				Ω(err).NotTo(HaveOccurred())

				// delete
				err = MK.Delete()
				Ω(err).NotTo(HaveOccurred())

				// verify delete
				model = Out
				err = ar.Find(MK.ID, &model)
				Expect(err).To(HaveOccurred())
				Ω(err.Error()).Should(Equal("record not found"))
			})

			It("should return all models", func() {
				// NOTE: there's a timing issue with deleting the collection
				// delete the collection
				//IntegrationTestAutomobile{}.ToActiveRecord().Truncate()
				//time.Sleep(5000 * time.Millisecond)

				Ω(Panamera.Save()).Should(BeTrue())
				Ω(Evoque.Save()).Should(BeTrue())

				var results []IntegrationTestAutomobile
				var dbPanamera IntegrationTestAutomobile
				var dbEvoque IntegrationTestAutomobile

				err := ar.All(&results, nil)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(BeNumerically(">=", 2))

				for _, model := range results {
					if model.ID == Panamera.ID {
						dbPanamera = model
					} else if model.ID == Evoque.ID {
						dbEvoque = model
					}
				}

				Ω(dbPanamera).ShouldNot(BeNil())
				Ω(dbEvoque).ShouldNot(BeNil())

				// verify property mappings for each automobile
				dbPanamera.AssertDbPropertyMappings(Panamera, false)
				dbEvoque.AssertDbPropertyMappings(Evoque, false)
			})
		})

		Context("Querying", func() {
			var results []IntegrationTestAutomobile
			var searchDataLoaded bool = false

			BeforeEach(func() {
				results = []IntegrationTestAutomobile{}
				if !searchDataLoaded {
					fmt.Println("Loading Orchestrate Search Data")
					// first, delete all autos that may have been generated during previous tests
					Panamera.Delete()
					Evoque.Delete()
					Bugatti.Delete()

					// next, create test data
					Ω(Panamera.Save()).Should(BeTrue())
					Ω(Evoque.Save()).Should(BeTrue())
					Ω(Bugatti.Save()).Should(BeTrue())

					// wait for the new test data to be indexed
					time.Sleep(1000 * time.Millisecond)

					searchDataLoaded = true
				}
			})

			Context("Relational Operators", func() {
				Context("Equal", func() {
					It("should query with two EQ operators", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
						err := ar.Where(QueryCondition{Key: "model", RelationalOperator: EQ, Value: "panamera"}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(1))

						auto := results[0]
						Ω(auto).ShouldNot(BeNil())
						Ω(auto.Model).Should(Equal("panamera"))
					})
				})
			})

			Context("Logical Operators", func() {
				Context("And", func() {
					It("should query with two AND operators", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
						err := ar.Where(QueryCondition{LogicalOperator: AND, Key: "model", RelationalOperator: EQ, Value: "panamera"}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(1))

						auto := results[0]
						Ω(auto).ShouldNot(BeNil())
						Ω(auto.Model).Should(Equal("panamera"))
					})
				})

				Context("Or", func() {
					It("should query with two OR operators", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
						ar.Where(QueryCondition{LogicalOperator: OR, Key: "model", RelationalOperator: EQ, Value: "veyron"})
						err := ar.Where(QueryCondition{LogicalOperator: OR, Key: "model", RelationalOperator: EQ, Value: "gobbledygook"}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(2))
					})
				})
			})

			Context("Query Transformations", func() {
				Context("Order Bys", func() {
					It("should order one field ASC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("evoque"))
						Ω(results[1].Model).Should(Equal("panamera"))
						Ω(results[2].Model).Should(Equal("veyron"))
					})

					It("should order one field DESC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("veyron"))
						Ω(results[1].Model).Should(Equal("panamera"))
						Ω(results[2].Model).Should(Equal("evoque"))
					})

					It("should order the first field ASC and a second field ASC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						ar.Order(OrderBy{Key: "year", SortOrder: ASC})
						err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("panamera"))
						Ω(results[1].Model).Should(Equal("evoque"))
						Ω(results[2].Model).Should(Equal("veyron"))
					})

					It("should order the first field ASC and a second field DESC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						ar.Order(OrderBy{Key: "year", SortOrder: ASC})
						err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("panamera"))
						Ω(results[1].Model).Should(Equal("veyron"))
						Ω(results[2].Model).Should(Equal("evoque"))
					})

					It("should order first field DESC and a second field ASC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						ar.Order(OrderBy{Key: "year", SortOrder: DESC})
						err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("evoque"))
						Ω(results[1].Model).Should(Equal("veyron"))
						Ω(results[2].Model).Should(Equal("panamera"))
					})

					It("should order the first field DESC and a second field DESC", func() {
						ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
						ar.Order(OrderBy{Key: "year", SortOrder: DESC})
						err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

						Ω(err).NotTo(HaveOccurred())
						Ω(results).ShouldNot(BeNil())
						Ω(len(results)).Should(Equal(3))

						Ω(results[0].Model).Should(Equal("veyron"))
						Ω(results[1].Model).Should(Equal("evoque"))
						Ω(results[2].Model).Should(Equal("panamera"))
					})
				})
			})
		})
	})
})
