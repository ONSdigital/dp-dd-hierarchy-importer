package sql

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	hierarchySQL = "insert into hierarchy (id, name, type) values (%s, %s, %s);\n"
	areaSQL      = "insert into hierarchy_level_type (id, name, level) values (%s, %s, %d) on conflict do nothing;\n"
	entrySQL     = "insert into hierarchy_entry (id, hierarchy_id, code, parent, name, hierarchy_level_type_id, display_order) values (%s, %s, %s, %s, %s, %s, %d);\n"
)

// WriteSQL writes sql insert statements to the given writer to create the given hierarchy in a db.
// Please note that there is no error handling of failed writes
// This is a stand-alone command line tool, so errors can just be reported to the user
func WriteSQL(writer io.Writer, hierarchy *Hierarchy) {

	if len(hierarchy.ID) == 0 {
		panic("Cannot write sql for a hierarchy without an id!")
	}
	io.WriteString(writer, fmt.Sprintf(hierarchySQL, quote(hierarchy.ID), quote(hierarchy.Names["en"]), quote(hierarchy.HierarchyType)))

	io.WriteString(writer, "\n")
	writeAreaTypes(writer, hierarchy.AreaTypes)

	io.WriteString(writer, "\n")
	writeEntries(writer, hierarchy.Entries, hierarchy.ID)

}

// Depth returns the maximum depth of the hierarchy
func (hierarchy Hierarchy) Depth() int {
	maxDepth := 0
	for _, entry := range hierarchy.Entries {
		depth := countLevel(entry, hierarchy.Entries) + 1
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
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

func writeAreaTypes(writer io.Writer, areas map[string]LevelType) {
	keys := make([]string, len(areas))
	i := 0
	for k := range areas {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		area := areas[key]
		io.WriteString(writer, fmt.Sprintf(areaSQL, quote(area.ID), quote(area.Name), area.Level))
	}
}

func writeEntries(writer io.Writer, entries map[string]Entry, hierarchyID string) {
	keys := make([]string, len(entries))
	i := 0
	for k := range entries {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	written := make(map[string]bool)
	for _, key := range keys {
		entry := entries[key]
		writeEntry(writer, &entry, hierarchyID, entries, written)
	}
}

func writeEntry(writer io.Writer, entry *Entry, hierarchyID string, entries map[string]Entry, written map[string]bool) {
	if written[entry.Code] == true {
		return
	}
	var parentId string
	if len(entry.ParentCode) > 0 {
		if parent, ok := entries[entry.ParentCode]; ok {
			writeEntry(writer, &parent, hierarchyID, entries, written)
			parentId = parent.UUID
		} else {
			fmt.Printf("!! Entry '%s' has an unknown parent '%s' - the hierarchy is incomplete and some insert statements will fail\n", entry.Code, entry.ParentCode)
		}
	}
	io.WriteString(writer, fmt.Sprintf(entrySQL, quote(entry.UUID), quote(hierarchyID), quote(entry.Code), quote(parentId), quote(entry.Names["en"]), quote(entry.AreaType), entry.DisplayOrder))
	written[entry.Code] = true
}

func quote(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "'" + strings.Replace(value, "'", "''", -1) + "'"
}
