restart: down up

up:
	docker compose up --build -d

down:
	docker compose down --remove-orphans

indexer:
	docker compose exec -it manticore indexer minjust_list
rotate:
	docker compose exec -it manticore indexer minjust_list --rotate

build:
	GOOS=windows GOARCH=amd64 go build -o ./dist/books-checker.exe ./app/*.go
	GOOS=linux GOARCH=amd64 go build -o ./dist/books-checker.linux.amd64 ./app/*.go
