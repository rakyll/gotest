all:
	GOOS=linux GOARCH=amd64 go build -o=./bin/gotest_linux
	gsutil cp bin/* gs://jbd-releases
