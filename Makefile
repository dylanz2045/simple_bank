

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres  mypostgres

dropdb:
	docker exec -it postgres dropdb mypostgres

sqlc:
	docker run --rm -v "C:\Users\zdlff\zdl\Project:/src" -w /src sqlc/sqlc generate

  .PHONY:  createdb dropdb sqlc  

