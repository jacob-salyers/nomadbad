
all: build gen kill start

local: build gen

build:
	(cd server ; go build && mv nomad ..)

gen:
	./scripts/generate.sh

kill:
	./scripts/kill.sh

start:
	./scripts/start.sh

