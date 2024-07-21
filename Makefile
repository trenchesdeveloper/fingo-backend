migratecreate:
	# Create a new migration file
	migrate create -ext sql -dir db/migrations $(name)

makepostgres:
	docker compose up -d

droppostgres:
	docker compose down

createdb:
	docker exec -it fingo_postgres createdb --username=root --owner=root fingo

dropdb:
	docker exec -it fingo_postgres dropdb fingo
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/fingo?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/fingo?sslmode=disable" -verbose down

sqlc:
	sqlc generate