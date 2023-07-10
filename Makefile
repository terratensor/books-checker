restart: down up

up:
	docker compose up --build -d

down:
	docker compose down --remove-orphans

indexer:
	docker compose exec -it manticore indexer minjust_list
rotate:
	docker compose exec -it manticore indexer minjust_list --rotate
