build:
	docker compose build

run1:
	docker compose run myapp1

run2:
	docker compose run myapp2

clean:
	docker ps -a | awk '{print $$1}' | xargs docker rm
