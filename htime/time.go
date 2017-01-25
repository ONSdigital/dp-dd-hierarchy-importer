package htime

import (
	"fmt"
	"strconv"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
)

const (
	en = "en"
)

var monthNames = map[int]string{
	1:  "January",
	2:  "February",
	3:  "March",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

func CreateHierarchy(start int, end int) *sql.Hierarchy {
	h := sql.NewHierarchy()
	h.Names[en] = "time"
	h.HierarchyType = "time"
	h.ID = "time"

	addAreaTypes(h)

	for y := start; y <= end; y++ {
		year := sql.NewEntry()
		year.Code = strconv.Itoa(y)
		year.Names[en] = year.Code
		year.AreaType = "year"
		h.Entries[year.Code] = year
		for m := 1; m <= 12; m++ {
			month := sql.NewEntry()
			month.Code = fmt.Sprintf("%d.%02d", y, m)
			month.AreaType = "month"
			month.Names[en] = monthNames[m]
			month.DisplayOrder = m
			month.ParentCode = year.Code
			h.Entries[month.Code] = month
		}
		for q := 1; q <= 4; q++ {
			quarter := sql.NewEntry()
			quarter.Code = fmt.Sprintf("%d.Q%d", y, q)
			quarter.AreaType = "quarter"
			quarter.Names[en] = fmt.Sprintf("Q%d", q)
			quarter.ParentCode = year.Code
			h.Entries[quarter.Code] = quarter
		}
	}
	return &h
}

func addAreaTypes(h sql.Hierarchy) {
	h.AreaTypes["year"] = sql.LevelType{ID: "year", Name: "year", Level: 0}
	h.AreaTypes["month"] = sql.LevelType{ID: "month", Name: "month", Level: 1}
	h.AreaTypes["quarter"] = sql.LevelType{ID: "quarter", Name: "quarter", Level: 1}
}
