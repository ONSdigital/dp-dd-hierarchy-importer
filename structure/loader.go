package structure

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/parser"
	. "github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
)

func LoadStructure(endpoint string) []*Hierarchy {
	reader := parser.OpenReader(endpoint)
	defer reader.Close()
	return readHierarchy(reader)
}

func readHierarchy(reader io.ReadCloser) []*Hierarchy {
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

func convertToHierarchy(data *StructuralData) []*Hierarchy {

	var hierarchies []*Hierarchy

	for i, codeList := range data.Structure.CodeLists.CodeList {
		hierarchy := NewHierarchy()
		hierarchies = append(hierarchies, &hierarchy)

		hierarchy.Id = codeList.Id + "_" + strconv.Itoa(i)
		for _, name := range codeList.Names {
			hierarchy.Names[name.Lang] = name.Name
		}
		for _, item := range codeList.Codes {
			entry := NewEntry()
			entry.Code = item.Value
			entry.ParentCode = item.Parent
			entry.Codename = item.Description.Name
			entry.Names[item.Description.Lang] = item.Description.Name
			hierarchy.Entries[entry.Code] = entry
		}

		for _, entry := range hierarchy.Entries {
			entry.Level = countLevel(entry, hierarchy.Entries)
		}
	}
	return hierarchies
}

func countLevel(entry Entry, entries map[string]Entry) int {
	if len(entry.ParentCode) == 0 {
		return 0
	}
	if parent, ok := entries[entry.ParentCode]; ok {
		return countLevel(parent, entries) + 1
	}
	return -1
}
