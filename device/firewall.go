package device

import (
	"github.com/zepryspet/GoPAN/utils"
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
