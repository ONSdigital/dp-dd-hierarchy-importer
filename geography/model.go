package geography

// GeographicData is a Top-level container, not represented in json
type GeographicData struct {
	ONS *ONS `json:"ons"`
}

// ONS is the top-level json object
type ONS struct {
	GeographyList *List `json:"geographyList"`
}

// List contains details of the Geography and a list of Items
type List struct {
	Geography *Geography `json:"geography"`
	Items     *Items     `json:"items"`
}

// Geography has ID and names of the geography
type Geography struct {
	ID    string `json:"id"`
	Names *Names `json:"names"`
}

// Items is A list of Items
type Items struct {
	Item []*Item `json:"item"`
}

// Item is a single entry in the geography hierarchy
type Item struct {
	Labels            *Labels   `json:"labels"`
	ItemCode          string    `json:"itemCode"`
	ParentCode        string    `json:"parentCode,omitempty"`
	AreaType          *AreaType `json:"areaType"`
	SubthresholdAreas string    `json:"subthresholdAreas"`
}

// AreaType is reference data for area types - country, region, etc
type AreaType struct {
	Abbreviation string  `json:"abbreviation"`
	Codename     string  `json:"codename"`
	Level        float64 `json:"level"`
}

// Labels Multi-lingual labels
type Labels struct {
	Label []*Label `json:"label"`
}

// Label in a single language
type Label struct {
	Lang  string `json:"@xml.lang"`
	Label string `json:"$"`
}

// Names multi-lingual names
type Names struct {
	Name []*Name `json:"name"`
}

// Name in a single language
type Name struct {
	Lang string `json:"@xml.lang"`
	Name string `json:"$"`
}
