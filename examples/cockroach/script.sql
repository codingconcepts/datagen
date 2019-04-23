-- REPEAT 10
-- NAME person
insert into "person" ("name", "date_of_birth") values
{{range $i, $e := $.times_100 }}
	{{if $i}},{{end}}
	(
		'{{s 10 10 "p-"}}',
		'{{d "1900-01-01" "2019-04-23" "2006-01-02" }}'
	)
{{end}}
returning "id";

-- REPEAT 10
-- NAME pet
insert into "pet" ("pid", "name") values
{{range $i, $e := .times_10 }}
	{{if $i}},{{end}}
	(
		'{{ref "person_id"}}',
		'{{s 10 10 "a-"}}'
	)
{{end}};

-- EOF
