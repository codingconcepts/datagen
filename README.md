![datagen logo](assets/logo.png)
================================

If you need to generate a lot of random data for your database tables but don't want to spend hours configuring a custom tool for the job, then `datagen` could work for you.

`datagen` takes its instructions from a configuration file.  These configuration files can execute any number of SQL queries, taking advantage of multi-row DML for fast inserts and use Go's [text/template](https://golang.org/pkg/text/template/) language to acheive this.

## Installation

``` bash
go get -u github.com/codingconcepts/datagen
```

## Usage

See the [examples](https://github.com/codingconcepts/datagen/tree/master/examples) directory for an example that works using the `make example` command.  When running the executable, use the following syntax:

```
datagen -script input.sql --driver postgres --conn postgres://root@localhost:26257/sml?sslmode=disable
```

`datagen` accepts the following arguments:

| Flag  | Description |
| ------------- | ------------- |
| -conn | The full database connection string (enclosed in quotes) |
| -datefmt | An optional string that determines the format of all database dates |
| -driver | The name of the database driver to use [postgres, mysql] |
| -script | The full path to the script file to use (enclosed in quotes) |


## Concepts

| Object  | Description |
| ------------- | ------------- |
| Block | A block of text within a configuration file that performs a series of operations against a database. |
| Script  | A script is a text file (typically called `input.sql` but this is optional) that contains a number of blocks. |

### Comments

`datagen` uses the Go text/templating engine where possible but where it's not possible to use that, it makes use of comments.  The following comments provide instructions to `datagen` during block parsing. 

`-- REPEAT N`

Repeat the block that directly follows the comment N times.  If this comment isn't provided, a block will be executed once.  Consider this when using the `.times_*` helpers to insert a large amount of data.  For example `-- REPEAT 100` when used in conjunction with `.times_1000` will result in 100,0000 rows being inserted using multi-row DML syntax as per the examples.

`-- NAME`

Assigns a given name to the block that directly follows the comment, allowing specific rows from blocks to be referenced and not muddled with others.  If this comment isn't provided, no distinction will be made between same-name columns from different tables, so issues will likely arise (e.g. `owner.id` and `pet.id` in the examples).  Only omit this for single-block configurations.

`-- EOF`

Causing block parsing to stop, essentially simulating the natural end-of-file.  If this comment isn't provided, the parse will parse all blocks in the script.

#### Custom functions

##### s

Generates a random string between a given minimum and maximum length with an optional prefix:

```
'{{string 5 10 "l-"}}'
```

`string` the name of the function<br/>
`5` the minimum string length including any prefix<br/>
`10` the maximum string length including any prefix<br/>
`"l-"` the prefix<br/>

Note that the apostrophes will wrap the string, turning it into a database string.

##### i

Generates a random 64 bit integer between a minimum and maximum value.

```
{{int 5 10}}
```

`int` the name of the function<br/>
`5` the minimum number to generate<br/>
`10` the maximum number to generate<br/>

##### d

Generates a random date between two dates.

```
'{{date "2018-01-02" "2019-01-02" }}'
```

`date` the name of the function<br/>
`"2018-01-02"` the minimum date to generate<br/>
`"2019-01-02"` the maximum date to generate<br/>

##### f

Generates a random 64 bit float between a minimum and maximum value. 

```
{{float 1.2345678901 2.3456789012}}
```

`float` the name of the function<br/>
`1.2345678901` the minimum number to generate<br/>
`2.3456789012` the maximum number to generate<br/>

##### uuid

Generates a random V4 UUID using Google's [uuid](github.com/google/uuid) package.

```
{{uuid}}
```

`uuid` the name of the function.

##### set

Selects a random string from a set of possible options.

```
'{{set "alice" "bob" "carol"}}'
```

`set` the name of the function<br/>
`"alice"`|`"bob"` etc. the available options to generate from.<br/>

##### ref

References a random value from a previous block's returned values (cached in memory).  For example, if you have two blocks, one named "owner" and another named "pet" and you insert a number of owners into the database, returning their IDs, then wish to assign pets to them, you can use the following syntax (assuming you've provided the value "owner" for the first block's `-- NAME` comment):

```
'{{ref "owner" "id"}}',
```

`ref` the name of the function<br/>

##### row

References a random row from a previous block's returned values and caches it so that values from the same row can be used for other column insert values.  For example, if you have two blocks, one named "owner" and another named "pet" and you insert a number of owners into the database, returning their IDs and names, you can use the following syntax to get the ID and name of a random row (assuming you've provided the value "owner" for the first block's `-- NAME` comment):

```
'{{row "owner" "id" $i}}',
'{{row "owner" "name" $i}}'
```

`row` the name of the function<br/>
`owner` the name of the block whose data we're referencing<br/>
`id` the name of the owner column we'd like<br/>
`$i` the group identifier for this insert statement.<br/>

#### Helper functions

##### times_*

Time functions can be used to generate multi-line DML.  The number after the underscore denotes the number of times something will be repeated.  Possible numbers are 1, 10, 100, 1000, 10000, and 100000.

```
{{range $i, $e := $.times_1}}
	...something
{{end}}
```

## Other database types:

### MySQL

```
datagen -script mysql.sql --driver mysql --conn root@/sandbox
```

With MySQL's lack of a `returning` clause, we instead select a random record from the `person` table when inserting pet records, which is less efficient but provides a workaround.

``` sql
-- REPEAT 10
-- NAME pet
insert into `pet` (`pid`, `name`) values
{{range $i, $e := .times_100 }}
	{{if $i}},{{end}}
	(
		(select `id` from `person` order by rand() limit 1),
		'{{string 10 10 "a-"}}'
	)
{{end}};

-- EOF
```

## Todos

* Refactor `parse.Blocks` function to it's easier to read.
