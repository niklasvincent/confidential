BINARY = confidential
GOARCH = amd64
DIR = apps/confidential
VERSION = 0.1.0

.PHONY: linux
linux:
	cd ${DIR} && GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY} . ; \

.PHONY: darwin
darwin:
	cd ${DIR} && GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY} . ; \

.PHONY: windows
windows:
	cd ${DIR} && GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}.exe . ; \

.PHONY: debian
debian: linux
	fpm -s dir -t deb -n $(BINARY) -v $(VERSION) --prefix /usr/local/bin -C ${DIR} ${BINARY}

.PHONY: rpm
rpm: linux
	fpm --rpm-os linux -s dir -t rpm -n $(BINARY) -v $(VERSION) --prefix /usr/local/bin -C ${DIR} ${BINARY}
	rpm --addsign confidential-${VERSION}-1.x86_64.rpm
