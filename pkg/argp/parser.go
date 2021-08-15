package argp

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

var argsRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)

type Parser struct {
	args []string
}

func New(args ...[]string) (p *Parser) {
	p = &Parser{
		args: os.Args[1:],
	}
	if len(args) > 0 {
		p.args = args[0]
	}
	p.args = resplit(p.args)
	return
}

func (p *Parser) Scan(param string, val interface{}) (ok bool, err error) {
	var (
		arg   string
		sval  string
		i     int
		pad   int
		found bool
	)

	for i, arg = range p.args {
		if strings.HasPrefix(arg, param) {
			found = true
			break
		}
	}
	if !found {
		return
	}

	if _, isBool := val.(*bool); isBool && len(arg) == len(param) {
		arg += "=true"
	}

	if len(arg) == len(param) {
		if len(p.args) < i+2 {
			return
		}
		sval = p.args[i+1]
		pad++
	} else {
		split := strings.SplitN(arg, "=", 2)
		if len(split) != 2 {
			return
		}
		sval = split[1]
	}

	if _, isStr := val.(*string); isStr {
		sval = "\"" + sval + "\""
	}

	err = json.Unmarshal([]byte(sval), val)
	ok = err == nil

	if ok {
		p.args = append(p.args[:i], p.args[i+1+pad:]...)
	}

	return
}

func (p *Parser) String(param string, def ...string) (val string, err error) {
	ok, err := p.Scan(param, &val)
	if err != nil {
		return
	}
	if !ok && len(def) > 0 {
		val = def[0]
	}
	return
}

func (p *Parser) Bool(param string, def ...bool) (val bool, err error) {
	ok, err := p.Scan(param, &val)
	if err != nil {
		return
	}
	if !ok && len(def) > 0 {
		val = def[0]
	}
	return
}

func (p *Parser) Int(param string, def ...int) (val int, err error) {
	ok, err := p.Scan(param, &val)
	if err != nil {
		return
	}
	if !ok && len(def) > 0 {
		val = def[0]
	}
	return
}

func (p *Parser) Float(param string, def ...float64) (val float64, err error) {
	ok, err := p.Scan(param, &val)
	if err != nil {
		return
	}
	if !ok && len(def) > 0 {
		val = def[0]
	}
	return
}

func (p *Parser) Args() []string {
	return p.args
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
