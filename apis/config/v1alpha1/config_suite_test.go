package v1alpha1_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config API Test Suite")
}
