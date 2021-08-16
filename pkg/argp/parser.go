package argp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var argsRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)

type flag struct {
	flag string
	def  interface{}
	help string
}

// Parser takes an array of arguments and provides
// functionalities to parse flags and values contained.
type Parser struct {
	args  []string
	flags []flag
}

// New initializes a new instance of Parser.
//
// It defaultly takes the value of os.Args[1:]
// as array of arguments. Optionally, you can
// pass a custom array of arguments you want
// to scan.
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

// Scan looks for the passed flag (unprefixed) in
// the arguments array. If the flag was found, the
// value of the flag is scanned into the pointer
// of val. If the flag and value was found and valid,
// true is returned. Otherwise, false is returned and
// if an error occurs, the error is returned as well.
//
// Example:
//   var config string
//   p := argp.New([]string{"--config", "myconfig.yml"})
//   ok, err := p.Scan("--config", &config)
//   // config = "myconfig.yml"
//   // ok     = true
//   // err    = nil
func (p *Parser) Scan(flag string, val interface{}) (ok bool, err error) {
	var (
		arg   string
		sval  string
		i     int
		pad   int
		found bool
	)

	for i, arg = range p.args {
		if strings.HasPrefix(arg, flag) {
			found = true
			break
		}
	}
	if !found {
		return
	}

	if _, isBool := val.(*bool); isBool && len(arg) == len(flag) {
		arg += "=true"
	}

	if len(arg) == len(flag) {
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

// String is shorthand for Scan with a string flag value.
// It returns the scanned value and an error if the parsing
// failed. If no flag or value was found and a def value was
// passed, def is returned as val.
func (p *Parser) String(flag string, def string, help ...string) (val string, err error) {
	ok, err := p.Scan(flag, &val)
	if err != nil {
		return
	}
	if !ok {
		val = def
	}
	p.register(flag, help, def)
	return
}

// Bool is shorthand for Scan with a bool flag value.
// If the flag was passed (with or wirhout value specified),
// true is returned. If the parsing fails, the error is
// returned. When def is passed and no flag was found, def
// is returned as val.
func (p *Parser) Bool(flag string, def bool, help ...string) (val bool, err error) {
	ok, err := p.Scan(flag, &val)
	if err != nil {
		return
	}
	if !ok {
		val = def
	}
	p.register(flag, help, def)
	return
}

// Int is shorthand for Scan with a integer flag value.
// It returns the scanned value and an error if the parsing
// failed. If no flag or value was found and a def value was
// passed, def is returned as val.
func (p *Parser) Int(flag string, def int, help ...string) (val int, err error) {
	ok, err := p.Scan(flag, &val)
	if err != nil {
		return
	}
	if !ok {
		val = def
	}
	p.register(flag, help, def)
	return
}

// Float is shorthand for Scan with a float flag value.
// It returns the scanned value and an error if the parsing
// failed. If no flag or value was found and a def value was
// passed, def is returned as val.
func (p *Parser) Float(flag string, def float64, help ...string) (val float64, err error) {
	ok, err := p.Scan(flag, &val)
	if err != nil {
		return
	}
	if !ok {
		val = def
	}
	p.register(flag, help, def)
	return
}

// Args returns all other un-scanned arguments of
// the passed arguments array.
//
// Example:
//   p := New([]string{"whats", "-n", "up"})
//   val, err := p.Bool("-n")
//   // val      = true
//   // err      = nil
//   // p.Args() = []string{"whats", "up"}
func (p *Parser) Args() []string {
	return p.args
}

// Help returns a compiled string listing all
// scanned flag specifications.
func (p *Parser) Help() string {
	lines := make([]string, len(p.flags))
	for i, f := range p.flags {
		lines[i] = fmt.Sprintf("  %s %s [default: %v]\n      %s",
			f.flag, typName(f.def), f.def, f.help)
	}
	return strings.Join(lines, "\n")
}

func (p *Parser) register(f string, help []string, def interface{}) {
	p.flags = append(p.flags, flag{
		flag: f,
		def:  def,
		help: optStr(help, ""),
	})
}
