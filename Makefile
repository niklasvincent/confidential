BINARY = confidential
GOARCH = amd64
DIR = apps/confidential

.PHONY: linux
linux:
	cd ${DIR} && GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}-linux-${GOARCH} . ; \

.PHONY: darwin
darwin:
	cd ${DIR} && GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY}-darwin-${GOARCH} . ; \

.PHONY: windows
windows:
	cd ${DIR} && GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}-windows-${GOARCH}.exe . ; \
