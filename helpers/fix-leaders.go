package helpers

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func FixLeader(leader, farm, application string) {
	res, err := ERPRequest(
		"PUT", "/api/resource/Farm Application/"+application,
		strings.NewReader(fmt.Sprintf(
			"{\"engineer_name\": \"%s\", \"leading_engineers\": 1}",
			leader)))
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		f, err := os.Create("res.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(f, res.Body)
		panic(res.StatusCode)
	}
	res, err = ERPRequest("PUT", "/api/resource/Farm/"+farm, strings.NewReader("{\"leading_engineers\": 1}"))
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		f, err := os.Create("res.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(f, res.Body)
		panic(res.StatusCode)
	}
}
