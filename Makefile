PACKAGE_NAME := openview
PACKAGE_DESCRIPTION := A modern image gallery
PACKAGE_VERSION := $(shell git describe --match "v*" | sed -e s/^v//)
PACKAGE_AUTHOR := Felix Kaiser <felix.kaiser@fxkr.net>
PACKAGE_URL := https://github.com/fxkr/openview
PACKAGE_ARCH := amd64

.PHONY: all \
	all-direct \
	build \
	version \
	deps-backend \
	deps-frontend \
	build-backend \
	build-frontend \
	test-gofmt \
	test-govet \
	test-gotest \
	install \
	package-deb

all:
	docker pull quay.io/fxkr/openview-ci:latest
	docker run \
		--tty \
		--interactive \
		--rm \
		--volume "$(realpath .)":/go/src/github.com/fxkr/openview \
		--workdir /go/src/github.com/fxkr/openview \
		quay.io/fxkr/openview-ci:latest \
		make all-direct DESTDIR=/tmp/package

all-direct: build install package-deb

build: deps-backend \
	deps-frontend \
	build-backend \
	build-frontend

version:
	echo ${PACKAGE_VERSION}

deps-backend:
	go get -v -t -d ./...

deps-frontend:
	yarn install

build-backend:
	go build github.com/fxkr/openview/backend/cmd/openview

build-frontend:
	node_modules/.bin/webpack

test-gofmt:
	bash -c "diff -u <(echo -n) <(gofmt -d ./)"

test-govet:
	go vet ./...

test-gotest:
	go test -v ./...

install:
	install -m 0755 -d "$(DESTDIR)/var/lib/openview"
	install -m 0755 -d "$(DESTDIR)/var/cache/openview"
	install -m 0755 -d "$(DESTDIR)/srv/images"
	install -m 0755 -d "$(DESTDIR)/usr/share/openview/static"
	install -m 0644 dist/* "$(DESTDIR)/usr/share/openview/static"
	install -m 0755 -D openview "$(DESTDIR)/usr/bin/openview"
	install -m 0755 -D openview.conf "$(DESTDIR)/usr/share/openview/default.conf"
	install -m 0755 -D openview.service "$(DESTDIR)/lib/systemd/system/openview.service"

package-deb:
	fpm \
		--name          "$(PACKAGE_NAME)" \
		--description   "$(PACKAGE_DESCRIPTION)" \
		--version       "$(PACKAGE_VERSION)" \
		--maintainer    "$(PACKAGE_AUTHOR)" \
		--vendor        "$(PACKAGE_AUTHOR)" \
		--architecture  "$(PACKAGE_ARCH)" \
		--url           "$(PACKAGE_URL)" \
		-s dir \
		-t deb \
		--deb-systemd   "openview.service" \
		--after-install "package/after_install.sh" \
		--depends       "libmagickwand-dev" \
		"$(DESTDIR)/=/"

package-deb-deploy:
	package_cloud push fxkr/openview/debian/jessie "$(PACKAGE_NAME)_$(PACKAGE_VERSION)_$(PACKAGE_ARCH).deb"
	package_cloud push fxkr/openview/debian/stretch "$(PACKAGE_NAME)_$(PACKAGE_VERSION)_$(PACKAGE_ARCH).deb"
