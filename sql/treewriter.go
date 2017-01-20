package sql

import (
	"fmt"
	"io"
	"sort"
)

type branch struct {
	code     string
	children map[string]branch
}

// WriteTree writes a tree representation of the given hierarchy to the writer
func WriteTree(writer io.Writer, hierarchy *Hierarchy, includeEmpty bool) {
	tree := assembleTree(hierarchy)
	keys := make([]string, len(tree))
	i := 0
	for k := range tree {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		b := tree[key]
		writeBranch(writer, b, "", hierarchy, includeEmpty)
	}
}

func writeBranch(writer io.Writer, b branch, indent string, hierarchy *Hierarchy, includeEmpty bool) {
	length := len(b.children)
	if includeEmpty || length > 0 {
		entry := hierarchy.Entries[b.code]
		fmt.Fprintf(writer, "%s%s - %s - %d children\n", indent, b.code, entry.Names["en"], length)
		keys := make([]string, length)
		i := 0
		for k := range b.children {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		for _, key := range keys {
			child := b.children[key]
			writeBranch(writer, child, indent+"    ", hierarchy, includeEmpty)
		}
	}
}

func assembleTree(hierarchy *Hierarchy) map[string]branch {
	topLevel := make(map[string]branch)
	all := make(map[string]branch)

	for _, entry := range hierarchy.Entries {
		b := ensureExistsInMap(entry, all)
		if parent, ok := hierarchy.Entries[entry.ParentCode]; ok {
			p := ensureExistsInMap(parent, all)
			p.children[b.code] = b
		}
		if len(entry.ParentCode) == 0 {
			topLevel[b.code] = b
		}
	}
	return topLevel
}

func ensureExistsInMap(entry Entry, all map[string]branch) branch {
	b, ok := all[entry.Code]
	if !ok {
		b = branch{code: entry.Code, children: make(map[string]branch)}
		all[entry.Code] = b
	}
	return b
}
