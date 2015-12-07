SHELL = /bin/bash
TARGETS = islint

# http://docs.travis-ci.com/user/languages/go/#Default-Test-Script
test: assets deps
	go test -v ./...

bench:
	go test -bench .

deps:
	go get ./...

imports:
	go get golang.org/x/tools/cmd/goimports
	goimports -w .

assets: assetutil/bindata.go

assetutil/bindata.go:
	go get -f -u github.com/jteeuwen/go-bindata/...
	go-bindata -o assetutil/bindata.go -pkg assetutil assets/...

vet:
	go vet ./...

cover:
	go test -cover ./...

generate:
	go generate

all: $(TARGETS)

islint: assets generate imports deps
	go build -o islint cmd/islint/main.go

clean:
	rm -f $(TARGETS)
	rm -f islint_*deb
	rm -f islint-*rpm
	rm -rf ./packaging/deb/islint/usr
	rm -f assetutil/bindata.go
	rm -f kind_string.go

deb: $(TARGETS)
	mkdir -p packaging/deb/islint/usr/sbin
	cp $(TARGETS) packaging/deb/islint/usr/sbin
	cd packaging/deb && fakeroot dpkg-deb --build islint .
	mv packaging/deb/islint_*.deb .

rpm: $(TARGETS)
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/rpm/islint.spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/rpm/buildrpm.sh islint
	cp $(HOME)/rpmbuild/RPMS/x86_64/islint*.rpm .
