package structure

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadDataToHierarchy(t *testing.T) {

	Convey("Given a reader containing json", t, func() {
		readcloser := ioutil.NopCloser(strings.NewReader(OriginalJSON))

		originalData := StructuralData{}

		json.Unmarshal([]byte(OriginalJSON), &originalData)

		Convey("When read into a hierarchy", func() {
			hierarchies := readHierarchy(readcloser)

			So(hierarchies, ShouldNotBeNil)
			So(len(hierarchies), ShouldEqual, 2)

			hierarchyMap := make(map[string]sql.Hierarchy)
			for _, h := range hierarchies {
				hierarchyMap[h.ID] = *h
			}

			for i, item := range originalData.Structure.CodeLists.CodeList {
				id := item.ID + "_" + strconv.Itoa(i)
				hierarchy := hierarchyMap[id]
				So(hierarchy, ShouldNotBeNil)
				So(hierarchy.ID, ShouldEqual, id)
				So(hierarchy.Names["en"], ShouldEqual, item.Names[0].Name)
				So(len(hierarchy.Entries), ShouldEqual, len(item.Codes))
				for _, item := range item.Codes {
					entry := hierarchy.Entries[item.Value]
					So(entry, ShouldNotBeNil)
					So(entry.Code, ShouldEqual, item.Value)
					So(entry.ParentCode, ShouldEqual, item.Parent)
					So(entry.AreaType, ShouldEqual, "")
					So(entry.Names[item.Description.Lang], ShouldEqual, item.Description.Name)
				}
			}

		})
	})
}

func TestReadEmptyDataToHierarchy(t *testing.T) {

	Convey("Given a reader containing nothing", t, func() {
		readcloser := ioutil.NopCloser(strings.NewReader("{}"))

		Convey("When read into a hierarchy", func() {
			defer func() {
				r := recover()
				So(r, ShouldNotBeNil)
				So(r, ShouldEqual, nilErrorMessage)
			}()
			readHierarchy(readcloser)

		})
	})
}
