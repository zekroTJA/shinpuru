package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/zekroTJA/shinpuru/internal/models"
)

const healthcheckEndpoint = "/api/v1/healthcheck"

var (
	fEndpoint = flag.String("addr", "http://localhost", "shinpuru instance address")
)

func main() {
	flag.Parse()

	addr := *fEndpoint + healthcheckEndpoint
	resp, err := http.Get(addr)
	checkErr(err)

	if resp.StatusCode != 200 {
		exit(2, "error: response status was %d\n", resp.StatusCode)
	}

	var status models.HealthcheckResponse
	err = json.NewDecoder(resp.Body).Decode(&status)
	checkErr(err)

	printStatus(status)

	if !status.AllOk {
		exit(3, "\nsystem state faulty\n")
	}
}

func checkErr(err error) {
	if err == nil {
		return
	}
	exit(1, "error: %s\n", err.Error())
}

func exit(code int, msg string, args ...any) {
	fmt.Printf(msg, args...)
	os.Exit(code)
}

func printStatusLine(name string, s models.HealthcheckStatus) {
	if s.Ok {
		fmt.Printf("%-14s ok\n", name)
	} else {
		fmt.Printf("%-14s faulty (%s)\n", name, s.Message)
	}
}

func printStatus(status models.HealthcheckResponse) {
	v := reflect.ValueOf(status)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if s, ok := f.Interface().(models.HealthcheckStatus); ok {
			printStatusLine(v.Type().Field(i).Name, s)
		}
	}
}
