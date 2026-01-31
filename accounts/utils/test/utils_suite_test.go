package test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAccountsUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Accounts Utils Suite")
}
