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

-- REPEAT 2
-- NAME two
-- DESCRIPTION Simulate XYZ.
insert into "two" (
	"one_id") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		'{{each "one" "id" $i}}'
	)
{{end}};