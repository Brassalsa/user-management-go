build:
	go build ./cmd/user-management-go

run:
	go run ./cmd/user-management-go	

dev: build
	 ./user-management-go

clean:
	rm ./user-management-go
