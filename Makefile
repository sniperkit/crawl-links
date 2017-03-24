
clean:
	rm -rf dist

build: clean
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o ./dist/crawl-links ./main.go
