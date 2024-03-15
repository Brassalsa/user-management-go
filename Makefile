build:
	go build ./cmd/user-management-go

run:
	go run ./cmd/user-management-go	

dev:
	go build ./cmd/user-management-go | ./user-management-go

clean:
	rm ./user-management-go
