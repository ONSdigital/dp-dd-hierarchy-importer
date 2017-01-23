package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"strconv"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/geography"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/structure"
)

var hierarchyType = flag.String("type", "", "'g' (geographical hierarchy) or 's' (structural hierarchy)")
var printTree = flag.String("tree", "", "If specified, 'b' will write a tree showing the hierarchy excluding leaf nodes, l will include leaf nodes")

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Something went wrong: %s\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	fmt.Println()
	checkCommandLineArgs()
	dir := getWorkingDir()

	hierarchies := loadHierarchies(*hierarchyType, flag.Arg(0))

	if len(hierarchies) == 0 {
		fmt.Println("No hierarchies found! Nothing to do")
	}
	var lastID string
	duplicates := false
	for i, h := range hierarchies {
		filename := filepath.Join(dir, h.ID)
		if h.ID == lastID {
			filename = filename + "_" + strconv.Itoa(i)
			duplicates = true
		}
		writeSQLForHierarchy(filename, h)
		writeTreeForHierarchy(filename, h)
		lastID = h.ID
		fmt.Println()
	}
	if duplicates {
		fmt.Println("!! Please note that there were multiple hierarchies with the same id. You should examine the files to decide which version to import - you cannot import both. You can use the -tree option to make comparison easier.")
	}
}

func checkCommandLineArgs() {
	flag.Parse()
	if len(flag.Args()) != 1 || (*hierarchyType != "g" && *hierarchyType != "s") {
		_, exe := filepath.Split(os.Args[0])
		fmt.Println("ONS hierarchy importer. Reads a json representation of a hierarchy or classification, and creates a set of sql insert statements to reconstruct a hierarchy in the db")
		fmt.Println("Please specify a type argument of 'g' (geographical hierarchy) or 's' (structural hierarchy/classification), and the location of the file to parse, e.g.")
		fmt.Println(exe + " -type=g 'http://web.ons.gov.uk/ons/api/data/hierarchies/hierarchy/2011WKWZH.json?apikey=XXXXX&levels=0,1,2'")
		fmt.Println("or")
		fmt.Println(exe + " -type=s 'http://web.ons.gov.uk/ons/api/data/classification/CL_0000641.json?apikey=XXXXX&context=Economic'")
		fmt.Println("or")
		fmt.Println(exe + " -type=g /tmp/localfile.json")
		fmt.Println("There is also a 'tree=b or -tree=l option, which will write a tree depiction of the hierarchy excluding leaves (b) or including them (l)")
		os.Exit(0)
	}
}

func loadHierarchies(t string, file string) []*sql.Hierarchy {
	var hierarchies []*sql.Hierarchy
	fmt.Printf("Importing hierarchies from %s\n", file)
	switch t {
	case "g":
		hierarchies = append(hierarchies, geography.LoadGeography(file))
	case "s":
		hierarchies = append(hierarchies, structure.LoadStructure(file)...)
	}
	return hierarchies
}

func writeSQLForHierarchy(filePrefix string, h *sql.Hierarchy) {
	depth := h.Depth()
	if depth < 3 {
		fmt.Printf("Hierarchy %s is only has a depth of %d - are you sure this qualifies as a hierarchy?\n", h.ID, depth)
	}
	filename := filePrefix + ".sql"
	fmt.Printf("Creating sql file %s\n", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	sql.WriteSQL(file, h)
	fmt.Printf("Finished writing %s with %d entries\n", filename, len(h.Entries))
}

func writeTreeForHierarchy(filePrefix string, h *sql.Hierarchy) {
	if len(*printTree) == 0 {
		return
	}
	filename := filePrefix + "_tree.txt"
	fmt.Printf("Creating tree %s\n", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	includeEmpty := *printTree == "l"
	sql.WriteTree(file, h, includeEmpty)
}

func getWorkingDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
