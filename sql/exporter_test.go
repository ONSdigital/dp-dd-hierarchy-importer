package sql

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"io/ioutil"

	. "github.com/smartystreets/goconvey/convey"
)

var sampleHierarchy = Hierarchy{
	Id:    "myHierarchyId",
	Names: map[string]string{"en": "my hierarchy name"},
	Entries: map[string]Entry{
		"c":      {Code: "c", Codename: "level c", Abbreviation: "abc", Level: 2, ParentCode: "b", Names: map[string]string{"en": "level c name", "cy": "welsh name"}},
		"b":      {Code: "b", Codename: "level b", Abbreviation: "ab", Level: 1, ParentCode: "a", Names: map[string]string{"en": "level b name", "cy": "welsh name"}},
		"a":      {Code: "a", Codename: "level a", Abbreviation: "a", Level: 0, ParentCode: "", Names: map[string]string{"en": "level a name", "cy": "welsh name"}},
		"orphan": {Code: "orphan", Codename: "orphan level", Abbreviation: "orphan", Level: 2, ParentCode: "x", Names: map[string]string{"en": "orphan 'name'", "cy": "welsh name"}},
	},
}

func TestWriteInserts(t *testing.T) {

	Convey("When WriteInserts is invoked with a hierarchy", t, func() {
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)
		WriteInserts(writer, &sampleHierarchy)
		writer.Flush()
		ioutil.WriteFile("/home/matt/temp/output.sql", buffer.Bytes(), 0644)
		var lines []string
		scanner := bufio.NewScanner(bufio.NewReader(&buffer))
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		Convey("Then insert statements should appear in the correct order, with hierarchy separated from the rest by a blank line", func() {
			So(len(lines), ShouldEqual, 6)

			// insert hierarchy
			line := lines[0]
			So(line, ShouldStartWith, "insert into hierarchy ")
			So(line, ShouldContainSubstring, sampleHierarchy.Names["en"])

			// blank line
			line = lines[1]
			So(len(strings.Trim(line, " ")), ShouldEqual, 0)

			// entry inserts should be a,b,c, with a commented-out orphan at any point
			var entries []string
			var orphan string
			for _, s := range lines[2:] {
				if strings.HasPrefix(s, "--") {
					orphan = s
				} else {
					entries = append(entries, s)
				}
			}
			line = entries[0]
			So(line, ShouldStartWith, "insert into hierarchy_entry ")
			So(line, ShouldContainSubstring, "'a', null")
			line = entries[1]
			So(line, ShouldStartWith, "insert into hierarchy_entry ")
			So(line, ShouldContainSubstring, "'b'")
			line = entries[2]
			So(line, ShouldStartWith, "insert into hierarchy_entry ")
			So(line, ShouldContainSubstring, "'c'")

			So(orphan, ShouldStartWith, "--insert into hierarchy_entry ")
			So(orphan, ShouldContainSubstring, "'orphan'")
			So(orphan, ShouldContainSubstring, "orphan ''name'''")
		})
	})
}
