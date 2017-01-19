package sql

type Hierarchy struct {
	Id      string
	Names   map[string]string
	Entries map[string]Entry
}

func NewHierarchy() Hierarchy {
	hierarchy := Hierarchy{}
	hierarchy.Entries = make(map[string]Entry)
	hierarchy.Names = make(map[string]string)
	return hierarchy
}

type Entry struct {
	Code         string
	ParentCode   string
	Codename     string
	Abbreviation string
	Level        int
	Names        map[string]string
}

func NewEntry() Entry {
	entry := Entry{}
	entry.Names = make(map[string]string)
	return entry
}
