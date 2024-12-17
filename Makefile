run:docker_up
	go build -o build/app ./cmd/main.go && ./build/app

docker_up: 
	docker compose up -d
	sleep 1

clean:
	rm -r logs build/*