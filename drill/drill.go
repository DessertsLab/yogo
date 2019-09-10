package drill

import (
	"encoding/json"
	"fmt"
	"regexp"
)

//https://stackoverflow.com/questions/24293790/golang-decoding-unmarshaling-invalid-unicode-in-json

// Map is a json struct
type mapRaw map[string]json.RawMessage

type mapListRaw map[string][]json.RawMessage

type mapString map[string]string

type listRaw []json.RawMessage

type listStr []string

// Data ...
type Data struct {
	raw []byte
}

// Rule ...
type Rule []json.RawMessage

type mapRule map[string]Rule

// rawtype return raw datatype of Data
func (d Data) rawtype() string {
	if len(d.raw) == 0 {
		return "empty"
	}
	first := string(d.raw)[0]
	if first == '{' {
		return "map"
	} else if first == '[' {
		return "list"
	} else {
		return "string"
	}
}

// reduce Data to subset Data
func (d *Data) reduce(r []byte) {
	var val string
	json.Unmarshal(r, &val)
	*d = d.get(val)
}

func (d *Data) get(key string) Data {
	if d.rawtype() == "string" {
		var id = regexp.MustCompile(key)
		res := id.FindString(string(d.raw))
		return Data{[]byte(res)}
	} else if d.rawtype() == "empty" {
		return Data{[]byte("")}
	} else {
		mr := mapRaw{}
		json.Unmarshal(d.raw, &mr)
		if mr[key] == nil {
			fmt.Println(key, "could not found in data")
			return Data{[]byte("")}
		}
		return Data{mr[key]}
	}
}

func (d *Data) search(keyvalue []byte) Data {
	ls := listStr{}
	lr := listRaw{}
	mr := mapRaw{}
	var val string
	json.Unmarshal(keyvalue, &ls)
	json.Unmarshal(d.raw, &lr)
	if len(ls) == 2 {
		for _, jsonItem := range lr {
			json.Unmarshal(jsonItem, &mr)
			json.Unmarshal(mr[ls[0]], &val)
			if val == ls[1] {
				return Data{jsonItem}
			}
		}
	} else if len(ls) == 3 {
		if ls[1] == "CONTAINS" {
			for _, jsonItem := range lr {
				json.Unmarshal(jsonItem, &mr)
				json.Unmarshal(mr[ls[0]], &val)
				var id = regexp.MustCompile(ls[2])
				if id.MatchString(val) {
					return Data{jsonItem}
				}
			}
		} else if ls[1] == "GETLIST" {
			jsonmr := mapRaw{}
			json.Unmarshal(d.raw, &jsonmr)
			lstr := listStr{}
			json.Unmarshal(jsonmr[ls[0]], &lstr)
			var id = regexp.MustCompile(ls[2])
			for _, s := range lstr {
				if id.MatchString(s) {
					return Data{[]byte(s)}
				}
			}
		}
	}
	return Data{}
}

func drill(data Data, rule Rule) string {

	for _, ruleItem := range rule {
		firstChar := string(ruleItem)[0]
		switch firstChar {
		/* data.raw should be a map and your rule should be string */
		case '"':
			if data.rawtype() == "string" {
				var val string
				json.Unmarshal(ruleItem, &val)
				var id = regexp.MustCompile(string(val))
				res := id.Find(data.raw)
				return string(res)
			}

			data.reduce(ruleItem)

			/* data.raw should be a list and your rule is a
			two length list rep key value to search in data
			*/
		case '[':
			data = data.search(ruleItem)
		default:
			data = Data{[]byte("")}
		}
	}

	var res string
	json.Unmarshal(data.raw, &res)

	return res
}

//FlattenByRule flat the json struct to mapstring
func (d Data) FlattenByRule(rule ...[]byte) map[string]string {
	if len(rule) == 0 {
		//TODO: auto flatten with no rules
		// fmt.Println("todo: auto flatten with no rules")
		return map[string]string{"empty": ""}
	}
	row := map[string]string{}

	mr := mapRule{}
	json.Unmarshal(rule[0], &mr)

	for colName, singleRule := range mr {
		row[colName] = drill(d, singleRule)
	}
	return row
}

// GetJSON bytes array of json string and  return Data type
func GetJSON(b []byte) Data {
	return Data{b}
}

/*

 */
