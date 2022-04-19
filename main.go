package main

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

var ipNettableoid []string
var interfaceTableoid []string

func main() {
	start := time.Now()
	params := &g.GoSNMP{
		Target:    "172.16.8.12",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(3) * time.Second,
	}
	err1 := params.Connect()
	if err1 != nil {
		fmt.Println("Unable to connect")

	}
	defer params.Conn.Close()
	var ipNetToMediaIfIndex = ".1.3.6.1.2.1.4.22.1.1.79."
	var ipNetToMediaPhysicalAddress = ".1.3.6.1.2.1.4.22.1.2.79."
	var ipNetToMediaType = ".1.3.6.1.2.1.4.22.1.3.79."
	var ipNetToMediaNetAddress = ".1.3.6.1.2.1.4.22.1.4.79."
	var Indexvalue = "1.3.6.1.2.1.4.22.1.3" // Isspe snmp walk chalyenga

	interfaceIndex := params.Walk(Indexvalue, WalkFunction1)
	if interfaceIndex != nil {
		fmt.Println("Walk Function failed")
	}
	var oidList []string

	for i := 0; i < len(ipNettableoid); i++ {
		oidList = append(oidList, (ipNetToMediaIfIndex + ipNettableoid[i]), (ipNetToMediaPhysicalAddress + ipNettableoid[i]))
		oidList = append(oidList, (ipNetToMediaType + ipNettableoid[i]), (ipNetToMediaNetAddress + ipNettableoid[i]))
	}

	var NetTableResult, _ = params.Get(oidList)
	netmap := make(map[string]interface{})
	var count = 0
	var index = 0
	innerMap := make(map[string]interface{})

	for _, variable := range NetTableResult.Variables {

		VariableName := strings.SplitAfter(variable.Name, ".1.3.6.1.2.1.4.22.1.")
		RootOid := VariableName[0] + strings.Split(VariableName[1], ".")[0] + "."
		if count == 4 {

			netmap[ipNettableoid[index]] = innerMap
			innerMap = make(map[string]interface{})
			index++
			count = 0
		}
		switch RootOid {
		case ".1.3.6.1.2.1.4.22.1.1.":
			value := g.ToBigInt(variable.Value)
			innerMap["IpNetToMediaIfIndex"] = value
		case ".1.3.6.1.2.1.4.22.1.2.":
			value := fmt.Sprintf("%x", variable.Value)
			innerMap["IfNetTOmediaPhyAddress"] = value
		case ".1.3.6.1.2.1.4.22.1.3.":
			value := variable.Value
			innerMap["ipNetToMediaNet"] = value
		case ".1.3.6.1.2.1.4.22.1.4.":
			value := variable.Value
			innerMap["ipNetToMediaNetAddress"] = value
		}
		count++
	}
	fmt.Println(netmap)
	end := time.Now()
	fmt.Println(end.Sub(start))

}
func WalkFunction1(pdu g.SnmpPDU) error {
	switch pdu.Type {
	case g.IPAddress:
		result := pdu.Value
		interfaceTableoid = append(interfaceTableoid, result.(string))
		break
	case g.Integer:
		result := g.ToBigInt(pdu.Value)
		interfaceTableoid = append(interfaceTableoid, result.String())
		break
	case g.OctetString:
		result := pdu.Value.([]byte)
		interfaceTableoid = append(interfaceTableoid, string(result))
		break
	default:
		result := pdu.Value
		interfaceTableoid = append(interfaceTableoid, result.(string))
	}

	return nil
}
