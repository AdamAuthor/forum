# Имя Docker-образа
IMAGE_NAME := forum
# Путь к файлу с исходным кодом сервера на Golang
SRC_DIR := .
DOCKERFILE := Dockerfile
.PHONY: all
all: build run
.PHONY: build
build:
	docker build -t $(IMAGE_NAME) -f $(DOCKERFILE) $(SRC_DIR)
.PHONY: run
run:
	docker run -d -p 8080:8080 $(IMAGE_NAME)
# Остановка Docker-контейнера
.PHONY: stop
stop:
	docker stop $$(docker ps -q --filter ancestor=$(IMAGE_NAME))
# Очистка собранных файлов и Docker-контейнеров
.PHONY: clean
clean: stop
	docker rm $$(docker ps -a -q --filter ancestor=$(IMAGE_NAME))
	docker rmi $(IMAGE_NAME)