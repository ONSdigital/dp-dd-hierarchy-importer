package structure

// StructuralData is a top-level container, not represented in json
type StructuralData struct {
	Structure *Structure
}

// Structure is the top level json object
type Structure struct {
	CodeLists *CodeLists
}

// CodeLists contains a slice of CodeList object - an artefact of converting xml to json
type CodeLists struct {
	CodeList []*CodeList
}

// CodeList contains a hierarchy of codes.
type CodeList struct {
	ID    string  `json:"@id"`
	Names []*Name `json:"Name"`
	Codes []*Code `json:"Code"`
}

// Name is a name in a single language
type Name struct {
	Lang string `json:"@xml.lang"`
	Name string `json:"$,omitempty"`
}

// Code is a single entry in the hierarchy
type Code struct {
	Value       string `json:"@value"`
	Urn         string `json:"@urn"`
	Parent      string `json:"@parentCode"`
	Description *Name
	Annotations *Annotations
}

// Annotations contains a slice of Annotations
type Annotations struct {
	Annotation []*Annotation
}

// Annotation is information such as DisplayOrder
type Annotation struct {
	AnnotationType string
	AnnotationText *Name
}
