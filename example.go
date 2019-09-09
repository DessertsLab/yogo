package main

import (
	"fmt"

	"github.com/DessertsLab/yogo/drill"
)

func main() {
	report := `{
    "code": "200",
    "data": {
        "ISPNUM": {
            "province": "湖北",
            "city": "武汉",
            "isp": "移动"
        },
        "RSL": [
            {
                "RS": {
                    "code": "0",
                    "desc": "三维验证一致"
                },
                "IFT": "B7"
            }
        ],
        "ECL": [
            {
                "code": "10000004",
                "IFT": "A3"
            }
        ]
    },
    "msg": "成功"
}`

	rule := `
        {
            "msg":["msg"],
            "prov":["data","ISPNUM","province"]
        }
        `

	s := `{
    "key": {
        "list": [{
            "one": 1,
            "qq": "2244"
        }, {
            "two": 2,
            "qq": "e32444"
        }],
        "list2": [{
            "three": 3,
            "qq": "44455"
        }],
        "list3": [
            "one eggs on the desk: 1",
            "two eggs on the desk: 2",
            "one eggs on the desk: 3"
        ],
        "detail": "feast on your life 345556"
    },
    "key2": "eeee"
}`

	rule2 := `
    {
        "col1":["key","list",["one",1],"qq"],
        "col2":["key2"],
        "col3":["key",["list3","GETLIST","one eggs on the desk: .*"],"\\d+$"],
        "col4":["key","list",["qq","CONTAINS","32"],"qq","^.*2"]
}
    `

	res := drill.GetJSON([]byte(report)).FlattenByRule([]byte(rule))
	fmt.Println(res["msg"])
	fmt.Println(res["prov"])

	res2 := drill.GetJSON([]byte(s)).FlattenByRule([]byte(rule2))
	fmt.Println(res2)

}
