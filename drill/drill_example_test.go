package drill

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func Example() {
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
    "key": [
            "one eggs on the desk: 1",
            "two eggs on the desk: 2",
            "one eggs on the desk: 3"
        ],
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

	gjexp := `{
    "name": {
        "first": "Tom",
        "last": "Anderson"
    },
    "age": 37,
    "children": ["Sara", "Alex", "Jack"],
    "fav.movie": "Deer Hunter",
    "friends": [{
            "first": "Dale",
            "last": "Murphy",
            "age": 44,
            "nets": ["ig", "fb", "tw"]
        },
        {
            "first": "Roger",
            "last": "Craig",
            "age": 68,
            "nets": ["fb", "tw"]
        },
        {
            "first": "Jane",
            "last": "Murphy",
            "age": 47,
            "nets": ["ig", "tw"]
        }
    ]
}`

	res := GetJSON([]byte(report)).FlattenByRule([]byte(rule))
	res2 := GetJSON([]byte(s)).FlattenByRule([]byte(rule2))
	res4 := gjson.Get(gjexp, "friends.2.nets.0")

	fmt.Println(res["msg"])
	fmt.Println(res["prov"])
	fmt.Println(res2)
	fmt.Println(res4)
	// Output:
	// 成功
	// 湖北
	// map[col1: col2:eeee col3: col4:]
	// ig

}
