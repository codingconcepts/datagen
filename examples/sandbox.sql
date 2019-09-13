-- REPEAT 1
-- NAME one
insert into "one" (
	"id",
    "name") values
{{range $i, $e := ntimes 10 10 }}
	{{if $i}},{{end}}
	(
		{{int 1 10000}},
		'{{wset 1.2 60 9.0 30 5.33 10}}'
	)
{{end}};