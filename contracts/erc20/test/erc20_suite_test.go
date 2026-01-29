package erc20_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestERC20(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ERC20 Suite")
}
