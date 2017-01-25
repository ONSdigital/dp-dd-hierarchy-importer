package htime

import (
	"testing"

	"fmt"
	"strconv"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateHierarchy(t *testing.T) {

	Convey("Given a start and end year", t, func() {

		start, end := 1900, 2100

		Convey("When CreateHierarchy is invoked", func() {
			result := CreateHierarchy(start, end)

			So(result, ShouldNotBeNil)
			So(result.ID, ShouldEqual, "time")
			So(result.HierarchyType, ShouldEqual, "time")
			So(result.Names["en"], ShouldEqual, "time")

			So(len(result.AreaTypes), ShouldEqual, 3)
			So(result.AreaTypes["year"], ShouldNotBeNil)
			So(result.AreaTypes["month"], ShouldNotBeNil)
			So(result.AreaTypes["quarter"], ShouldNotBeNil)

			for y := start; y <= end; y++ {
				year := result.Entries[strconv.Itoa(y)]
				So(year, ShouldNotBeNil)
				So(len(year.ParentCode), ShouldEqual, 0)
				So(year.Names["en"], ShouldEqual, strconv.Itoa(y))
				So(year.AreaType, ShouldEqual, "year")
				for m := 1; m <= 12; m++ {
					month := result.Entries[fmt.Sprintf("%d.%02d", y, m)]
					So(month, ShouldNotBeNil)
					So(month.ParentCode, ShouldEqual, year.Code)
					So(month.Names["en"], ShouldEqual, monthNames[m])
					So(month.AreaType, ShouldEqual, "month")
				}
				for q := 1; q <= 4; q++ {
					quarter := result.Entries[fmt.Sprintf("%d.Q%d", y, q)]
					So(quarter, ShouldNotBeNil)
					So(quarter.ParentCode, ShouldEqual, year.Code)
					So(quarter.Names["en"], ShouldEqual, fmt.Sprintf("Q%d", q))
					So(quarter.AreaType, ShouldEqual, "quarter")
				}
			}
		})
	})
}
