package geography

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadDataToHierarchy(t *testing.T) {

	Convey("Given a reader containing json", t, func() {
		readcloser := ioutil.NopCloser(strings.NewReader(OriginalJson))

		originalData := GeographicData{}
		json.Unmarshal([]byte(OriginalJson), &originalData)

		Convey("When read into a hierarchy", func() {
			hierarchy := readHierarchy(readcloser)

			So(hierarchy, ShouldNotBeNil)
			So(hierarchy.Id, ShouldEqual, "2011STATH")
			So(hierarchy.Names["en"], ShouldEqual, "2011 Statistical Geography Hierarchy")
			So(hierarchy.Names["cy"], ShouldEqual, "Hierarchaeth Daearyddiaeth Ystadegol 2011")
			So(len(hierarchy.Entries), ShouldEqual, 5)

			for _, item := range originalData.ONS.GeographyList.Items.Item {
				entry := hierarchy.Entries[item.ItemCode]
				So(entry, ShouldNotBeNil)
				So(entry.Code, ShouldEqual, item.ItemCode)
				So(entry.ParentCode, ShouldEqual, item.ParentCode)
				So(entry.AreaType, ShouldEqual, item.AreaType.Abbreviation)
				areaType := hierarchy.AreaTypes[entry.AreaType]
				So(areaType, ShouldNotBeNil)
				So(areaType.Id, ShouldEqual, item.AreaType.Abbreviation)
				So(areaType.Name, ShouldEqual, item.AreaType.Codename)
				So(areaType.Level, ShouldEqual, int(item.AreaType.Level))
				for _, label := range item.Labels.Label {
					So(entry.Names[label.Lang], ShouldEqual, label.Label)
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
