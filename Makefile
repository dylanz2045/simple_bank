

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres  mypostgres

dropdb:
	docker exec -it postgres dropdb mypostgres

sqlc:
	docker run --rm -v "C:\Users\zdlff\zdl\Project:/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

migratedown:
	migrate -path db/migration -database "postgresql://postgres:cst4Ever@localhost:5432/mypostgres?sslmode=disable" -verbose down

migrateup:
	migrate -path db/migration -database "postgresql://postgres:cst4Ever@localhost:5432/mypostgres?sslmode=disable" -verbose up

  .PHONY:  createdb dropdb sqlc  test migrateup migratedown

