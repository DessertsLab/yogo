package drill

import (
	"fmt"
	"reflect"
	"testing"
)

var DATAJSON3 = `[{"msg": "成功101"},{"msg": "失败404"}]`

var DATAJSON1 = `{
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

var DATAJSON2 = `{"list3": [
            "one eggs on the desk: 1",
            "two eggs on the desk: 2",
            "one eggs on the desk: 3"
        ]}`

var RULEJSON1 = `
{
	"msg":["msg"],
 	"prov":["data","ISPNUM","province"]
}
`

func TestGetJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want Data
	}{
		{"case1", args{[]byte("abc")}, Data{[]byte{97, 98, 99}}},
		{"case2", args{[]byte("")}, Data{[]byte{}}},
		{"case3", args{[]byte(" ")}, Data{[]byte{32}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetJSON(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_rawtype(t *testing.T) {
	case2str := `{
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

	case3str := `"ECL": [
            {
                "code": "10000004",
                "IFT": "A3"
            }
		]`

	case4str := `[
            {
                "code": "10000004",
                "IFT": "A3"
            }
        ]`

	type fields struct {
		raw []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"case1", fields{[]byte("{key:value}")}, "map"},
		{"case2", fields{[]byte(case2str)}, "map"},
		{"case3", fields{[]byte(case3str)}, "string"},
		{"case4", fields{[]byte(case4str)}, "list"},
		{"case5", fields{[]byte("")}, "empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				raw: tt.fields.raw,
			}
			if got := d.rawtype(); got != tt.want {
				t.Errorf("Data.rawtype() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_reduce(t *testing.T) {
	type fields struct {
		raw []byte
	}
	type args struct {
		r []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Data
	}{
		{"case1", fields{[]byte(`{"hello":"world"}`)}, args{[]byte(`"hello"`)}, Data{[]byte(`"world"`)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				raw: tt.fields.raw,
			}
			fmt.Println(tt.args.r)
			d.reduce(tt.args.r)
			if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("after reduce data =  %v, but want %v", d, tt.want)
			}
		})
	}
}

func TestData_get(t *testing.T) {
	type fields struct {
		raw []byte
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Data
	}{
		{"case1", fields{[]byte(`{"hello":"world"}`)}, args{"hello"}, Data{[]byte(`"world"`)}},
		// FIXME: there is a quote at the end
		{"case2", fields{[]byte(`"互联网金融门户:1"`)}, args{"互联网金融门户.*"}, Data{[]byte(`互联网金融门户:1"`)}},
		{"case3", fields{[]byte(`互联网金融门户:1`)}, args{"\\d+$"}, Data{[]byte(`1`)}},
		{"case4", fields{[]byte("")}, args{"\\d+$"}, Data{[]byte(``)}},
		{"case5", fields{[]byte(`{"hello":"world"}`)}, args{"notthere"}, Data{[]byte("")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				raw: tt.fields.raw,
			}
			if got := d.get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_FlattenByRule(t *testing.T) {
	type fields struct {
		raw []byte
	}
	type args struct {
		rule [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		// TODO: Add test cases.
		{"case1_rulejson1", fields{[]byte(DATAJSON1)}, args{[][]byte{[]byte(RULEJSON1)}}, map[string]string{"msg": "成功", "prov": "湖北"}},
		{"case2_norule", fields{[]byte(DATAJSON1)}, args{}, map[string]string{"empty": ""}},
		{"case3_emptyrule", fields{[]byte(DATAJSON1)}, args{[][]byte{}}, map[string]string{"empty": ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				raw: tt.fields.raw,
			}
			if got := d.FlattenByRule(tt.args.rule...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.FlattenByRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_search(t *testing.T) {
	type fields struct {
		raw []byte
	}
	type args struct {
		keyvalue []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Data
	}{
		// TODO: Add test cases.
		{"case1", fields{[]byte(DATAJSON2)}, args{[]byte(`["list3", "GETLIST", "one eggs on the desk.*"]`)}, Data{[]byte(`one eggs on the desk: 1`)}},
		{"case2", fields{[]byte(DATAJSON3)}, args{[]byte(`["msg", "CONTAINS", "成功"]`)}, Data{[]byte(`{"msg": "成功101"}`)}},
		{"case3", fields{[]byte(DATAJSON3)}, args{[]byte(`["msg", "失败404"]`)}, Data{[]byte(`{"msg": "失败404"}`)}},
		{"case4", fields{[]byte(DATAJSON3)}, args{[]byte(`["msg"]`)}, Data{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				raw: tt.fields.raw,
			}
			if got := d.search(tt.args.keyvalue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.search() = %v, want %v", string(got.raw[:]), string(tt.want.raw[:]))
			}
		})
	}
}

func Test_drill(t *testing.T) {
	type args struct {
		data Data
		rule Rule
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"case1", args{Data{[]byte(DATAJSON1)}, Rule{[]byte(`"code"`)}}, `200`},
		{"case2", args{Data{[]byte(`"find the number: 101"`)}, Rule{[]byte(`".*"`)}}, `"find the number: 101"`},
		{"case3", args{Data{[]byte(DATAJSON3)}, Rule{[]byte(`["msg","CONTAINS","101"]`)}}, ``},
		{"case4", args{Data{[]byte(DATAJSON3)}, Rule{[]byte(`?"msg","CONTAINS","101"]`)}}, ``},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := drill(tt.args.data, tt.args.rule); got != tt.want {
				t.Errorf("drill() = %v, want %v", got, tt.want)
			}
		})
	}
}
