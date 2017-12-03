BINARY = confidential
GOARCH = amd64

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}-linux-${GOARCH} . ; \

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY}-darwin-${GOARCH} . ; \

windows:
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY}-windows-${GOARCH}.exe . ; \
