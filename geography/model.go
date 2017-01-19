package geography

type GeographicData struct {
	ONS *ONS `json:"ons"`
}

type ONS struct {
	GeographyList *GeographyList `json:"geographyList"`
}

type GeographyList struct {
	Geography *Geography `json:"geography"`
	Items     *Items     `json:"items"`
}

type Geography struct {
	Id    string `json:"id"`
	Names *Names `json:"names"`
}

type Items struct {
	Item []*Item `json:"item"`
}

type Item struct {
	Labels            *Labels   `json:"labels"`
	ItemCode          string    `json:"itemCode"`
	ParentCode        string    `json:"parentCode,omitempty"`
	AreaType          *AreaType `json:"areaType"`
	SubthresholdAreas string    `json:"subthresholdAreas"`
}

type AreaType struct {
	Abbreviation string  `json:"abbreviation"`
	Codename     string  `json:"codename"`
	Level        float64 `json:"level"`
}

type Labels struct {
	Label []*Label `json:"label"`
}

type Label struct {
	Lang  string `json:"@xml.lang"`
	Label string `json:"$"`
}

type Names struct {
	Name []*Name `json:"name"`
}

type Name struct {
	Lang string `json:"@xml.lang"`
	Name string `json:"$"`
}
