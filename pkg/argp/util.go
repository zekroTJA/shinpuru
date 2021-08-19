package argp

import "strings"

func optStr(v []string, def string) (r string) {
	r = def
	if len(v) > 0 {
		r = v[0]
	}
	return
}

func split(v string) (split []string) {
	split = argsRx.FindAllString(v, -1)
	if len(split) == 0 {
		return
	}

	for i, k := range split {
		if strings.Contains(k, "\"") {
			split[i] = strings.Replace(k, "\"", "", -1)
		}
	}

	return
}

func resplit(args []string) []string {
	join := strings.Join(args, " ")
	return split(join)
}

func typName(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case float32:
	case float64:
		return "float"
	case bool:
		return "bool"
	}
	return "<unknown type>"
}
