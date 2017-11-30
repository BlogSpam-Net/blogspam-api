
blogspam-api: $(wildcard *.go)
	go build .

clean:
	rm blogspam-api fmt cover.out foo.html || true

fmt:
	go fmt .

test:
	go test -coverprofile fmt

html:
	go test -coverprofile=cover.out .
	go tool cover -html=cover.out -o foo.html
	firefox foo.html
