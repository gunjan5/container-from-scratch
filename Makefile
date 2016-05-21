default: image

get:
	goimports -w .
	go get -t -d -v ./...

src:
	docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.5.1 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o app'

fmt:
	gofmt -w .
	#TODO: go lint, go vet
test:
	go test -v -race ./...
	go test -cover -v ./...

image:
	docker build -t gunjan5/container-from-scratch .

run:
	docker run --rm -it gunjan5/container-from-scratch

build:
	GOOS=linux GOARCH=amd64 go build -o app

depsave:
	godep save
	
depupdate:
	go get -t -v ./...
	godep update