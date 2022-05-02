SA:=$(CURDIR)/secrets/sa.json

deps:
	pip3 install semgrep
	go get -v -t -d ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest

clean:
	rm -f report.xml
	rm -f coverage.txt
	rm -f coverage.xml

lint:
	staticcheck ./...

ci: lint security test

security:
	semgrep --error --metrics=on --strict --config=p/golang -o sast.json --json

test: clean
	GOOGLE_APPLICATIO_CREDENTIALS=${SA} go test -cover -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out
