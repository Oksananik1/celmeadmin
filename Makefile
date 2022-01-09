.PHONY: build
build: bindata
		 env GOOS=linux GOARCH=386 go build  celme/bin/celme

.PHONY: bindata
bindata:
		go mod vendor
		mkdir  -p vendor/celme/bindata/blank/
		go-bindata -pkg blank -o vendor/celme/bindata/blank/bindata.go   blank/static/...  blank/views/...


#5030588209:AAF8Cyvc8CXn683GVDLLz2WHA889CHyUpks
