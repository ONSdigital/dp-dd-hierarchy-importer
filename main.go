package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/ONSdigital/dp-dd-hierarchy-importer/geography"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/sql"
	"github.com/ONSdigital/dp-dd-hierarchy-importer/structure"
)

var hierarchyType = flag.String("type", "", "'g' (geographical hierarchy) or 's' (structural hierarchy)")

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Something went wrong: %s\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	checkCommandLineArgs()
	dir := getWorkingDir()

	hierarchies := loadHierarchies(*hierarchyType, flag.Arg(0))

	for _, h := range hierarchies {
		writeSqlForHierarchy(dir, h)
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
		fmt.Println(exe + " -type=s 'http://web.ons.gov.uk/ons/api/data/classification/CL_0001363.json?apikey=XXXXX&context=Census'")
		fmt.Println("or")
		fmt.Println(exe + " -type=g /tmp/localfile.json")
		os.Exit(0)
	}
}

func loadHierarchies(t string, file string) []*sql.Hierarchy {
	var hierarchies []*sql.Hierarchy
	switch t {
	case "g":
		hierarchies = append(hierarchies, geography.LoadGeography(file))
	case "s":
		hierarchies = append(hierarchies, structure.LoadStructure(file)...)
	}
	return hierarchies
}

func writeSqlForHierarchy(dir string, h *sql.Hierarchy) {
	filename := filepath.Join(dir, h.Id+".sql")
	fmt.Printf("Creating sql file %s\n", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	sql.WriteInserts(file, h)
	fmt.Printf("Finished writing %s with %d entries\n", filename, len(h.Entries))
}

func getWorkingDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
