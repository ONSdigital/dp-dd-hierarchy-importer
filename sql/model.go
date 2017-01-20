package sql

type Hierarchy struct {
	Id      string
	Names   map[string]string
	Entries map[string]Entry
	AreaTypes map[string]LevelType
}

func NewHierarchy() Hierarchy {
	hierarchy := Hierarchy{}
	hierarchy.Entries = make(map[string]Entry)
	hierarchy.Names = make(map[string]string)
	hierarchy.AreaTypes = make(map[string]LevelType)
	return hierarchy
}

type LevelType struct {
	Id string
	Name string
	Level int
}

type Entry struct {
	Code         string
	ParentCode   string
	AreaType string
	Names        map[string]string
}

func NewEntry() Entry {
	entry := Entry{}
	entry.Names = make(map[string]string)
	return entry
}
