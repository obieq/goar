package goar

import (
	"errors"
	"reflect"

	. "github.com/obieq/goar/tests/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ActiveRecordVehicle struct {
	ActiveRecord
	Timestamps
	Vehicle
}

type ActiveRecordAutomobile struct {
	ActiveRecordVehicle
}

type ActiveRecordMotorcycle struct {
	ActiveRecordVehicle
}

type CallbackErrorModel struct {
	ActiveRecordVehicle
	Name string
}

func (ar *ActiveRecordVehicle) SetKey(key string) {

}

func (ar *ActiveRecordVehicle) All(interface{}, map[string]interface{}) error {
	return nil
}

func (ar *ActiveRecordVehicle) Truncate() (numRowsDeleted int, err error) {
	return 0, nil
}

func (ar *ActiveRecordVehicle) Find(id interface{}, out interface{}) error {
	return nil
}

func (model ActiveRecordAutomobile) ToActiveRecord() *ActiveRecordAutomobile {
	return ToAR(&model).(*ActiveRecordAutomobile)
}

func (model ActiveRecordMotorcycle) ToActiveRecord() *ActiveRecordMotorcycle {
	return ToAR(&model).(*ActiveRecordMotorcycle)
}

func (model CallbackErrorModel) ToActiveRecord() *CallbackErrorModel {
	return ToAR(&model).(*CallbackErrorModel)
}

func (model *ActiveRecordAutomobile) DbSave() (err error) {
	return nil
}

func (model *ActiveRecordAutomobile) DbDelete() (err error) {
	return nil
}

func (model *ActiveRecordAutomobile) DbSearch(results interface{}) error {
	return nil
}

func (model *CallbackErrorModel) DbSave() (err error) {
	return nil
}

func (model *CallbackErrorModel) DbDelete() (err error) {
	return nil
}

func (model *CallbackErrorModel) DbSearch(results interface{}) error {
	return nil
}

func (m *ActiveRecordAutomobile) Validate() {
	m.Validation.Required("Year", m.Year)
	m.Validation.Required("Make", m.Make)
	m.Validation.Required("Model", m.Model)
}

func (m *ActiveRecordMotorcycle) Validate() {
}

func (m *CallbackErrorModel) Validate() {
}

var beforeSaveCallback = func() error {
	return errors.New("some error")
}

func (m *CallbackErrorModel) BeforeSave() error {
	return beforeSaveCallback()
}

func (m *CallbackErrorModel) AfterSave() error {
	return errors.New("some error")
}

func (m *ActiveRecordMotorcycle) ModelName() string {
	return "motocicletas"
}

func (m *ActiveRecordAutomobile) BeforeSave() error {
	m.Year = 1888
	return nil
}

func (m *ActiveRecordAutomobile) BeforeSaveError() error {
	return errors.New("some error")
}

func validAutomobileFactory() *ActiveRecordAutomobile {
	automobile := ActiveRecordAutomobile{}
	automobile.Year = 2007
	automobile.Make = "porsche"
	automobile.Model = "carrera gt"
	return automobile.ToActiveRecord()
}

