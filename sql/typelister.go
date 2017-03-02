package sql

import (
	"fmt"
	"io"
	"sort"
)

const (
	noAreaType = "[No Area Type]"
)

// WriteLists writes flat lists of all entries of each level type
func WriteLists(writer io.Writer, hierarchy *Hierarchy) {

	list := createFlatList(hierarchy)
	fmt.Fprintf(writer, "%s contains %d AreaTypes:\n", hierarchy.ID, len(hierarchy.AreaTypes))
	fmt.Fprintf(writer, "\n")

	areas := hierarchy.AreaTypes
	keys := make([]string, len(areas))
	i := 0
	for k := range areas {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		area := areas[key]
		fmt.Fprintf(writer, "%s - %s, level %d: %d entries\n", area.ID, area.Name, area.Level, len(list[area.ID]))
	}
	fmt.Fprintf(writer, "%s: %d entries\n", noAreaType, len(list[noAreaType]))
	fmt.Fprintf(writer, "\n\n")

	for _, key := range keys {
		entries := list[key]
		fmt.Fprintf(writer, "\n******************************************\n%s: %d entries\n", key, len(entries))
		for _, e := range entries {
			fmt.Fprintf(writer, "%s - %s\n", e, hierarchy.Entries[e].Names["en"])
		}
	}
	fmt.Fprintf(writer, "\n******************************************\n%s: %d entries\n", noAreaType, len(noAreaType))
	for _, e := range list[noAreaType] {
		fmt.Fprintf(writer, "%s - %s\n", e, hierarchy.Entries[e].Names["en"])
	}

}

func createFlatList(hierarchy *Hierarchy) map[string][]string {
	list := make(map[string][]string)

	for _, e := range hierarchy.Entries {
		areaType := e.AreaType
		if len(areaType) == 0 {
			areaType = noAreaType
		}
		list[areaType] = append(list[areaType], e.Code)
	}

	return list
}
