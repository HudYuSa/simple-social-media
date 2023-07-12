run_server:
	nodemon --exec go run main.go --signal SIGTERM

migrate_up:
	migrate -path database/migrations/ -database "postgres://postgres:jempolbesar@localhost:5432/comments?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migrations/ -database "postgres://postgres:jempolbesar@localhost:5432/comments?sslmode=disable" -verbose down

gen_migrate:
	migrate create -ext sql -dir database/migrations -seq $(name)