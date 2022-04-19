package main

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"time"
)

func main() {
	start := time.Now()
	params := g.GoSNMP{
		Target:    "172.16.8.2",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}

	table1 := make(map[string]string)
	table1["ifIndex"] = ".1.3.6.1.2.1.2.2.1.1."
	table1["ifDescription"] = "1.3.6.1.2.1.2.2.1.2."
	table1["OutputStatus"] = "1.3.6.1.2.1.2.2.1.8."
	table1["AdminStatus"] = "1.3.6.1.2.1.2.2.1.7."
	table2 := make(map[string]string)
	table2["column1"] = "1.3.6.1.2.1.4.20.1.2."
	table2["column2"] = "1.3.6.1.2.1.4.20.1.3."
	table2["column3"] = "1.3.6.1.2.1.4.20.1.4."
	table2["column4"] = "1.3.6.1.2.1.4.20.1.5."

	var walkId = "1.3.6.1.2.1.2.2.1.1"
	//defer profile.Start().Stop()
	c := make(chan []interface{}, 1)
	for i := 0; i < 3; i++ {
		go SnmpValue(params, table1, walkId, c)
		go SnmpValue(params, table2, "1.3.6.1.2.1.4.20.1.1", c)
	}
	for i := 0; i < 6; i++ {
		fmt.Println(<-c)
		fmt.Println("...................")
	}
	end := time.Now()
	fmt.Println(end.Sub(start))

}

func SnmpValue(params g.GoSNMP, table map[string]string, walkId string, c chan []interface{}) {
	err1 := params.Connect()
	if err1 != nil {
		fmt.Println("Unable to connect")
	}
	var walkOidArray []string
	walkIdResult := params.Walk(walkId, func(pdu g.SnmpPDU) error {
		switch pdu.Type {
		case g.IPAddress:
			result := pdu.Value
			walkOidArray = append(walkOidArray, result.(string))
			break
		case g.Integer:
			result := g.ToBigInt(pdu.Value)
			walkOidArray = append(walkOidArray, result.String())
			break
		case g.OctetString:
			result := pdu.Value.([]byte)
			walkOidArray = append(walkOidArray, string(result))
			break
		default:
			result := pdu.Value
			walkOidArray = append(walkOidArray, result.(string))
		}
		return nil
	},
	)

	if walkIdResult != nil {
		fmt.Println("Error in walkResult")
	}
	var oidDescriptionArray []string

	for key := range table {
		oidDescriptionArray = append(oidDescriptionArray, key)
	}
	var listOfOids []string

	var resultArray []interface{}
	for i := 0; i < len(walkOidArray); i++ {
		for oid := range oidDescriptionArray {
			listOfOids = append(listOfOids, table[oidDescriptionArray[oid]]+walkOidArray[i])
		}
	}

	var startIndex = 0
	var endIndex = 40

	for {
		var result, err = params.Get(listOfOids[startIndex:endIndex])
		if err != nil {
			fmt.Println(err)
		}
		for _, variable := range result.Variables {
			resultArray = append(resultArray, SnmpData(variable))
		}
		startIndex = endIndex + 1
		endIndex = endIndex + 60

		if endIndex > len(listOfOids) {
			endIndex = len(listOfOids)
		}
		if startIndex == len(listOfOids)+1 {
			break
		}
	}
	c <- resultArray

}

func SnmpData(pdu g.SnmpPDU) interface{} {

	/*if pdu.Value == nil {
		return "empty"
	}*/
	switch pdu.Type {
	case g.IPAddress:
		return pdu.Value
	case g.Integer:
		return g.ToBigInt(pdu.Value)

	case g.OctetString:
		return string(pdu.Value.([]byte))
	default:
		return pdu.Value
	}
	
}
