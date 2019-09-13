-- REPEAT 1
-- NAME one
insert into "one" (
    "id",
    "name") values
{{range $i, $e := $.times_1_10 }}
	{{if $i}},{{end}}
	(
		{{int 1 10000}},
		'{{string 5 20 "" ""}}'
	)
{{end}}
returning "id";