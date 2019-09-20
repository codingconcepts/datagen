-- REPEAT 10
-- NAME owner
insert into "owner" ("email", "date_of_birth") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		'{{email}}',
		'{{date "1900-01-01" "now" ""}}'
	)
{{end}}
returning "id";

-- REPEAT 20
-- NAME pet
insert into "pet" ("pid", "name", "type") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		'{{ref "owner" "id"}}',
		'{{adj}} {{noun}}',
		'{{wset "dog" 60 "cat" 40}}'
	)
{{end}};
