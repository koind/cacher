run: build up

build:
	docker build -t "cacher:1.0" .

up:
	docker run --name cacher -d -p 8080:8080 -e "CONFIG_PATH=--config=./config/config.toml" cacher:1.0

down:
	docker stop cacher

docs:
	swag init -g ./cmd/cacher/main.go


