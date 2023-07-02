run_server:
	nodemon --exec go run main.go --signal SIGTERM

migrate_up:
	migrate -path database/migrations/ -database "postgresql://postgres:jempolbesar123@localhost:5432/namadb?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migrations/ -database "postgresql://postgres:jempolbesar123@localhost:5432/namadb?sslmode=disable" -verbose down

gen_migrate:
	migrate create -ext sql -dir database/migrations -seq $(name)