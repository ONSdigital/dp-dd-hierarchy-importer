package structure

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/parser"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
)

const nilErrorMessage = "Structure or CodeLists is null"

// LoadStructure reads the contents of an endpoint as a classification structure, returns a slice of hierarchies
func LoadStructure(endpoint string) []*sql.Hierarchy {
	reader := parser.OpenReader(endpoint)
	defer reader.Close()
	return readHierarchy(reader)
}

func readHierarchy(reader io.ReadCloser) []*sql.Hierarchy {
	data := readData(reader)
	return convertToHierarchy(data)
}

// Reads from the reader into a StructuralData object and converts the result into a Hierarchy object
func readData(reader io.ReadCloser) *StructuralData {
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading body!")
		panic(err.Error())
	}

	var data *StructuralData
	err = json.Unmarshal(body, &data)
	return data

}

func convertToHierarchy(data *StructuralData) []*sql.Hierarchy {

	if data == nil || data.Structure == nil || data.Structure.CodeLists == nil {
		panic(nilErrorMessage)
	}
	var hierarchies []*sql.Hierarchy

	for i, codeList := range data.Structure.CodeLists.CodeList {
		hierarchy := sql.NewHierarchy()
		hierarchies = append(hierarchies, &hierarchy)

		hierarchy.ID = codeList.ID + "_" + strconv.Itoa(i)
		for _, name := range codeList.Names {
			hierarchy.Names[name.Lang] = name.Name
		}
		for _, item := range codeList.Codes {
			entry := sql.NewEntry()
			entry.Code = item.Value
			entry.ParentCode = item.Parent
			entry.Names[item.Description.Lang] = item.Description.Name
			hierarchy.Entries[entry.Code] = entry
		}

	}
	return hierarchies
}
