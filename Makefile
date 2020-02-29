target=$1

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o ${target} ./cmd/${target}
	
build-win:
	GOOS=win GOARCH=amd64 go build -o ${target} ./cmd/${target}

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${target} ./cmd/${target}