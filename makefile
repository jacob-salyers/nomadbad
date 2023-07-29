

build:
	(cd server ; go build && mv nomad ..)

gen:
	./scripts/generate.sh
