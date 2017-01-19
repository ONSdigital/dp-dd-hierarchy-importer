package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func OpenReader(endpoint string) io.ReadCloser {
	if isUrl(endpoint) {
		response, err := http.Get(endpoint)
		if err != nil {
			fmt.Println("Error calling endpoint!")
			panic(err)
		}
		return response.Body
	}

	file, err := os.Open(endpoint)
	if err != nil {
		fmt.Printf("Error opening file '%s': %s\n", endpoint, err)
		panic(err)
	}
	return file

}

func Parse(reader io.ReadCloser, data interface{}) {
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error reading body! %s\n", err)
		panic(err.Error())
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		fmt.Printf("Unable to marshal data into %T: %s\n", data, err)
		fmt.Printf("Tried to read:\n%s\n", string(body))
		panic(err.Error())
	}

}

func isUrl(file string) bool {
	return strings.HasPrefix(file, "http")
}