var _ = Describe("ActiveRecord", func() {

	var (
		automobile *ActiveRecordAutomobile
		motorcycle *ActiveRecordMotorcycle
	)

	BeforeEach(func() {
		automobile = validAutomobileFactory()
		motorcycle = ActiveRecordMotorcycle{}.ToActiveRecord()
	})

	It("should convert a struct to an ActiveRecord", func() {
		automobile := &ActiveRecordAutomobile{}
		Ω(automobile.Self()).Should(BeNil())
		Ω(automobile.Query()).Should(BeNil())

		Ω(ToAR(automobile)).ShouldNot(BeNil())

		Ω(automobile.Self()).ShouldNot(BeNil())
		Ω(automobile.Query()).ShouldNot(BeNil())
	})

	Context("Model Name", func() {
		It("should derive model name", func() {
			Ω(automobile.ModelName()).Should(Equal("active_record_automobiles"))
		})

		It("should override model name", func() {
			Ω(motorcycle.ModelName()).Should(Equal("motocicletas"))
		})
	})

	Context("Self", func() {
		It("should get self", func() {
			Ω(automobile.Self()).ShouldNot(BeNil())
		})

		It("should set self", func() {
			automobile.SetSelf(nil)
			Ω(automobile.Self()).Should(BeNil())

			automobile.SetSelf(automobile)
			Ω(automobile.Self()).ShouldNot(BeNil())
		})
	})

	Context("Validation", func() {
		It("should be valid", func() {
			Ω(automobile.Valid()).Should(BeTrue())
			Ω(automobile.HasErrors()).Should(BeFalse())
		})

		It("should be invalid", func() {
			Ω(automobile.Valid()).Should(BeTrue())

			// invalidate
			automobile.Year = 0
			automobile.Make = ""
			automobile.Model = ""

			// validate
			Ω(automobile.Valid()).Should(BeFalse())

			// verify
			Ω(automobile.Errors).ShouldNot(BeNil())
			Ω(automobile.NumErrors()).Should(Equal(3))

			Ω(automobile.ErrorMap()["Year"].Message).Should(Equal("Required"))
			Ω(automobile.ErrorMap()["Make"].Message).Should(Equal("Required"))
			Ω(automobile.ErrorMap()["Model"].Message).Should(Equal("Required"))
		})
	})

	Context("Persistance", func() {
		It("should insert", func() {
			Ω(automobile.Valid()).Should(BeTrue())

			success, err := automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())
		})

		It("should update", func() {
			Ω(automobile.Valid()).Should(BeTrue())

			success, err := automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())

			automobile.Model += " updated"
			success, err = automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())
		})

		It("should delete", func() {
			err := automobile.Delete()
			Ω(err).NotTo(HaveOccurred())
		})
	})

	Context("Query", func() {
		It("should get query", func() {
			q := automobile.Query()
			Ω(q).ShouldNot(BeNil())
		})

		It("should set query", func() {
			q := NewQuery()
			automobile.SetQuery(q)
			Ω(automobile.Query()).ShouldNot(BeNil())
		})

		It("should specify where condition", func() {
			Ω(automobile.Query().WhereConditions).Should(BeNil())
			automobile.Where(QueryCondition{Key: "Year", RelationalOperator: EQ, Value: 2007})
			Ω(automobile.Query().WhereConditions).ShouldNot(BeNil())
			Ω(motorcycle.Query().WhereConditions).Should(BeNil()) // ensure query isn't shared among multiple instances
		})

		It("should sum", func() {
			Ω(len(automobile.Query().Aggregations)).Should(Equal(0))
			automobile.Sum("Year")
			Ω(len(automobile.Query().Aggregations)).Should(Equal(1))
		})

		It("should pluck", func() {
			Ω(automobile.Query().Plucks).Should(BeNil())
			automobile.Pluck("doesn't matter")
			Ω(automobile.Query().Plucks).ShouldNot(BeNil())
		})

		It("should order by", func() {
			Ω(automobile.Query().OrderBys).Should(BeNil())
			automobile.Order(OrderBy{Key: "Year", SortOrder: DESC})
			automobile.Order(OrderBy{Key: "Model"})
			Ω(len(automobile.Query().OrderBys)).Should(Equal(2))
		})

		It("should specify distinct", func() {
			Ω(automobile.Query().Distinct).Should(BeFalse())
			automobile.Distinct()
			Ω(automobile.Query().Distinct).Should(BeTrue())
		})

		It("should run query", func() {
			automobile.Pluck("doesn't matter")
			var results interface{}
			err := automobile.Run(&results)
			Ω(err).NotTo(HaveOccurred())
			Ω(automobile.Query().Plucks).Should(BeNil())
		})
	})

	Context("Timestamps", func() {
		It("should set timestamps after create", func() {
			Ω(automobile.CreatedAt).Should(BeNil())
			Ω(automobile.UpdatedAt).Should(BeNil())
			//Ω(automobile.CreatedAt.IsZero()).Should(BeTrue())

			success, err := automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())

			Ω(automobile.CreatedAt).ShouldNot(BeNil())
			Ω(automobile.UpdatedAt).Should(BeNil())
		})

		It("should set timestamps after update", func() {
			Ω(automobile.CreatedAt).Should(BeNil())
			Ω(automobile.UpdatedAt).Should(BeNil())

			success, err := automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())
			Ω(automobile.CreatedAt).ShouldNot(BeNil())
			Ω(automobile.UpdatedAt).Should(BeNil())
			createdAt := automobile.CreatedAt

			automobile.Model += " updated"
			success, err = automobile.Save()
			Ω(success).Should(BeTrue())
			Ω(err).NotTo(HaveOccurred())
			Ω(automobile.CreatedAt).ShouldNot(BeNil())
			Ω(automobile.UpdatedAt).ShouldNot(BeNil())
			Ω(createdAt).Should(Equal(automobile.CreatedAt))
		})
	})

	Context("Callbacks", func() {
		var (
			callbackErrorModel *CallbackErrorModel
		)

		BeforeEach(func() {
			callbackErrorModel = CallbackErrorModel{}.ToActiveRecord()
		})

		It("should fire a callback", func() {
			e := reflect.ValueOf(automobile).Elem()
			err := Callback("BeforeSave", e.Addr(), nil)
			Ω(err).NotTo(HaveOccurred())
			Ω(automobile.Year).Should(Equal(1888))
		})

		It("should not return an error if callback method is not defined", func() {
			e := reflect.ValueOf(automobile).Elem()
			err := Callback("ImATeapot", e.Addr(), nil)
			Ω(err).NotTo(HaveOccurred())
		})

		It("should raise an error while calling BeforeSave and return the error", func() {
			success, err := callbackErrorModel.Save()
			Ω(err).To(HaveOccurred())
			Ω(success).To(BeFalse())
		})

		It("should raise an error while calling AfterSave and log the error", func() {
			// redefine the beforeSaveCallback method, otherwise we won't reach the AfterSave method
			beforeSaveCallback = func() error {
				return nil
			}

			auto := CallbackErrorModel{Name: "professor moriarty"}.ToActiveRecord()
			success, err := auto.Save()
			Ω(err).NotTo(HaveOccurred())
			Ω(success).To(BeTrue())
		})
	})
})
