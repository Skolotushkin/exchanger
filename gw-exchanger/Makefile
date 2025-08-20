APP_NAME=gw-exchanger

build:
	GOOS=linux GOARCH=amd64 go build -o main ./cmd

run:
	go run ./cmd -c config.env

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run --rm -it \
		--env-file config.env \
		-p 50051:50051 \
		$(APP_NAME):latest