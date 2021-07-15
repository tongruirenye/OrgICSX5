all: check
	go build -o icsd main.go


check:
	go vet github.com/tongruirenye/OrgICSX5/server

clean:
	rm -f icsd

.PHONY: clean check
