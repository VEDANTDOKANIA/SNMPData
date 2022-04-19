package main

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"github.com/pkg/profile"
	"time"
)

func main() {
	params := &g.GoSNMP{
		Target:    "172.16.8.12",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(1) * time.Second,
	}
	err1 := params.Connect()
	if err1 != nil {
	}
	defer profile.Start().Stop()

	var oidlist []string
	oidlist = append(oidlist, ".1.3.6.1.2.1.2.2.1.1.1")
	for i := 0; i < 500; i++ {

		var result, _ = params.Get(oidlist)
		fmt.Println(result)
	}

	/*for _, variables := range result.Variables {
		var data = SnmpData1(variables)
		fmt.Println(data)
	}*/

}
func SnmpData1(pdu g.SnmpPDU) string {
	var result string
	if pdu.Value == nil {
		return "empty"
	}
	switch pdu.Type {
	case g.IPAddress:
		result = pdu.Value.(string)
		break
	case g.Integer:
		result = g.ToBigInt(pdu.Value).String()
		break

	case g.OctetString:
		result = string(pdu.Value.([]byte))
		break
	default:
		result = pdu.Value.(string)
	}

	return result
}
