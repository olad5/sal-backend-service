test: 
		go test -count=1 -tags integration  ./tests/integration/

test.verbose: 
		go test -count=1 -tags integration -v  ./tests/integration/

run: 
		PORT=4000 go run cmd/main.go 
