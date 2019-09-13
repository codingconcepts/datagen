-- REPEAT 10
-- NAME one
insert into "one" (
    "id",
    "name") values
{{range $i, $e := $.times_1 }}
	{{if $i}},{{end}}
	(
		{{int 1 10000}},
		'{{string 5 20 "" ""}}'
	)
{{end}}
returning "id";

-- REPEAT 10
-- NAME two
insert into "two" (
	"one_id") values
{{range $i, $e := $.times_1 }}
	{{if $i}},{{end}}
	(
		'{{each "one" "id" $i}}'
	)
{{end}};