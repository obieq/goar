package goar

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestActiveRecord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ActiveRecord Suite")
}
