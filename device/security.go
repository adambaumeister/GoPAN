package device

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/zepryspet/GoPAN/utils"
	"strings"
)

type Rule struct {
	Name        string
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

func (r *Rule) Print() {
	fmt.Printf("%v	%v	%v:%v	%v:%v", r.Name, r.Action,
		strings.Join(MembersToText(r.From), ","),
		strings.Join(MembersToText(r.Source), ","),
		strings.Join(MembersToText(r.To), ","),
		strings.Join(MembersToText(r.Destinaton), ","))
}

type RuleBase struct {
	Pre  []Rule
	Post []Rule

	Device []Rule
}

func MembersToText(elements []*etree.Element) []string {
	var result []string
	for _, element := range elements {
		result = append(result, element.Text())
	}
	return result
}

func GetSecurityRules(fqdn string, apikey string, path string) []Rule {
	resp := pan.GetXpath(fqdn, apikey, path)

	rules := []Rule{}
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)

	//fmt.Print(string(resp))
	rulesRoot := doc.FindElement("./response/result/security/rules")
	for _, r := range rulesRoot.SelectElements("entry") {
		rule := Rule{}

		rule.Name = r.SelectAttr("name").Value
		rule.To = r.FindElements("./to/member")
		rule.From = r.FindElements("./from/member")
		rule.Source = r.FindElements("./source/member")
		rule.Service = r.FindElements("./service/member")
		rule.Application = r.FindElements("./application/member")
		rule.Destinaton = r.FindElements("./destination/member")
		rule.Action = r.FindElement("./action").Text()
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
