package sql

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"strconv"

	. "github.com/smartystreets/goconvey/convey"
)

var sampleHierarchy = Hierarchy{
	ID:    "myHierarchyId",
	Names: map[string]string{"en": "my hierarchy name"},
	Entries: map[string]Entry{
		"c":      {Code: "c", AreaType: "C", ParentCode: "b", Names: map[string]string{"en": "level c name", "cy": "welsh name"}},
		"b":      {Code: "b", AreaType: "B", ParentCode: "a", Names: map[string]string{"en": "level b name", "cy": "welsh name"}},
		"a":      {Code: "a", AreaType: "A", ParentCode: "", Names: map[string]string{"en": "level a name", "cy": "welsh name"}},
		"orphan": {Code: "orphan", AreaType: "C", ParentCode: "x", Names: map[string]string{"en": "orphan 'name'", "cy": "welsh name"}},
	},
	AreaTypes: map[string]LevelType{
		"A": {ID: "A", Level: 1, Name: "Level 1 -A"},
		"B": {ID: "B", Level: 2, Name: "Level 2 -B"},
		"C": {ID: "C", Level: 3, Name: "Level 3 -C"},
	},
}

func TestWriteInserts(t *testing.T) {

	Convey("When WriteSQL is invoked with a hierarchy", t, func() {
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		WriteSQL(writer, &sampleHierarchy)
		writer.Flush()
		var lines []string
		scanner := bufio.NewScanner(bufio.NewReader(&buffer))
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		Convey("Then insert statements should appear in the correct order: hierarchy, area types; entries", func() {
			So(len(lines), ShouldEqual, (2 + len(sampleHierarchy.AreaTypes) + 1 + len(sampleHierarchy.Entries)))

			idx := 0
			// insert hierarchy
			line := lines[idx]
			idx++
			So(line, ShouldStartWith, "insert into hierarchy ")
			So(line, ShouldContainSubstring, sampleHierarchy.Names["en"])

			// blank line
			line = lines[idx]
			idx++
			So(len(strings.Trim(line, " ")), ShouldEqual, 0)

			// area types
			for _, area := range []string{"A", "B", "C"} {
				line := lines[idx]
				idx++
				So(line, ShouldStartWith, "insert into hierarchy_level_type ")
				So(line, ShouldContainSubstring, "'"+sampleHierarchy.AreaTypes[area].Name+"'")
				So(line, ShouldContainSubstring, strconv.Itoa(sampleHierarchy.AreaTypes[area].Level))
			}

			// blank line
			line = lines[idx]
			idx++
			So(len(strings.Trim(line, " ")), ShouldEqual, 0)

			// entries
			for _, e := range []string{"a", "b", "c"} {
				line := lines[idx]
				idx++
				entry := sampleHierarchy.Entries[e]
				So(line, ShouldStartWith, "insert into hierarchy_entry ")
				So(line, ShouldContainSubstring, "'"+entry.Code+"'")
				if len(entry.ParentCode) == 0 {
					So(line, ShouldContainSubstring, "null")
				} else {
					So(line, ShouldContainSubstring, "'"+entry.ParentCode+"'")
				}
			}

			orphan := lines[idx]
			So(orphan, ShouldStartWith, "insert into hierarchy_entry ")
			So(orphan, ShouldContainSubstring, "'orphan'")
			So(orphan, ShouldContainSubstring, "orphan ''name'''")
		})
	})
}

var flat = Hierarchy{
	ID:    "flat",
	Names: map[string]string{"en": "my hierarchy name"},
	Entries: map[string]Entry{
		"c": {Code: "c", ParentCode: "a", Names: map[string]string{"en": "level c name", "cy": "welsh name"}},
		"b": {Code: "b", ParentCode: "a", Names: map[string]string{"en": "level b name", "cy": "welsh name"}},
		"a": {Code: "a", ParentCode: "", Names: map[string]string{"en": "level a name", "cy": "welsh name"}},
	},
	AreaTypes: map[string]LevelType{},
}

func TestHierarchyDepth(t *testing.T) {

	Convey("When WriteSQL is invoked with a hierarchy", t, func() {
		So(sampleHierarchy.Depth(), ShouldEqual, 3)
		So(flat.Depth(), ShouldEqual, 2)
		So(Hierarchy{}.Depth(), ShouldEqual, 0)
	})
}
