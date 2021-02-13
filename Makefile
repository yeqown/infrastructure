# target=$1

release: build-osx archive
	echo "release done"

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o package/${target} ./cmd/${target}

archive:
	- mkdir -p package/osx
	mv package/${target} package/osx/
	cd package/osx && tar -zcvf ../${target}.tar.gz .
	# clear