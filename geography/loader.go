package geography

import (
	"io"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/parser"
	. "github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
)

func LoadGeography(endpoint string) *Hierarchy {
	reader := parser.OpenReader(endpoint)
	defer reader.Close()
	return readHierarchy(reader)
}

func readHierarchy(reader io.ReadCloser) *Hierarchy {
	var data GeographicData
	parser.Parse(reader, &data)
	return convertToHierarchy(&data)
}

func convertToHierarchy(data *GeographicData) *Hierarchy {

	hierarchy := NewHierarchy()
	geog := data.ONS.GeographyList

	hierarchy.Id = geog.Geography.Id
	for _, name := range geog.Geography.Names.Name {
		hierarchy.Names[name.Lang] = name.Name
	}

	for _, item := range geog.Items.Item {
		code := item.ItemCode
		entry := NewEntry()
		entry.Code = code
		if item.AreaType != nil {
			entry.Level = int(item.AreaType.Level)
			entry.Abbreviation = item.AreaType.Abbreviation
			entry.Codename = item.AreaType.Codename
		}
		entry.ParentCode = item.ParentCode
		for _, label := range item.Labels.Label {
			entry.Names[label.Lang] = label.Label
		}
		hierarchy.Entries[code] = entry
	}
	return &hierarchy
}
