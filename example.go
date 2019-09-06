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

	res := drill.GetJSON([]byte(report)).FlattenByRule([]byte(rule))
	fmt.Println(res["msg"])
	fmt.Println(res["prov"])
}
