package structure

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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
	CodeList CodeListHolder `json:"CodeList"`
}

// CodeListHolder is an artefact of the fact that the data api can return either a single object or an array
type CodeListHolder struct {
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
	Annotation AnnotationHolder `json:"Annotation"`
}

// AnnotationHolder is an artefact of the fact that the data api can return a single object or an array for the 'Annotation' property
type AnnotationHolder struct {
	Annotation []*Annotation
}

// Annotation is information such as DisplayOrder
type Annotation struct {
	AnnotationType string
	AnnotationText *Name
}

// UnmarshalJSON can unmarshal a single CodeList or an array
func (holder *CodeListHolder) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewBuffer(data))
	t, err := d.Token()
	if err != nil {
		return err
	}
	if t == json.Delim('[') {
		if err := json.Unmarshal(data, &holder.CodeList); err != nil {
			fmt.Printf("Unable to Unmarshal []*CodeList from: %s", string(data))
			return err
		}
	} else {
		var c *CodeList
		if err := json.Unmarshal(data, &c); err != nil {
			fmt.Printf("Unable to Unmarshal *CodeList from: %s", string(data))
			return err
		}
		holder.CodeList = append(holder.CodeList, c)
	}
	return nil
}

// MarshalJSON removes the CodeListHolder from the json
func (holder *CodeListHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(holder.CodeList)
}

// GetCodeLists is a convenience method to avoid writing s.CodeLists.CodeList.CodeList
func (s *Structure) GetCodeLists() []*CodeList {
	return s.CodeLists.CodeList.CodeList
}

// UnmarshalJSON can unmarshal a single Annotation or an array
func (holder *AnnotationHolder) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewBuffer(data))
	t, err := d.Token()
	if err != nil {
		return err
	}
	if t == json.Delim('[') {
		if err := json.Unmarshal(data, &holder.Annotation); err != nil {
			fmt.Printf("Unable to Unmarshal []*Annotation from: %s", string(data))
			return err
		}
	} else {
		var a *Annotation
		if err := json.Unmarshal(data, &a); err != nil {
			fmt.Printf("Unable to Unmarshal *Annotation from: %s", string(data))
			return err
		}
		holder.Annotation = append(holder.Annotation, a)
	}
	return nil
}

// MarshalJSON removes the AnnotationHolder from the resulting json
func (holder *AnnotationHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(holder.Annotation)
}

// GetAnnotations is a convenience method to avoid having to write a.Annotations.Annotation.Annotation
func (c *Code) GetAnnotations() []*Annotation {
	return c.Annotations.Annotation.Annotation
}
