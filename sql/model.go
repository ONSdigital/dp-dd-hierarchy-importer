package sql

import "github.com/satori/go.uuid"

// Hierarchy represents a complete hierarchy
type Hierarchy struct {
	ID            string
	Names         map[string]string
	Entries       map[string]Entry
	AreaTypes     map[string]LevelType
	HierarchyType string
}

// NewHierarchy returns a properly initialised Hierarchy
func NewHierarchy() Hierarchy {
	hierarchy := Hierarchy{}
	hierarchy.Entries = make(map[string]Entry)
	hierarchy.Names = make(map[string]string)
	hierarchy.AreaTypes = make(map[string]LevelType)
	return hierarchy
}

// LevelType is reference data for a type such as an Area - country, region etc
type LevelType struct {
	ID    string
	Name  string
	Level int
}

// Entry is a single entry in a hierarchy
type Entry struct {
	Code         string
	ParentCode   string
	AreaType     string
	Names        map[string]string
	DisplayOrder int
	UUID         string
}

// NewEntry returns a properly initialised Entry
func NewEntry() Entry {
	entry := Entry{}
	entry.Names = make(map[string]string)
	entry.UUID = uuid.NewV4().String()
	return entry
}
