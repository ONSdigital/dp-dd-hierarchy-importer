dp-dd-hierarchy-importer
================

A command-line utility to read a hierarchy file (in json format) from the ONS website
and create a set of sql insert statements to recreate that hierarchy in a db.
The ddl to create the tables to import into can be found in table_structure.sql

### Compilation and use

You'll need to have go version (>=1.7) installed. Then run:
	go build
	dp-dd-hierarchy-importer
You should see additional instructions and example usages

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

# Hierarchy notes:

Root of the data api: http://web.ons.gov.uk/ons/api/data/?apikey=API_KEY

### Contexts:
There are 4 contexts, within which most hierarchies reside (geographical info is context-free but can be referenced within a context): 
	http://web.ons.gov.uk/ons/api/data/contexts?apikey=API_KEY
- Census
- Economic
- Socio-Economic
- Social

## (Geographical) Hierarchies
 Can exist outside of any context, but can be relevant to a context (url works with or without "&amp;context=xxx")
 http://web.ons.gov.uk/ons/api/data/hierarchies?apikey=API_KEY

This is the primary type of hierarchy.

## Classifications
 are refined by Context
- Census: http://web.ons.gov.uk/ons/api/data/classifications.xml?apikey=API_KEY&context=Census
- Economic: http://web.ons.gov.uk/ons/api/data/classifications.xml?apikey=API_KEY&context=Economic
- Socio-Economic: none
- Social: http://web.ons.gov.uk/ons/api/data/classifications.xml?apikey=API_KEY&context=Social

Classifications can be hierarchical, having a parent, order, IsTotal and SubTotal flags. 
Some classifications can be flat lists, and we will only import classifications with a deep tree (such as COICOP: http://web.ons.gov.uk/ons/api/data/classification/CL_0000641.xml?apikey=API_KEY&context=Economic).

One possible (flat) classification list is time: http://web.ons.gov.uk/ons/api/data/classification/CL_0000635.xml?apikey=API_KEY&context=Economic

There can be multiple classifications in a single json file, with the same id and overlapping entries.

**Note that the existing data api returns inconsistent json depending on whether there are multiple classifications.** 
If there is a single classification the api returns a single CodeList rather than an array of Codelists. This will break the importer - you must save the json and fix the structure before importing.

## Concepts
 are refined by Context
- Census:  http://web.ons.gov.uk/ons/api/data/concepts.xml?apikey=API_KEY&context=Census&find=&concept=
- Economic: none
- Socio-Economic: none
- Social: none

Concepts are hierarchical, having a parent, level, order and collection count. 
However, the hierarchy is very flat (no entry ever has a grandparent) and seems to be more of a list of multi-value properties. We won't import these.

## Non-hierachical data:

### Collections
 are refined by Context
- Census: http://web.ons.gov.uk/ons/api/data/collections.xml?apikey=API_KEY&context=Census&find=&concept=
- Economic: http://web.ons.gov.uk/ons/api/data/collections.xml?apikey=API_KEY&context=Economic&find=&concept=
- Socio-Economic: none
- Social: http://web.ons.gov.uk/ons/api/data/collections.xml?apikey=API_KEY&context=Social&find=&concept=

Each collection can have a geographic hierarchy but are not themselves hierarchical.


### Datasets
 are refined by Context, are not hierarchical
