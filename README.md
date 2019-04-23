# datagen
A snappy CLI data generator for databases.

## Installation

``` bash
go get -u /github.com/codingconcepts/datagen
```

## Usage

datagen accepts a configuration file that can execute any number of SQL queries.  SQL queries are built up using Go's [text/template](https://golang.org/pkg/text/template/) language, so can take advantage of multi-row DML for faster inserts etc.

It comes with the following functions to make generating data easy.  If there's a missing generator, please raise an issue or a PR:

## Working example

### Setup

The following example assumes a database called `sandbox` is running in CockroachDB.

Create a table called "person" that defines someone who can have zero or many pets.

``` sql
CREATE TABLE person (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    date_of_birth TIMESTAMP NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC)
)
```

Create a table called "pet" that defines a pet that can belong to a person.

``` sql
CREATE TABLE pet (
    pid UUID NOT NULL,
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (pid ASC, id ASC),
    CONSTRAINT fk_pid_ref_person FOREIGN KEY (pid) REFERENCES person (id)
) INTERLEAVE IN PARENT person (pid)
```

### Script

The following script defines two "blocks":

* One that inserts 10,000 people into the `person` table by executing the first block 10 times (as specified by the `-- REPEAT 10` comment)
* One that inserts 1,000 pets into the `pet` table, referencing IDs of the people previously inserted.

``` sql
-- REPEAT 10
-- NAME person
insert into "person" ("name", "date_of_birth") values
{{range $i, $e := $.times_1000 }}
	{{if $i}},{{end}}
	(
		'{{s 10 10 "p-"}}',
		'{{d "2018-01-02" "2019-01-02" "2006-01-02" }}'
	)
{{end}}
returning "id";

-- REPEAT 10
-- NAME pet
insert into "pet" ("pid", "name") values
{{range $i, $e := .times_100 }}
	{{if $i}},{{end}}
	(
		'{{ref "person_id"}}',
		'{{s 10 10 "a-"}}'
	)
{{end}};

-- EOF

```

### Execute

```
go run main.go -script input.sql --driver postgres --conn postgres://root@localhost:26257/sml?sslmode=disable
```

#### Comments

`-- REPEAT N`

Repeat the block that directly follows the comment N times.  If this comment isn't provided, a block will be executed once.

`-- NAME`

Assigns a given name to the block that directly follows the comment, allowing specific IDs from blocks to be used and not muddled with others.  If this comment isn't provided, no distinction will be made between same-name columns from different tables, so issues will likely arise.  Only omit this for single-block configurations.

`-- EOF`

Causing block parsing to stop, essentially simulating the natural end-of-file.  If this comment isn't provided, the parse will parse all blocks in the script.

#### Custom functions

##### s

Generates a random string between a given minimum and maximum length with an optional prefix:

```
'{{s 5 10 "l-"}}'
```

`s` the name of the function<br/>
`5` the minimum string length including any prefix<br/>
`10` the maximum string length including any prefix<br/>
`"l-"` the prefix<br/>

Note that the apostrophes will wrap the string, turning it into a database string.

##### i

Generates a random integer (32 or 64 bit dependent on architecture) between a minimum and maximum value.

```
{{i 5 10}}
```

`i` the name of the function<br/>
`5` the smallest possible number<br/>
`10` the largest possible number<br/>

##### d

Generates a random date between two dates.  Optionally takes a date format to use, allowing you to target different databases.  Note that all dates will be in UTC format when they enter the database.

```
'{{d "2018-01-02" "2019-01-02" "2006-01-02" }}'
```

`d` the name of the function<br/>
`"2018-01-02"` the earliest possible date<br/>
`"2019-01-02"` the latest possible date<br/>
`"2006-01-02"` the date format to use for parsing the min and max dates and also the date that will be sent to the database (uses Go's `time.Time` formatting rules).  If not provided, this will default to `time.RFC3339` for min, max and output.<br/>

##### f32

```
{{f32 1.23 2.34}}
```

`f32` the name of the function<br/>
`1.23` the smallest possible number<br/>
`2.34` the largest possible number<br/>

##### f64

```
{{f64 1.2345678901 2.3456789012}}
```

`f64` the name of the function<br/>
`1.2345678901` the smallest possible number<br/>
`2.3456789012` the largest possible number<br/>

##### uuid

Generates a random V4 UUID using Google's [uuid](github.com/google/uuid) package.

```
{{uuid}}
```

`uuid` the name of the function.

##### set

Selects a random string from a set of possible options.

```
'{{set "alice" "bob"}}'
```

`set` the name of the function<br/>
`"alice"` the first possible option<br/>
`"bob"` the second possible option<br/>

##### ref

References a random value from a previous block's returned values (cached in memory).  For example, if you have two blocks, one named "person" and another named "pet" and you insert a number of people into the database, returning their IDs, then wish to assign pets to them, you can use the following syntax (assuming you've provided the value "person" for the first block's `-- NAME` comment):

```
'{{ref "person_id"}}',
```

`ref` the name of the function<br/>

#### Helper functions

##### times_1

Time functions can be used to generate multi-line DML.

Execute something once.

```
{{range $i, $e := $.times_1}}
```

##### times_10

Execute something 10 times.

```
{{range $i, $e := $.times_10}}
```

##### times_100

Execute something 100 times.

```
{{range $i, $e := $.times_100}}
```

##### times_1000

Execute something 1,000 times.

```
{{range $i, $e := $.times_1000}}
```

##### times_10000

Execute something 10,000 times.

```
{{range $i, $e := $.times_10000}}
```

##### times_100000

Execute something 100,000 times.

```
{{range $i, $e := $.times_100000}}
```

## Todos

* Runtime test against different types of databases.
* Refactor `parse.Blocks` function to it's easier to read.