-- REPEAT 10
-- NAME person
insert into "person" ("s", "i", "d", "f32", "f64") values
{{range $i, $e := $.times_100 }}
	{{if $i}},{{end}}
	(
		'{{s 10 10 "p-"}}',
		{{i 1 100}},
		'{{d "2006-01-02T15:04:05+07:00" "2019-01-02T15:04:05+07:00"}}',
		{{f32 1 10}},
		{{f64 1 100}}
	)
{{end}}
returning "id";

-- REPEAT 10
-- NAME pet
insert into "pet" ("pid", "name") values
{{range $i, $e := .times_10 }}
	{{if $i}},{{end}}
	(
		'{{ref "id"}}',
		'{{s 10 10 "a-"}}'
	)
{{end}};

-- EOF
