package structure

type StructuralData struct {
	Structure *Structure
}

type Structure struct {
	CodeLists *CodeLists
}

type CodeLists struct {
	CodeList []*CodeList
}

type CodeList struct {
	Id    string  `json:"@id"`
	Names []*Name `json:"Name"`
	Codes []*Code `json:"Code"`
}

type Name struct {
	Lang string `json:"@xml.lang"`
	Name string `json:"$,omitempty"`
}

type Code struct {
	Value       string `json:"@value"`
	Urn         string `json:"@urn"`
	Parent      string `json:"@parentCode"`
	Description *Name
	Annotations *Annotations
}

type Annotations struct {
	Annotation []*Annotation
}

type Annotation struct {
	AnnotationType string
	AnnotationText *Name
}
