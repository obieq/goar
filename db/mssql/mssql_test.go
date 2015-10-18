package mssql_test

import (
	. "github.com/obieq/goar"
	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsSql", func() {
	var (
		ModelS, MK, Sprite, Panamera, Evoque, Bugatti, Out MsSqlAutomobile
		ar                                                 *MsSqlAutomobile
	)

	BeforeEach(func() {
		ar = MsSqlAutomobile{}.ToActiveRecord()
		ar.Truncate()
	})

	Context("DB Interactions", func() {
		BeforeEach(func() {
			//ModelS = MsSqlAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2009, Model: "model s"}}}.ToActiveRecord()
			ModelS = MsSqlAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "tesla", Year: 2009, Model: "model s"}}}
			ToAR(&ModelS)
			//ModelS.SetKey("id1")
			Ω(ModelS.Valid()).Should(BeTrue())

			MK = MsSqlAutomobile{SafetyRating: 3, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "3000"}}}
			ToAR(&MK)
			//MK.SetKey("id2")
			Ω(MK.Valid()).Should(BeTrue())

			Sprite = MsSqlAutomobile{SafetyRating: 2, Automobile: Automobile{Vehicle: Vehicle{Make: "austin healey", Year: 1960, Model: "sprite"}}}
			ToAR(&Sprite)
			//Sprite.SetKey("id3")
			Ω(Sprite.Valid()).Should(BeTrue())

			Panamera = MsSqlAutomobile{SafetyRating: 5, Automobile: Automobile{Vehicle: Vehicle{Make: "porsche", Year: 2010, Model: "panamera"}}}
			ToAR(&Panamera)
			//Panamera.SetKey("id4")
			Ω(Panamera.Valid()).Should(BeTrue())

			Evoque = MsSqlAutomobile{SafetyRating: 1, Automobile: Automobile{Vehicle: Vehicle{Make: "land rover", Year: 2013, Model: "evoque"}}}
			ToAR(&Evoque)
			//Evoque.SetKey("id5")
			Ω(Evoque.Valid()).Should(BeTrue())

			Bugatti = MsSqlAutomobile{SafetyRating: 4, Automobile: Automobile{Vehicle: Vehicle{Make: "bugatti", Year: 2013, Model: "veyron"}}}
			ToAR(&Bugatti)
			//Bugatti.SetKey("id6")
			Ω(Bugatti.Valid()).Should(BeTrue())
		})

		Context("DB Operations", func() {
			It("should truncate a table", func() {
				var autos []MsSqlAutomobile
				var autos2 []MsSqlAutomobile

				Ω(ModelS.Save()).Should(BeTrue())
				Ω(ModelS.ID).Should(BeNumerically(">", 0))

				Ω(Sprite.Save()).Should(BeTrue())
				Ω(Sprite.ID).Should(BeNumerically(">", ModelS.ID))

				err := MsSqlAutomobile{}.ToActiveRecord().All(&autos, nil)
				Ω(len(autos)).Should(Equal(2))

				_, err = MsSqlAutomobile{}.ToActiveRecord().Truncate()
				Ω(err).NotTo(HaveOccurred())

				err = MsSqlAutomobile{}.ToActiveRecord().All(&autos2, nil)
				Ω(len(autos2)).Should(Equal(0))
			})
		})

		Context("Persistance", func() {
			It("should create a model and find it by id", func() {
				Ω(ModelS.Save()).Should(BeTrue())
				Ω(ModelS.ID).Should(BeNumerically(">", 0))

				Ω(Sprite.Save()).Should(BeTrue())
				Ω(Sprite.ID).Should(BeNumerically(">", ModelS.ID))

				// verify
				model := Out
				err := MsSqlAutomobile{}.ToActiveRecord().Find(ModelS.ID, &model)
				Ω(err).NotTo(HaveOccurred())
				Ω(model.ID).Should(Equal(ModelS.ID))
			})

			//It("should not create a model using an existing id", func() {
			//Sprite.Delete()
			//Ω(Sprite.Save()).Should(BeTrue())

			//// reset CreatedAt
			//Sprite.CreatedAt = nil
			//success, err := Sprite.Save() // id is still the same, so save should fail
			//Ω(err).To(HaveOccurred())
			//Ω(success).Should(BeFalse())
			//})

			//It("should update an existing model", func() {
			//Sprite.Delete()
			//Ω(Sprite.Save()).Should(BeTrue())
			//year := Sprite.Year
			//modelName := Sprite.Model

			//// create
			//result := Out
			//err := ar.Find(Sprite.ID, &result)
			//Ω(err).NotTo(HaveOccurred())
			//Ω(result.ID).Should(Equal(Sprite.ID))
			//Ω(result.CreatedAt).ShouldNot(BeNil())
			//Ω(result.UpdatedAt).Should(BeNil())

			//// update
			//dbModel := result.ToActiveRecord()
			//dbModel.Year += 1
			//dbModel.Model += " updated"
			//Ω(dbModel.Save()).Should(BeTrue())

			//// verify updates
			//result = Out
			//err = ar.Find(Sprite.ID, &result)
			//Expect(err).NotTo(HaveOccurred())
			//Ω(result.Year).Should(Equal(year + 1))
			//Ω(result.Model).Should(Equal(modelName + " updated"))
			//Ω(result.CreatedAt).ShouldNot(BeNil())
			//Ω(result.UpdatedAt).ShouldNot(BeNil())
			//})

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
		})

		Context("Stored Procedures", func() {
			It("should execute a non-parameterized stored procedure that returns a results set (array)", func() {
				var autos []MsSqlAutomobile

				Ω(ModelS.Save()).Should(BeTrue())
				Ω(ModelS.ID).Should(BeNumerically(">", 0))

				Ω(Sprite.Save()).Should(BeTrue())
				Ω(Sprite.ID).Should(BeNumerically(">", ModelS.ID))

				err := MsSqlAutomobile{}.ToActiveRecord().SpExecResultSet(AUTO_LIST_SP_NAME, nil, &autos)
				Ω(err).NotTo(HaveOccurred())

				//verify
				// NOTE: the stored proc returns resultset sorted by id ASC
				Ω(len(autos)).Should(Equal(2))
				Ω(autos[0].ID).Should(Equal(ModelS.ID))
				Ω(autos[1].ID).Should(Equal(Sprite.ID))
			})

			It("should execute a parameterized stored procedure that returns a results set (array)", func() {
				var autos []MsSqlAutomobile

				Ω(ModelS.Save()).Should(BeTrue())
				Ω(ModelS.ID).Should(BeNumerically(">", 0))

				Ω(Sprite.Save()).Should(BeTrue())
				Ω(Sprite.ID).Should(BeNumerically(">", ModelS.ID))

				Ω(MK.Save()).Should(BeTrue())
				Ω(MK.ID).Should(BeNumerically(">", Sprite.ID))

				Ω(Panamera.Save()).Should(BeTrue())
				Ω(Panamera.ID).Should(BeNumerically(">", MK.ID))

				Ω(Evoque.Save()).Should(BeTrue())
				Ω(Evoque.ID).Should(BeNumerically(">", Panamera.ID))

				params := map[string]interface{}{"Id": 1, "Model": Sprite.Model}

				err := MsSqlAutomobile{}.ToActiveRecord().SpExecResultSet(AUTO_LIST_WITH_PARAMS_SP_NAME, params, &autos)
				Ω(err).NotTo(HaveOccurred())

				//verify
				// NOTE: the stored proc returns resultset sorted by id ASC
				Ω(len(autos)).Should(Equal(3))
				Ω(autos[0].ID).Should(Equal(MK.ID))
				Ω(autos[1].ID).Should(Equal(Panamera.ID))
				Ω(autos[2].ID).Should(Equal(Evoque.ID))
			})
		}) // end Context("Stored Procedures")

		Context("Querying", func() {
			BeforeEach(func() {
				// truncate
				_, err := MsSqlAutomobile{}.ToActiveRecord().Truncate()
				Ω(err).NotTo(HaveOccurred())

				// create test data
				Ω(Panamera.Save()).Should(BeTrue())
				Ω(Evoque.Save()).Should(BeTrue())
				Ω(Bugatti.Save()).Should(BeTrue())
			})

			Context("All", func() {
				It("should return all models", func() {
					var results []MsSqlAutomobile
					var dbPanamera, dbEvoque, dbBugatti MsSqlAutomobile

					err := ar.All(&results, nil)
					Ω(err).NotTo(HaveOccurred())
					Ω(len(results)).Should(Equal(3))

					for _, model := range results {
						if model.ID == Panamera.ID {
							dbPanamera = model
						} else if model.ID == Evoque.ID {
							dbEvoque = model
						} else if model.ID == Bugatti.ID {
							dbBugatti = model
						}
					}

					Ω(dbPanamera).ShouldNot(BeNil())
					Ω(dbEvoque).ShouldNot(BeNil())
					Ω(dbBugatti).ShouldNot(BeNil())

					// verify property mappings for each automobile
					dbPanamera.AssertDbPropertyMappings(Panamera, false)
					dbEvoque.AssertDbPropertyMappings(Evoque, false)
					dbBugatti.AssertDbPropertyMappings(Bugatti, false)
				})

				It("should limit the number of records returned", func() {
					var results []MsSqlAutomobile
					limit := 2

					// no limit
					err := ar.All(&results, nil)
					Ω(err).NotTo(HaveOccurred())
					Ω(len(results)).Should(Equal(3))

					// reset results
					results = []MsSqlAutomobile{}

					// limit
					err = ar.All(&results, map[string]interface{}{"limit": limit})
					Ω(err).NotTo(HaveOccurred())
					Ω(len(results)).Should(Equal(limit))
				})

				It("should return an error if limit is > 1000", func() {
					var results []MsSqlAutomobile
					limit := 1001

					err := ar.All(&results, map[string]interface{}{"limit": limit})
					Ω(err).To(HaveOccurred())
				})

				It("should return an error if limit is < 1", func() {
					var results []MsSqlAutomobile
					limit := 0

					err := ar.All(&results, map[string]interface{}{"limit": limit})
					Ω(err).To(HaveOccurred())
				})
			}) // Context: All

			//Context("Relational Operators", func() {
			//Context("Equal", func() {
			//It("should query with two EQ operators", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
			//err := ar.Where(QueryCondition{Key: "model", RelationalOperator: EQ, Value: "panamera"}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(1))

			//auto := results[0]
			//Ω(auto).ShouldNot(BeNil())
			//Ω(auto.Model).Should(Equal("panamera"))
			//})
			//})
			//})

			//Context("Logical Operators", func() {
			//Context("And", func() {
			//It("should query with two AND operators", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
			//err := ar.Where(QueryCondition{LogicalOperator: AND, Key: "model", RelationalOperator: EQ, Value: "panamera"}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(1))

			//auto := results[0]
			//Ω(auto).ShouldNot(BeNil())
			//Ω(auto.Model).Should(Equal("panamera"))
			//})
			//})

			//Context("Or", func() {
			//It("should query with two OR operators", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: EQ, Value: 2010})
			//ar.Where(QueryCondition{LogicalOperator: OR, Key: "model", RelationalOperator: EQ, Value: "veyron"})
			//err := ar.Where(QueryCondition{LogicalOperator: OR, Key: "model", RelationalOperator: EQ, Value: "gobbledygook"}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(2))
			//})
			//})
			//})

			//Context("Query Transformations", func() {
			//Context("Order Bys", func() {
			//It("should order one field ASC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("evoque"))
			//Ω(results[1].Model).Should(Equal("panamera"))
			//Ω(results[2].Model).Should(Equal("veyron"))
			//})

			//It("should order one field DESC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("veyron"))
			//Ω(results[1].Model).Should(Equal("panamera"))
			//Ω(results[2].Model).Should(Equal("evoque"))
			//})

			//It("should order the first field ASC and a second field ASC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//ar.Order(OrderBy{Key: "year", SortOrder: ASC})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("panamera"))
			//Ω(results[1].Model).Should(Equal("evoque"))
			//Ω(results[2].Model).Should(Equal("veyron"))
			//})

			//It("should order the first field ASC and a second field DESC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//ar.Order(OrderBy{Key: "year", SortOrder: ASC})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("panamera"))
			//Ω(results[1].Model).Should(Equal("veyron"))
			//Ω(results[2].Model).Should(Equal("evoque"))
			//})

			//It("should order first field DESC and a second field ASC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//ar.Order(OrderBy{Key: "year", SortOrder: DESC})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: ASC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("evoque"))
			//Ω(results[1].Model).Should(Equal("veyron"))
			//Ω(results[2].Model).Should(Equal("panamera"))
			//})

			//It("should order the first field DESC and a second field DESC", func() {
			//ar.Where(QueryCondition{Key: "year", RelationalOperator: GTE, Value: 2010})
			//ar.Order(OrderBy{Key: "year", SortOrder: DESC})
			//err := ar.Order(OrderBy{Key: "model", SortOrder: DESC}).Run(&results)

			//Ω(err).NotTo(HaveOccurred())
			//Ω(results).ShouldNot(BeNil())
			//Ω(len(results)).Should(Equal(3))

			//Ω(results[0].Model).Should(Equal("veyron"))
			//Ω(results[1].Model).Should(Equal("evoque"))
			//Ω(results[2].Model).Should(Equal("panamera"))
			//})
			//})
			//})
		})
	})
})
