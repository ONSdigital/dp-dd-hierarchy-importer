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
		readcloser := ioutil.NopCloser(strings.NewReader(OriginalJSON))

		originalData := GeographicData{}
		json.Unmarshal([]byte(OriginalJSON), &originalData)

		Convey("When read into a hierarchy", func() {
			hierarchy := readHierarchy(readcloser)

			So(hierarchy, ShouldNotBeNil)
			So(hierarchy.ID, ShouldEqual, "2011STATH")
			So(hierarchy.HierarchyType, ShouldEqual, "geography")
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
				So(areaType.ID, ShouldEqual, item.AreaType.Abbreviation)
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

func TestReadInvalidDataToHierarchy(t *testing.T) {

	Convey("Given a reader containing invalid json", t, func() {
		readcloser := ioutil.NopCloser(strings.NewReader(`{"ons":{"base":{"@href":"http://web.ons.gov.uk/ons/api/data/"},"node":{"urls":{"url":[{"@representation":"xml","href":"hierarchies/hierarchy/2011WKWZH.xml?apikey=Y6Xs59zXU0&levels=0,1,2"},]},"description":"","name":"Geography Classifications"},"linkedNodes":{"linkedNode":{`))

		Convey("When read into a hierarchy", func() {
			defer func() {
				r := recover()
				So(r, ShouldNotBeNil)
			}()
			readHierarchy(readcloser)

		})
	})
}
