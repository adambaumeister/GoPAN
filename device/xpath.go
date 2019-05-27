package device

import (
	"fmt"
	"strings"
)

const PANOS_ROOT_XPATH = "/config/devices"

func (fw *Firewall) MakeXPath(path []string) string {
	device := fmt.Sprintf("entry[@name='%v']", fw.Device)
	vsys := fmt.Sprintf("vsys/entry[@name='%v']", fw.Vsys)
	xpList := []string{
		PANOS_ROOT_XPATH,
		device,
		vsys,
	}

	xpList = append(xpList, path...)

	return fmt.Sprintf("%v", strings.Join(xpList, "/"))
}
