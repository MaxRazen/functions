## help: Show makefile commands
help: Makefile
	@echo "---- Project: MaxRazen/notifier ----"
	@echo " Usage: make COMMAND"
	@echo
	@echo " Available Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## tidy: Ensures fresh go.mod and go.sum
tidy:
	go mod tidy
	go mod verify

## test: Runs tests across the project with no cache
test:
	go test -count=1 ./...

testall:
	./bin/testall.sh
