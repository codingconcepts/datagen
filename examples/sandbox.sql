-- REPEAT 5
-- NAME owner
insert into "table" ("c1", "c2") values
{{range $i, $e := $.times_10 }}
	{{if $i}},{{end}}
	(
		'{{stringf "%d %s" 1 10 5 5 "hijklmnop" }}'
	)
{{end}};

-- EOF

