package device

import (
	"github.com/beevik/etree"
	"github.com/zepryspet/GoPAN/utils"
)

type Rule struct {
	To          []*etree.Element
	From        []*etree.Element
	Source      []*etree.Element
	Destinaton  []*etree.Element
	SourceUser  []*etree.Element
	Category    []*etree.Element
	Application []*etree.Element
	Service     []*etree.Element
	HipProfiles []*etree.Element
	Action      string
}

type RuleBase struct {
	Pre  []Rule
	Post []Rule

	Device []Rule
}

func GetSecurityRules(fqdn string, apikey string, path string) []Rule {

	//path := "/config/devices/entry[@name='localhost.localdomain']/vsys/entry[@name='vsys1']/rulebase/security"
	resp := pan.GetXpath(fqdn, apikey, path)

	rules := []Rule{}
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)

	//fmt.Print(string(resp))
	rulesRoot := doc.FindElement("./response/result/security/rules")
	for _, r := range rulesRoot.SelectElements("entry") {
		rule := Rule{}

		rule.To = r.FindElements("./to/member")
		rule.From = r.FindElements("./from/member")
		rule.Source = r.FindElements("./source/member")
		rule.Destinaton = r.FindElements("./destination/member")
		rule.SourceUser = r.FindElements("./source-user/member")

		rules = append(rules, rule)
	}

	return rules
}

func GetElementsText(elements []*etree.Element) []string {
	r_elements := []string{}
	for _, element := range elements {
		r_elements = append(r_elements, element.Text())
	}
	return r_elements
}
