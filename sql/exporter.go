package sql

import (
	"fmt"
	"io"
	"strings"
)

const (
	hierarchySql = "insert into hierarchy (hierarchy_id, hierarchy_name) values ('%s', %s);\n"
	entrySql     = "insert into hierarchy_entry (hierarchy_id, entry_code, parent_code, code_name, abbreviation, description, level) values ('%s', '%s', %s, %s, %s, %s, %d);\n"
)

// write sql to the given writer.
// Please note that there is no error handling of failed writes
// This is a stand-alone command line tool, so errors can just be reported to the user
func WriteInserts(writer io.Writer, hierarchy *Hierarchy) {

	if len(hierarchy.Id) == 0 {
		panic("Cannot write sql for a hierarchy without an id!")
	}
	io.WriteString(writer, fmt.Sprintf(hierarchySql, hierarchy.Id, quote(hierarchy.Names["en"])))
	io.WriteString(writer, "\n")

	written := make(map[string]bool)
	for _, entry := range hierarchy.Entries {
		writeEntry(writer, &entry, hierarchy.Id, hierarchy.Entries, written)
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
			fmt.Printf("!! Entry '%s' has an unknown parent '%s' - commenting out the insert")
			io.WriteString(writer, "--")
		}
	}
	io.WriteString(writer, fmt.Sprintf(entrySql, hierarchyId, entry.Code, quote(entry.ParentCode), quote(entry.Codename), quote(entry.Abbreviation), quote(entry.Names["en"]), entry.Level))
	written[entry.Code] = true
}

func quote(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "'" + strings.Replace(value, "'", "''", -1) + "'"
}
