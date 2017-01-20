package geography

import (
	"io"

	"fmt"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/parser"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
)

const nilErrorMessage = "ons or ons.geographyList is nil"

// LoadGeography Reads the content of the endpoint as a geographyList, returns a hierarchy
func LoadGeography(endpoint string) *sql.Hierarchy {
	reader := parser.OpenReader(endpoint)
	defer reader.Close()
	return readHierarchy(reader)
}

func readHierarchy(reader io.ReadCloser) *sql.Hierarchy {
	var data GeographicData
	parser.Parse(reader, &data)
	return convertToHierarchy(&data)
}

func convertToHierarchy(data *GeographicData) *sql.Hierarchy {

	if data == nil || data.ONS == nil || data.ONS.GeographyList == nil {
		panic(nilErrorMessage)
	}

	hierarchy := sql.NewHierarchy()
	geog := data.ONS.GeographyList

	hierarchy.ID = geog.Geography.ID
	for _, name := range geog.Geography.Names.Name {
		hierarchy.Names[name.Lang] = name.Name
	}

	for _, item := range geog.Items.Item {
		code := item.ItemCode
		entry := sql.NewEntry()
		entry.Code = code
		if item.AreaType != nil {
			key := item.AreaType.Abbreviation
			if area, exists := hierarchy.AreaTypes[key]; exists {
				if area.Name != item.AreaType.Codename {
					fmt.Println("WARNING: AreaType %s is defined multiple times with different names - '%s', '$s'", key, area.Name, item.AreaType.Codename)
				}
			} else {
				hierarchy.AreaTypes[key] = sql.LevelType{ID: key, Name: item.AreaType.Codename, Level: int(item.AreaType.Level)}
			}
			entry.AreaType = key
		}
		entry.ParentCode = item.ParentCode
		for _, label := range item.Labels.Label {
			entry.Names[label.Lang] = label.Label
		}
		hierarchy.Entries[code] = entry
	}
	return &hierarchy
}
