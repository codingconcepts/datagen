-- REPEAT 1
-- NAME one
insert into "one" (
	"id",
    "name") values
{{range $i, $e := ntimes 10 10 }}
	{{if $i}},{{end}}
	(
		{{int 1 10000}},
		'{{set "a" 1 2.3}}'
	)
{{end}};