up:
	docker-compose up -d --build

down:
	docker-compose down

migrateup:
	docker run -v ${PWD}/migrations:/migrations --network checklist_default migrate/migrate -path=/migrations/ -database "mysql://root:password@tcp(mysql_checklist:3306)/checklist" up $(numb)
migratedown:
	docker run -v ${PWD}/migrations:/migrations --network checklist_default migrate/migrate -path=/migrations/ -database "mysql://root:password@tcp(mysql_checklist:3306)/checklist" down $(numb)

migrateforce:
	docker run -v ${PWD}/migrations:/migrations --network checklist_default migrate/migrate -path=/migrations/ -database "mysql://root:password@tcp(mysql_checklist:3306)/checklist" force $(numb)
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)