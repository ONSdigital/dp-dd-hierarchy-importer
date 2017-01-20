package sql

import (
	"fmt"
	"io"
	"strings"
	"sort"
)

const (
	hierarchySql = "insert into hierarchy (hierarchy_id, hierarchy_name) values (%s, %s);\n"
	areaSql = "insert into hierarchy_area_type (id, name, level) select %s, %s, %d where not exists (select id from hierarchy_area_type where id=%[1]s);\n"
	entrySql     = "insert into hierarchy_entry (hierarchy_id, entry_code, parent_code, name, area_type) values (%s, %s, %s, %s, %s);\n"
)


// write sql to the given writer.
// Please note that there is no error handling of failed writes
// This is a stand-alone command line tool, so errors can just be reported to the user
func WriteSql(writer io.Writer, hierarchy *Hierarchy) {

	if len(hierarchy.Id) == 0 {
		panic("Cannot write sql for a hierarchy without an id!")
	}
	io.WriteString(writer, fmt.Sprintf(hierarchySql, quote(hierarchy.Id), quote(hierarchy.Names["en"])))

	io.WriteString(writer, "\n")
	writeAreaTypes(writer, hierarchy.AreaTypes)

	io.WriteString(writer, "\n")
	writeEntries(writer, hierarchy.Entries, hierarchy.Id)

}

func ShouldWriteSql(hierarchy *Hierarchy) bool {
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
		io.WriteString(writer, fmt.Sprintf(areaSql, quote(area.Id), quote(area.Name), area.Level))
	}
}

func writeEntries(writer io.Writer, entries map[string]Entry, hierarchyId string) {
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
		writeEntry(writer, &entry, hierarchyId, entries, written)
	}
}

func writeEntry(writer io.Writer, entry *Entry, hierarchyId string, entries map[string]Entry, written map[string]bool) {
	if written[entry.Code] == true {
		return
	}
	if len(entry.ParentCode) > 0 {
		if parent, ok := entries[entry.ParentCode]; ok {
			writeEntry(writer, &parent, hierarchyId, entries, written)
		} else {
			fmt.Printf("!! Entry '%s' has an unknown parent '%s' - commenting out the insert\n", entry.Code, entry.ParentCode)
			io.WriteString(writer, "--")
		}
	}
	io.WriteString(writer, fmt.Sprintf(entrySql, quote(hierarchyId), quote(entry.Code), quote(entry.ParentCode), quote(entry.Names["en"]), quote(entry.AreaType)))
	written[entry.Code] = true
}

func quote(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "'" + strings.Replace(value, "'", "''", -1) + "'"
}
