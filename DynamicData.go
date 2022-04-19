package main

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"github.com/pkg/profile"
	"time"
)

var walkoidarray []string

func main() {
	start := time.Now()
	params := g.GoSNMP{
		Target:    "172.16.8.12",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(1) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err1 := params.Connect()
	if err1 != nil {
		fmt.Println("Unable to connect")
	}

	tablelist := make(map[string]string)
	tablelist["ifIndex"] = "1.3.6.1.2.1.2.2.1.1."
	tablelist["ifDescription"] = "1.3.6.1.2.1.2.2.1.2."
	/*tablelist["ifType"] = "1.3.6.1.2.1.2.2.1.3."
	tablelist["ifMtu"] = ".1.3.6.1.2.1.2.2.1.4."*/
	tablelist["OutputStatus"] = ".1.3.6.1.2.1.2.2.1.8."
	tablelist["AdminStatus"] = ".1.3.6.1.2.1.2.2.1.7."

	var walkId = "1.3.6.1.2.1.2.2.1.1"
	defer profile.Start().Stop()
	for i := 0; i < 1; i++ {
		SnmpValue(params, tablelist, walkId)
		fmt.Println(".............................")
		walkoidarray = walkoidarray[:0]
	}
	/*walkIdResult := params.Walk(walkId, WalkFunction)

	if walkIdResult != nil {
		fmt.Println("Error in walkResult")
	}
	var oidDescriptionArray []string

	for k := range tablelist {
		oidDescriptionArray = append(oidDescriptionArray, k)
	}

	var listofoid []string
	var outerMap = make(map[string]interface{})
	for i := 0; i < len(walkoidarray); i++ {
		var innerMap = make(map[string]interface{})
		for oid := range oidDescriptionArray {
			listofoid = append(listofoid, tablelist[oidDescriptionArray[oid]]+walkoidarray[i])
			var result, _ = params.Get(listofoid)
			for _, variables := range result.Variables {
				var data = SnmpData(variables)
				innerMap[oidDescriptionArray[oid]] = data
			}
			listofoid = listofoid[:0]
			outerMap[walkoidarray[i]] = innerMap
		}
	}
	fmt.Println(outerMap)*/

	end := time.Now()
	fmt.Println(end.Sub(start))

}

func SnmpValue(params g.GoSNMP, table map[string]string, walkId string) {

	err1 := params.Connect()
	if err1 != nil {
		fmt.Println("Unable to connect")
	}

	walkIdResult := params.Walk(walkId, WalkFunction)

	if walkIdResult != nil {
		fmt.Println("Error in walkResult")
	}
	var oidDescriptionArray []string

	for key := range table {
		oidDescriptionArray = append(oidDescriptionArray, key)
	}
	var listofoid []string
	var outerMap = make(map[string]interface{})

	for i := 0; i < len(walkoidarray); i++ {
		var innerMap = make(map[string]interface{})
		for oid := range oidDescriptionArray {
			listofoid = append(listofoid, table[oidDescriptionArray[oid]]+walkoidarray[i])
			var result, _ = params.Get(listofoid)

			for _, variables := range result.Variables {
				var data = SnmpData(variables)
				innerMap[oidDescriptionArray[oid]] = data
			}
			listofoid = listofoid[:0]
			outerMap[walkoidarray[i]] = innerMap
		}
	}
	fmt.Println(outerMap)
	walkoidarray = nil

}

func SnmpData(pdu g.SnmpPDU) string {
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

func WalkFunction(pdu g.SnmpPDU) error {
	switch pdu.Type {
	case g.IPAddress:
		result := pdu.Value
		walkoidarray = append(walkoidarray, result.(string))
		break
	case g.Integer:
		result := g.ToBigInt(pdu.Value)
		walkoidarray = append(walkoidarray, result.String())
		break
	case g.OctetString:
		result := pdu.Value.([]byte)
		walkoidarray = append(walkoidarray, string(result))
		break
	default:
		result := pdu.Value
		walkoidarray = append(walkoidarray, result.(string))
	}
	return nil
}
