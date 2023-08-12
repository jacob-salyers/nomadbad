
all: build gen kill start

build:
	(cd server ; go build && mv nomad ..)

gen:
	./scripts/generate.sh

kill:
	./scripts/kill.sh

start:
	./scripts/start.sh
