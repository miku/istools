SHELL = /bin/bash
TARGETS = islint iscov islabel

# find go-bindata executable on vm
export PATH := /home/vagrant/bin:$(PATH)

# http://docs.travis-ci.com/user/languages/go/#Default-Test-Script
test: assets generate deps
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
	go get -f -u golang.org/x/tools/cmd/stringer
	go generate

all: $(TARGETS)

islint: assets imports generate deps
	go build -o islint cmd/islint/main.go

iscov: assets imports generate deps
	go build -o iscov cmd/iscov/main.go

islabel: assets imports generate deps
	go build -o islabel cmd/islabel/main.go

clean:
	rm -f $(TARGETS)
	rm -f istools_*deb
	rm -f istools-*rpm
	rm -rf ./packaging/deb/istools/usr
	rm -f assetutil/bindata.go
	rm -f kind_string.go

deb: $(TARGETS)
	mkdir -p packaging/deb/istools/usr/sbin
	cp $(TARGETS) packaging/deb/istools/usr/sbin
	cd packaging/deb && fakeroot dpkg-deb --build istools .
	mv packaging/deb/istools_*.deb .

rpm: $(TARGETS)
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/rpm/istools.spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/rpm/buildrpm.sh istools
	cp $(HOME)/rpmbuild/RPMS/x86_64/istools*.rpm .

PORT = 2222
SSHCMD = ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p $(PORT)
SCPCMD = scp -o port=$(PORT) -o StrictHostKeyChecking=no -i vagrant.key

# Helper to build RPM on a RHEL6 VM, to link against glibc 2.12
vagrant.key:
	curl -sL "https://raw.githubusercontent.com/mitchellh/vagrant/master/keys/vagrant" > vagrant.key
	chmod 0600 vagrant.key

# Don't forget to vagrant up :) - and add your public key to the guests authorized_keys
setup: vagrant.key
	$(SSHCMD) "sudo yum install -y sudo yum install http://ftp.riken.jp/Linux/fedora/epel/6/i386/epel-release-6-8.noarch.rpm"
	$(SSHCMD) "sudo yum install -y golang git rpm-build gcc-c++"
	$(SSHCMD) "mkdir -p /home/vagrant/src/github.com/miku"
	$(SSHCMD) "cd /home/vagrant/src/github.com/miku && git clone /vagrant/.git istools"

rpm-compatible: vagrant.key
	$(SSHCMD) "GOPATH=/home/vagrant go get -f -u github.com/jteeuwen/go-bindata/... golang.org/x/tools/cmd/goimports"
	$(SSHCMD) "cd /home/vagrant/src/github.com/miku/istools && git pull origin master && pwd && GOPATH=/home/vagrant make clean rpm"
	$(SCPCMD) vagrant@127.0.0.1:/home/vagrant/src/github.com/miku/istools/*rpm .
