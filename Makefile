run\:piximan:
	go run cmd/piximan/main.go

run\:piximanctl:
	go run cmd/piximanctl/main.goma

build\:piximan:
	go build -o ./bin/piximan ./cmd/piximan/main.go

build\:piximanctl:
	go build -o ./bin/piximanctl ./cmd/piximanctl/main.go

build:
	make build:piximan
	make build:piximanctl
