test: kld.test
	sudo ./kld.test -test.v -test.coverprofile=coverage.out

kld.test: *.go
	go test -o kld.test -cover -c -v github.com/go-freebsd/kld

cover: coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	@-rm -f kld.test

