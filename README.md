# ![datagen logo](assets/logo.png)

[![Build Status](https://travis-ci.org/codingconcepts/datagen.svg?branch=master)](https://travis-ci.org/codingconcepts/datagen) [![Go Report Card](https://goreportcard.com/badge/github.com/codingconcepts/datagen)](https://goreportcard.com/report/github.com/codingconcepts/datagen)

If you need to generate a lot of random data for your database tables but don't want to spend hours configuring a custom tool for the job, then `datagen` could work for you.

`datagen` takes its instructions from a configuration file. These configuration files can execute any number of SQL queries, taking advantage of multi-row DML for fast inserts and Go's [text/template](https://golang.org/pkg/text/template/) language is used to acheive this.

> "[modelgen](https://github.com/LUSHDigital/modelgen) saves sooo much time. coupled with datagen it just gives you a crazy head start"

> "right now datagen and [modelgen](https://github.com/LUSHDigital/modelgen) are god sends to me"

## Installation

```bash
go get -u github.com/codingconcepts/datagen
```

## Usage

See the [examples](https://github.com/codingconcepts/datagen/tree/master/examples) directory for a CockroachDB example that works using the `make example` command. When running the executable, use the following syntax:

```
datagen -script script.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable
```

`datagen` accepts the following arguments:

| Flag       | Description |
| ---------- | ----------- |
| `-conn`    | The full database connection string (enclosed in quotes) |
| `-driver`  | The name of the database driver to use [postgres, mysql] |
| `-script`  | The full path to the script file to use (enclosed in quotes) |
| `-datefmt` | _(optional)_ `time.Time` format string that determines the format of all database and template dates. Defaults to "2006-01-02" |
| `-debug`   | _(optional)_ If set, the SQL generated will be written to stout. Note that `ref`, `row`, and `each` won't work. |

## Concepts

| Object | Description |
| ------ | ---------- |
| Block  | A block of text within a configuration file that performs a series of operations against a database. |
| Script | A script is a text file that contains a number of blocks. |

### Comments

`datagen` uses Go's [text/template](https://golang.org/pkg/text/template/) engine where possible but where it's not possible to use that, it parses and makes use of comments. The following comments provide instructions to `datagen` during block parsing.

| Comment       | Description |
| ------------- | ----------- |
| `-- REPEAT N` | Repeat the block that directly follows the comment N times. If this comment isn't provided, a block will be executed once. Consider this when using the `ntimes` function to insert a large amount of data. For example `-- REPEAT 100` when used in conjunction with `ntimes 1000` will result in 100,0000 rows being inserted using multi-row DML syntax as per the examples.               |
| `-- NAME`     | Assigns a given name to the block that directly follows the comment, allowing specific rows from blocks to be referenced and not muddled with others. If this comment isn't provided, no distinction will be made between same-name columns from different tables, so issues will likely arise (e.g. `owner.id` and `pet.id` in the examples). Only omit this for single-block configurations. |
| `-- EOF`      | Causing block parsing to stop, essentially simulating the natural end-of-file. If this comment isn't provided, the parse will parse all blocks in the script. |

#### Helper functions

##### ntimes

Expresses the number of multi-row DML statements that will be generated:

```
{{range $i, $e := ntimes 1 10 }}
	{{if $i}},{{end}}
	(
		...something
	)
{{end}}
```

`ntimes` the name of the function.<br/>
`1` the minimum value.<br/>
`10` _(optional)_ the maximum value. If omitted, the number will be exactly equal to the minimum value.<br/>

The following script generates 5 entries into the `one` table and between 5 and 10 entries into the `two` table as a result of the combination of the `-- REPEAT` and `ntimes` configured:

```
-- REPEAT 1
-- NAME one
insert into "one" (
    "id",
    "name") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		{{int 1 10000}},
		'{{string 5 20 "" ""}}'
	)
{{end}}
returning "id";

-- REPEAT 1
-- NAME two
insert into "two" (
	"one_id") values
{{range $i, $e := ntimes 5 10 }}
	{{if $i}},{{end}}
	(
		'{{each "one" "id" $i}}'
	)
{{end}};
```

The `ntimes` and `REPEAT` values for table one's insert totalled 5, so you'll see 5 rows in table one:

| id |
| -- |
| 1977 |
| 2875 |
| 6518 |
| 6877 |
| 9425 |

The `ntimes` and `REPEAT` values for table two's insert totalled 7 (`ntimes` generated 7 and we `REPEATE` once):

| one_id | count |
| ------ | ----- |
| 1977 | 2 |
| 2875 | 1 |
| 6518 | 2 |
| 6877 | 1 |
| 9425 | 1 |

By increasing the `REPEAT` value to 2, we'll generate a total of 14 (`ntimes` is 7 multiplied by two this time):

| one_id | count |
| ------ | ----- |
| 1977 | 3 |
| 2875 | 2 |
| 6518 | 3 |
| 6877 | 3 |
| 9425 | 3 |

##### string

Generates a random string between a given minimum and maximum length with an optional prefix:

```
'{{string 5 10 "l-" "abcABC"}}'
```

`string` the name of the function.<br/>
`5` the minimum string length including any prefix.<br/>
`10` the maximum string length including any prefix.<br/>
`"l-"` the prefix.<br/>
`"abcABC"` _(optional)_ the set of characters to select from.<br/>

Note that the apostrophes will wrap the string, turning it into a database string.

##### stringf

Generates a formatted string using placeholder syntax:

```
'{{stringf "%s.%d@acme.com" 5 10 "abc" 10000 20000}}',
```

`stringf` the name of the function.<br/>
`"%s.%i@acme.com"` the format string.<br/>
`5` the minimum string length for the first string placeholder.<br/>
`10` the minimum string length for the first string placeholder.<br/>
`"abc"` the characters to use for the first string placeholder (leave blank to use defaults).<br/>
`10000` the minimum value for the integer placeholder.<br/>
`20000` the minimum value for the integer placeholder.<br/>

Note that at present only the following verbs are supported:

- %s - a string
- %d - an integer

##### int

Generates a random 64 bit integer between a minimum and maximum value.

```
{{int 5 10}}
```

`int` the name of the function.<br/>
`5` the minimum number to generate.<br/>
`10` the maximum number to generate.<br/>

##### date

Generates a random date between two dates.

```
'{{date "2018-01-02" "now" "" }}'
```

`date` the name of the function.<br/>
`"2018-01-02"` the minimum date to generate.<br/>
`"2019-01-02"` the maximum date to generate.<br/>
`""` the format to use for input dates, left blank to use the value specified by the `-datefmt` flag date.  If overridden, both the minimum and maximum date arguments should be in the overridden format.

Note that `"now"` can be passed to both the minimum and maximum dates if required.

```
'{{date "now" "now" "2006-01-02 15:04:05" }}'
```

`"2006-01-02 15:04:05"` the date format you which to be generated

##### float

Generates a random 64 bit float between a minimum and maximum value.

```
{{float 1.2345678901 2.3456789012}}
```

`float` the name of the function.<br/>
`1.2345678901` the minimum number to generate.<br/>
`2.3456789012` the maximum number to generate.<br/>

##### uuid

Generates a random V4 UUID using Google's [uuid](github.com/google/uuid) package.

```
{{uuid}}
```

`uuid` the name of the function.

##### set

Selects a random value from a set of possible values.

```
'{{set "alice" 1 2.3"}}'
```

`set` the name of the function.<br/>
`"alice"`|`"bob"` etc. the available options to generate from.<br/>

##### wset

Selects a random value from a set of possible values using weighting.

```
'{{wset "a" 60 "b" 30 "c" 10}}'
```

`wset` the name of the function.<br/>
`"a"` the first option.<br/>
`60` a weight of 60 for the first option.<br/>
`"b"` the second option.<br/>
`30` a weight of 30 for the second option.<br/>
`"c"` the third option.<br/>
`10` a weight of 10 for the first option.<br/>

Weights can be any number.

##### fset

Selects a random value from a set of possible values contained within a file.

```
'{{fset "./examples/types.txt"}}'
```

`fset` the name of the function.<br/>
`"./examples/types.txt"` the path to the file containing the options.<br/>

##### ref

References a random value from a previous block's returned values (cached in memory). For example, if you have two blocks, one named "owner" and another named "pet" and you insert a number of owners into the database, returning their IDs, then wish to assign pets to them, you can use the following syntax (assuming you've provided the value "owner" for the first block's `-- NAME` comment):

```
'{{ref "owner" "id"}}',
```

`ref` the name of the function.<br/>

##### row

References a random row from a previous block's returned values and caches it so that values from the same row can be used for other column insert values. For example, if you have two blocks, one named "owner" and another named "pet" and you insert a number of owners into the database, returning their IDs and names, you can use the following syntax to get the ID and name of a random row (assuming you've provided the value "owner" for the first block's `-- NAME` comment):

```
'{{row "owner" "id" $i}}',
'{{row "owner" "name" $i}}'
```

`row` the name of the function.<br/>
`owner` the name of the block whose data we're referencing.<br/>
`id` the name of the owner column we'd like.<br/>
`$i` the group identifier for this insert statement (ensures columns get taken from the same row).<br/>

##### each

Works in a simliar way to `row` but references _sequential_ rows from a previous block's returned values, allowing all of a previous block's rows to have associated rows in a related table, provided the product of `--REPEAT` and `ntimes` is the same as the previous block's.

```
'{{each "owner" "id" $i "pet"}}',
'{{each "owner" "name" $i "pet"}}',
```

`each` the name of the function.<br/>
`owner` the name of the block whose data we're referencing.<br/>
`id` the name of the owner column we'd like.<br/>
`$i` the group identifier for this insert statement (ensures columns get taken from the same row).<br/>

```
{{range $i, $e := ntimes 1}}
	...something
{{end}}
```

## Other database types:

### MySQL

```
datagen -script mysql.sql --driver mysql --conn root@/sandbox
```

With MySQL's lack of a `returning` clause, we instead select a random record from the `person` table when inserting pet records, which is less efficient but provides a workaround.

```sql
-- REPEAT 10
-- NAME pet
insert into `pet` (`pid`, `name`) values
{{range $i, $e := ntimes 100 }}
	{{if $i}},{{end}}
	(
		(select `id` from `person` order by rand() limit 1),
		'{{string 10 10 "a-" ""}}'
	)
{{end}};
```

## Todos

* Ability to generate specific types of data (first name / last name etc).

* Better handling of connection issues during run.

* Integration tests.

* Migrate to travis-ci.com and add coveralls support back in.
