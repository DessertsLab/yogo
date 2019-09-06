package drill

import (
	"fmt"
	"reflect"
	"testing"
)

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
		// TODO: Add test cases.
		{"case1", fields{[]byte(`{"hello":"world"}`)}, args{"hello"}, Data{[]byte(`"world"`)}},
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
