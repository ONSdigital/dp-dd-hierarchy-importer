package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/geography"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/parser"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/structure"
)

const (
	DIMENSION_START_INDEX  = 3
	HIERARCHY_ID_OFFSET    = 0
	DIMENSION_VALUE_OFFSET = 2
	STRUCTURE_URL          = "http://web.ons.gov.uk/ons/api/data/classification/{id}.json?apikey={key}&context={context}"
	GEOGRAPHY_URL          = "http://web.ons.gov.uk/ons/api/data/hierarchies/hierarchy/{id}.json?apikey={key}&levels=0,1,2,3,4,5,6,7,8,9,10,11"
)

var CLASSIFICATION_CONTEXTS = [...]string{"Economic", "Census", "Socio-Economic", "Social"}

func ParseHierarchiesFromCSV(filename string) (map[string]map[string]bool, map[int]map[string]bool) {

	hierarchiesById := make(map[string]map[string]bool)
	hierarchiesByDimension := make(map[int]map[string]bool)

	r := parser.OpenReader(filename)
	defer r.Close()
	reader := csv.NewReader(r)

	reader.Read()

csvLoop:
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF reached, no more records to process", err.Error())
				break csvLoop
			} else {
				fmt.Println("Error occurred and cannot process anymore entry", err.Error())
				panic(err)
			}
		}
		dimensionIndex := 0
		for i := DIMENSION_START_INDEX; i < len(row); i = i + 3 {
			dimensionIndex++
			hierarchyId := strings.TrimSpace(row[i+HIERARCHY_ID_OFFSET])
			d, dimensionExists := hierarchiesByDimension[dimensionIndex]
			if !dimensionExists {
				d = make(map[string]bool)
				hierarchiesByDimension[dimensionIndex] = d
			}
			d[hierarchyId] = true
			if len(hierarchyId) > 0 {
				m, exists := hierarchiesById[hierarchyId]
				if !exists {
					m = make(map[string]bool)
					hierarchiesById[hierarchyId] = m
				}
				m[row[i+DIMENSION_VALUE_OFFSET]] = true
			}
		}
	}
	delete(hierarchiesById, "time")
	return hierarchiesById, hierarchiesByDimension
}

func FindAllHierarchies(filename string, apikey string) {
	hierarchiesInCsv, hierarchiesInDimensions := ParseHierarchiesFromCSV(filename)

	structureUrl := strings.Replace(STRUCTURE_URL, "{key}", apikey, -1)
	geographyUrl := strings.Replace(GEOGRAPHY_URL, "{key}", apikey, -1)

	found := make(map[string]bool)
	//missing := []string{}
	for hierarchyId := range hierarchiesInCsv {
		// try structural
		for _, context := range CLASSIFICATION_CONTEXTS {
			hierarchies := tryToLoadStructure(structureUrl, hierarchyId, context)
			if len(hierarchies) > 0 {
				fmt.Println(fmt.Sprintf("hierarchy %s in context %s contains %d hierarchies", hierarchyId, context, len(hierarchies)))
				for i, h := range hierarchies {
					containsAll := true
				inCSVLoop:
					for required := range hierarchiesInCsv[hierarchyId] {
						if _, exists := h.Entries[required]; !exists {
							containsAll = false
							break inCSVLoop
						}
					}
					if containsAll {
						fmt.Println(fmt.Sprintf("hierarchy %s[%d] in context %s contains all entries", hierarchyId, i, context))
						writeSQLForHierarchy(fmt.Sprintf("%s_%s_%d", hierarchyId, context, i), h)
						found[hierarchyId] = true
					}
				}
			}
		}
		if _, exists := found[hierarchyId]; !exists {
			// try geographic
			hierarchies := tryToLoadGeography(geographyUrl, hierarchyId)
			if len(hierarchies) > 0 {
				fmt.Println(fmt.Sprintf("Geographic hierarchy %s contains %d hierarchies", hierarchyId, len(hierarchies)))
				for i, h := range hierarchies {
					containsAll := true
				inCSVLoop2:
					for required := range hierarchiesInCsv[hierarchyId] {
						if _, exists := h.Entries[required]; !exists {
							containsAll = false
							break inCSVLoop2
						}
					}
					if containsAll {
						fmt.Println(fmt.Sprintf("Geographic hierarchy %s[%d] in contains all entries", hierarchyId, i))
						writeSQLForHierarchy(fmt.Sprintf("%s_%d", hierarchyId, i), h)
						found[hierarchyId] = true
					}
				}
			}
		}
	}
	fmt.Println()
	for idx, hierarchies := range hierarchiesInDimensions {
		fmt.Print(fmt.Sprintf("Dimension %d contains %d hierarchies: ", idx, len(hierarchies)))
		for hierarchyId := range hierarchies {
			fmt.Print(fmt.Sprintf("%s: (%d matched values) ", hierarchyId, len(hierarchiesInCsv[hierarchyId])))
		}
		fmt.Println()
	}
	fmt.Println()
	for hierarchyId := range hierarchiesInCsv {
		if _, exists := found[hierarchyId]; !exists {
			fmt.Println(fmt.Sprintf("Unable to find hierarchy matching all entries for %s: %v", hierarchyId, hierarchiesInCsv[hierarchyId]))
		}
	}
}

func tryToLoadStructure(baseUrl string, hierarchyId string, context string) []*sql.Hierarchy {
	var result []*sql.Hierarchy
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("Hierarchy %s is not in context %s: %s", hierarchyId, context, r)
			result = nil
		}
	}()
	url := strings.Replace(strings.Replace(baseUrl, "{id}", hierarchyId, -1), "{context}", context, -1)
	result = structure.LoadStructure(url)
	return result
}

func tryToLoadGeography(baseUrl string, hierarchyId string) []*sql.Hierarchy {
	var result []*sql.Hierarchy
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("Hierarchy %s is not in context %s: %s", hierarchyId, context, r)
			result = nil
		}
	}()
	url := strings.Replace(baseUrl, "{id}", hierarchyId, -1)
	result = append(result, geography.LoadGeography(url))
	return result
}

func writeSQLForHierarchy(filePrefix string, h *sql.Hierarchy) {
	depth := h.Depth()
	if depth < 3 {
		fmt.Printf("Hierarchy %s has a maximum depth of %d\n", h.ID, depth)
	}
	filename := filePrefix + ".sql"
	fmt.Printf("Creating sql file %s\n", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	sql.WriteSQL(file, h)
	fmt.Printf("Finished writing %s with %d entries\n", filename, len(h.Entries))
}
