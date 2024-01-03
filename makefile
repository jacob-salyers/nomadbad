
all: build gen kill start

dev: build gen

local: build genlocal startlocal

build:
	(cd server ; go build && mv nomad ..)

gen:
	./scripts/generate.sh

genlocal: 
	./scripts/generate_local.sh

kill:
	./scripts/kill.sh

start:
	./scripts/start.sh

startlocal:
	./nomad -l
