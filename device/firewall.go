package device

import (
	"encoding/xml"
	"fmt"
	"github.com/zepryspet/GoPAN/utils"
	"regexp"
)

type Firewall struct {
	// Basic stuff, required
	Fqdn   string
	User   string
	Pass   string
	Apikey string

	// Optional
	Vsys   string
	Device string

	// Data
	Rules []Rule
}

func Connect(fqdn string, user string, pass string) *Firewall {
	/*
		Connect to a Firewall and return it's containing Struct
	*/

	fw := Firewall{
		Fqdn:   fqdn,
		Vsys:   "vsys1",
		Device: "localhost.localdomain",
	}
	fw.Apikey = pan.Keygen(fqdn, user, pass)
	fw.User = user
	fw.Pass = pass
	return &fw
}

func (fw *Firewall) GetRules() *RuleBase {
	/*
		Get the rulebase for this device
		As a firewall, it only has one set of rules, global to the device
	*/
	rb := RuleBase{}

	path := fw.MakeXPath([]string{"rulebase", "security"})
	rb.Device = GetSecurityRules(fw.Fqdn, fw.Apikey, path)

	fw.Rules = rb.Device
	return &rb
}

func (fw *Firewall) SearchRules(query string) []Rule {
	fw.GetRules()
	result := []Rule{}
	for _, r := range fw.Rules {
		m, _ := regexp.Match(query, []byte(r.Name))
		if m {
			result = append(result, r)
		}
	}
	return result
}

func (fw *Firewall) SearchAndPrint(context string, query string) {

	switch context {
	case "rules":
		result := fw.SearchRules(query)
		for _, r := range result {
			r.Print()
		}
	}
}

func (fw *Firewall) Test(query []string) {
	testParams := pan.KvCmdGenSlice(query[1:])
	test := fmt.Sprintf("<test><%v>%v</%v></test>", query[0], testParams, query[0])
	resp := pan.RunOp(fw.Fqdn, fw.Apikey, test)

	v := Response{}
	xml.Unmarshal(resp, &v)
	//fmt.Print(string(resp))
	fmt.Printf("From: %v\n", v.Result.Rules.Entry[0].Name)
}

type Response struct {
	Cmd    string     `xml:"cmd,attr"`
	Status string     `xml:"status,attr"`
	Result TestResult `xml:"result"`
}
type TestResult struct {
	Rules TestResultRules `xml:"rules"`
}

type TestResultRules struct {
	Entry []TestResultEntry `xml:"entry"`
}
type TestResultEntry struct {
	Name string `xml:"name,attr"`
}
