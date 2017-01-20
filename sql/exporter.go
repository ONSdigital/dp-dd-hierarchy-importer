package sql

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	hierarchySQL = "insert into hierarchy (hierarchy_id, hierarchy_name) values (%s, %s);\n"
	areaSQL      = "insert into hierarchy_area_type (id, name, level) values (%s, %s, %d) on conflict do nothing;\n"
	entrySQL     = "insert into hierarchy_entry (hierarchy_id, entry_code, parent_code, name, area_type) values (%s, %s, %s, %s, %s);\n"
)

// WriteSQL writes sql insert statements to the given writer to create the given hierarchy in a db.
// Please note that there is no error handling of failed writes
// This is a stand-alone command line tool, so errors can just be reported to the user
func WriteSQL(writer io.Writer, hierarchy *Hierarchy) {

	if len(hierarchy.ID) == 0 {
		panic("Cannot write sql for a hierarchy without an id!")
	}
	io.WriteString(writer, fmt.Sprintf(hierarchySQL, quote(hierarchy.ID), quote(hierarchy.Names["en"])))

	io.WriteString(writer, "\n")
	writeAreaTypes(writer, hierarchy.AreaTypes)

	io.WriteString(writer, "\n")
	writeEntries(writer, hierarchy.Entries, hierarchy.ID)

}

// ShouldWriteSQL returns true if the hierarchy has sufficient depth (at least one entry has grandchildren)
func ShouldWriteSQL(hierarchy *Hierarchy) bool {
	for _, entry := range hierarchy.Entries {
		if countLevel(entry, hierarchy.Entries) > 1 {
			return true
		}
	}
	return false
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
	if len(entry.ParentCode) > 0 {
		if parent, ok := entries[entry.ParentCode]; ok {
			writeEntry(writer, &parent, hierarchyID, entries, written)
		} else {
			fmt.Printf("!! Entry '%s' has an unknown parent '%s' - commenting out the insert\n", entry.Code, entry.ParentCode)
			io.WriteString(writer, "--")
		}
	}
	io.WriteString(writer, fmt.Sprintf(entrySQL, quote(hierarchyID), quote(entry.Code), quote(entry.ParentCode), quote(entry.Names["en"]), quote(entry.AreaType)))
	written[entry.Code] = true
}

func quote(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "'" + strings.Replace(value, "'", "''", -1) + "'"
}
