hello:
	@ echo "The following make targets are available"
	@ echo "  help - print this message"
	@ echo "  init - build docker image and install dependencies"
	@ echo "  update - update docker image and dependencies"

init:
	docker build -t remote-manager ./docker
	docker run -it --rm -v "${PWD}":/usr/src/remote-manager remote-manager composer install
	cp .env .env.local
	cp config.json.dist config.json

update:
	docker pull php:7.4-cli-alpine
	docker build -t remote-manager ./docker
	docker run -it --rm -v "${PWD}":/usr/src/remote-manager remote-manager composer install
