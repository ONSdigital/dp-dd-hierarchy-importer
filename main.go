package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"strconv"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/csvparser"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/geography"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/htime"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/structure"
)

var hierarchyType = flag.String("type", "", "'g' (geographical hierarchy), 's' (structural hierarchy), or 't' (time)")
var start = flag.Int("start", 1900, "If type=t, the start year. Default 1900")
var end = flag.Int("end", 2100, "If type=t, the end year. Default 2100")
var printTree = flag.String("tree", "", "If specified, 'b' will write a tree showing the hierarchy excluding leaf nodes, l will include leaf nodes")
var csvFile = flag.String("csvFile", "", "The name of a csv file to find hierarchies for")
var apiKey = flag.String("apiKey", "", "The api key ")

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

	if len(*csvFile) > 0 {
		csvparser.FindAllHierarchies(*csvFile, *apiKey)
		os.Exit(0)
	}

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
	validJSON := len(flag.Args()) == 1 && (*hierarchyType == "g" || *hierarchyType == "s")
	validTime := len(flag.Args()) == 0 && *hierarchyType == "t"
	validCsv := len(*csvFile) > 0 && len(*apiKey) > 0
	if !validJSON && !validTime && !validCsv {
		_, exe := filepath.Split(os.Args[0])
		fmt.Println("ONS hierarchy importer. Reads a json representation of a hierarchy or classification, and creates a set of sql insert statements to reconstruct a hierarchy in the db")
		fmt.Println("Please specify a type argument of 'g' (geographical hierarchy) or 's' (structural hierarchy/classification), and the location of the file to parse, e.g.")
		fmt.Println(exe + " -type=g 'http://web.ons.gov.uk/ons/api/data/hierarchies/hierarchy/2011WKWZH.json?apikey=XXXXX&levels=0,1,2'")
		fmt.Println("or")
		fmt.Println(exe + " -type=s 'http://web.ons.gov.uk/ons/api/data/classification/CL_0000641.json?apikey=XXXXX&context=Economic'")
		fmt.Println("or")
		fmt.Println(exe + " -type=g /tmp/localfile.json")
		fmt.Println("Or a type of 't' and a start and end year. Creates a time hierarchy where each year contains months and quarters")
		fmt.Println(exe + " -type=t -start=1900 -end=2100")
		fmt.Println("There is also a 'tree=b or -tree=l option, which will write a tree depiction of the hierarchy excluding leaves (b) or including them (l)")
		fmt.Println()
		fmt.Println("Alternatively, provide two arguments:")
		fmt.Println("  -csvFile=/path/to/file.csv")
		fmt.Println("  -apiKey=yourApiKey")
		fmt.Println("This wil analyse a csv file, retrieve all hierarchies associated with the file and output some information about the dimensions")
		os.Exit(0)
	}
}

func loadHierarchies(t string, file string) []*sql.Hierarchy {
	var hierarchies []*sql.Hierarchy
	switch t {
	case "g":
		fmt.Printf("Importing hierarchies from %s\n", file)
		hierarchies = append(hierarchies, geography.LoadGeography(file))
	case "s":
		fmt.Printf("Importing hierarchies from %s\n", file)
		hierarchies = append(hierarchies, structure.LoadStructure(file)...)
	case "t":
		hierarchies = append(hierarchies, htime.CreateHierarchy(*start, *end))
	}
	return hierarchies
}

func writeSQLForHierarchy(filePrefix string, h *sql.Hierarchy) {
	depth := h.Depth()
	if depth < 3 {
		fmt.Printf("Hierarchy %s has a maximum depth of %d\n", h.ID, depth)
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
	filename = filePrefix + "_arealist.txt"
	fmt.Printf("Creating area list %s\n", filename)
	file, err = os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	sql.WriteLists(file, h)
}

func getWorkingDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
