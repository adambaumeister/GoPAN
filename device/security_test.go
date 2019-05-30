package device

import (
	"fmt"
	"github.com/zepryspet/GoPAN/utils"
	"strings"
	"testing"
)

/*
Test suite for the object-style firewall rulebase functions

To use this suite, you can run "go test" inside this package directory.

It requires a Firewall or Panorama to be accessible at TESTS_FW_ADDR with default credentials.
*/

func TestGetSecurityRules(t *testing.T) {
	fw := Connect("localhost:8443", "admin", "admin")
	fw.GetRules()

	for _, rule := range fw.Rules {
		sources := GetElementsText(rule.Source)
		fmt.Printf("%v\n", strings.Join(sources, ", "))
	}
}

func TestCmdGenSlice(t *testing.T) {
	fmt.Print(pan.KvCmdGenSlice([]string{"f", "v"}))
}
