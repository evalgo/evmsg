.PHONY: all
all: test

test:
	go get -v github.com/jstemmer/go-junit-report
	go build -o go-junit-report github.com/jstemmer/go-junit-report
	go get -v
	go test -v -run=Test_Unit 2>&1 | ./go-junit-report > report.xml

build:
	GOOS=linux GOARCH=amd64 go build -o evmsg.linux.amd64 cmd/evmsg/main.go
	GOOS=darwin GOARCH=amd64 go build -o evmsg.darwin.amd64 cmd/evmsg/main.go
	GOOS=windows GOARCH=amd64 go build -o evmsg.windows.amd64 cmd/evmsg/main.go

.PHONY: clean 
clean:
	find . -name "*~" | xargs rm -fv
	rm -fv go-junit-report report.xml
	rm -fv evmsg.*.amd64

