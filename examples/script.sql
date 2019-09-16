-- REPEAT 1
-- NAME owner
insert into "owner" ("email", "date_of_birth") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		'{{stringf "%s.%s@acme.co.uk" 5 5 "abcdefg" 5 5 "hijklmnop" }}',
		'{{date "1900-01-01" "now" ""}}'
	)
{{end}}
returning "id";

-- REPEAT 2
-- NAME pet
insert into "pet" ("pid", "name", "type") values
{{range $i, $e := ntimes 5 }}
	{{if $i}},{{end}}
	(
		'{{ref "owner" "id"}}',
		'{{name}}',
		'{{wset "dog" 60 "cat" 40}}'
	)
{{end}};
