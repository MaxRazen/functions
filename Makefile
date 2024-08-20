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

## test: Runs tests across the packages (./pkg) with no cache
test:
	go test -count=1 ./...

## testapp: Runs tests across the project including packages and functions
testall:
	./bin/testall.sh
