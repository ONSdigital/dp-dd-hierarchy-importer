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
