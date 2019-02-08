package cps

import (
	"fmt"
	g "github.com/soniah/gosnmp"
	"errors"
	"os"
	"strconv"
	"github.com/zepryspet/GoPAN/utils"
	"time"
)

var zones []string
var cps []int

//function to SNMP walk on the  CPS OIDs and save the output into text files based on the zone name.
func Snmpgen(fqdn string, community string, seconds int, version int, authv3 string, privacyv3 string) {
	//Validating SNMPv3 password and snmp version
	if version == 3{
		if authv3 =="" || privacyv3 ==""{
			e := "empty snmpv3 password"
			pan.Logerror(errors.New(e), true)
		}
	}else if version !=2{
		e := "non-supported snmp version"
		pan.Logerror(errors.New(e), true)
	}
	sec := time.Duration( seconds)
	loop := true
	for loop {
		port, _ := strconv.ParseUint("161", 10, 16)
		params := &g.GoSNMP{}
		if version ==2{
			params = &g.GoSNMP{
				Target:    fqdn,
				Port:      uint16(port),
				Community: community,
				Version:   g.Version2c,
				Timeout:   time.Duration(10) * time.Second,
				Retries:	3,
				//Logger:    log.New(os.Stdout, "", 0), removed verbose output, too much noise
			}
		}else if version ==3{
			//fmt.Println(passv3)
			params = &g.GoSNMP{
				Target:        fqdn,
				Port:          uint16(port),
				Version:       g.Version3,
				Timeout:       time.Duration(10) * time.Second,
				SecurityModel: g.UserSecurityModel,
				MsgFlags:      g.AuthPriv,
				Retries:	3,
				SecurityParameters: &g.UsmSecurityParameters{
					UserName: community,
					AuthenticationProtocol:   g.SHA,
					AuthenticationPassphrase: authv3,
					PrivacyProtocol:          g.AES,
					PrivacyPassphrase:        privacyv3,
				},
			}
		}
		err := params.Connect()
		if err != nil {
			fmt.Println("Error establishing SNMP connection to firewall. "+ err.Error())
			pan.Logerror(err, true)
		}
		defer params.Conn.Close()
		//oids := []string{".1.3.6.1.4.1.25461.2.1.2.3.10.1.1", ".1.3.6.1.4.1.25461.2.1.2.3.10.1.1.5.84.114.117.115.116"}
		var oid [3]string
		//zone names OID
		oidZone := ".1.3.6.1.4.1.25461.2.1.2.3.10.1.1"
		//TCP CPS per zone oID, returns unsigned 32bit integer
		oid[0] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.2"
		//UDP CPS per zone oID, returns unsigned 32bit integer
		oid[1] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.3"
		//other IP CPS per zone oID, returns unsigned 32bit integer
		oid[2] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.4"

		err = params.BulkWalk(oidZone, saveZone)
		if err != nil {
			fmt.Println("Error doing SNMP bulkwalk to get zone names. "+ err.Error())
			pan.Logerror(err, true)
		}
		fmt.Printf("\n%v", zones)

		for _, element := range oid {
			err = params.BulkWalk(element, saveCPS)
			if err != nil {
				fmt.Println("Error doing SNMP bulkwalk to get CPS. "+ err.Error())
				pan.Logerror(err, true)
			}
			fmt.Printf("%v", cps)
			//Saving CPS info to the respective zone name using CSV
			for i := range zones {
				pan.Wlog(zones[i]+".csv", strconv.Itoa(cps[i])+",", false)
			}
			//Clearing the slice to reuse it for the next SNMP  walk
			cps = nil
		}
		//Ading time and a breakline at the end of all zone statistic CPS zone files
		for i := range zones {
			t := time.Now()
			pan.Wlog(zones[i]+".csv", t.Format(time.UnixDate) +"\n", false)
		}
		//clearing slice for re-use
		zones = nil
		time.Sleep(sec * time.Second)
	}
}

func saveZone(pdu g.SnmpPDU) error {
	//fmt.Printf("%s = ", pdu.Name)
	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.([]byte)
		zones = append(zones, string(b))
	default:
		fmt.Println("received zone name NOT as string, please check the OIDs for CPS zone names")
		pan.Wlog("error.txt", "received zone name NOT as string, please check the OIDs for CPS zone names", true)
		os.Exit(1)
	}
	return nil
}

func saveCPS(pdu g.SnmpPDU) error {
	//fmt.Printf("%s = ", pdu.Name)
	switch pdu.Type {
	case g.OctetString:
		fmt.Println("received CPS as string instead as int, please check the OIDs for CPS")
		pan.Wlog("error.txt", "received CPS as string instead as int, please check the OIDs for CPS", true)
		os.Exit(1)
	default:
		cps = append(cps, pdu.Value.(int))
	}
	return nil
}
