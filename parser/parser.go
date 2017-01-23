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

// OpenReader opens a reader on a local file or a url, as appropriate
func OpenReader(endpoint string) io.ReadCloser {
	if isURL(endpoint) {
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

// Parse parses the content from the Reader into the data object
func Parse(reader io.ReadCloser, data interface{}) {
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error reading body! %s\n", err)
		panic(err.Error())
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Character)"
			fmt.Printf("Invalid character at offset %v\n %s\n", t.Offset, jsn)
		case *json.UnmarshalTypeError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Type)"
			fmt.Printf("Invalid value at offset %v\n %s\n", t.Offset, jsn)
			fmt.Println("You might need to save the json file locally and validate the format - e.g. classifications with a single CodeList are known to be in an invalid format. See the readme for details.")
		default:
			fmt.Printf("Unable to unmarshal data. Error=%T, data=%s\n", t, string(body))
		}
		panic(err)
	}

}

func isURL(file string) bool {
	return strings.HasPrefix(file, "http")
}
