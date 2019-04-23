-- REPEAT 10
-- NAME stock
insert into "stock" ("lin", "location_type", "location_id", "department_id", "quantity", "quantity_in_transit", "quantity_on_order") values
{{range $i, $e := $.times_1000 }}
	{{if $i}},{{end}}
	(
		'{{s 10 10 "l-"}}',
		'{{set "retail" "online"}}',
		{{i 1000 9999}},
		{{i 10000 99999}},
		{{i 1 1000}},
		{{i 1 10}},
		{{i 1 100}}
	)
{{end}}
returning "id";

-- REPEAT 10
-- NAME reservation
insert into "reservation" ("stock_id", "id", "order_number", "line_item_uuid", "quantity", "expires_at", "is_paid") values
{{range $i, $e := $.times_100 }}
	{{if $i}},{{end}}
	(
		'{{ref "stock_id"}}',
		'{{uuid}}',
		'{{s 20 20 "on-"}}',
		'{{s 20 20 "li-"}}',
		{{i 1 5}},
		'{{d "2019-04-23" "2019-04-24" "2006-01-02" }}',
		'{{set "true" "false"}}'
	)
{{end}}
returning "id";