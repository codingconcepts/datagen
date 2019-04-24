-- REPEAT 10
-- NAME owner
insert into "owner" ("name", "date_of_birth") values
{{range $i, $e := $.times_1000 }}
	{{if $i}},{{end}}
	(
		'{{s 10 10 "o-"}}',
		'{{d "1900-01-01" "2019-04-23" "2006-01-02" }}'
	)
{{end}}
returning "id", "name";

-- REPEAT 10
-- NAME pet
insert into "pet" ("pid", "name", "owner_name") values
{{range $i, $e := .times_100 }}
	{{if $i}},{{end}}
	(
		'{{row "owner" "id" $i}}',
		'{{s 10 10 "p-"}}',
		'{{row "owner" "name" $i}}'
	)
{{end}};

-- EOF
