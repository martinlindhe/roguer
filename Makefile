bench:
	go test -bench=.

run:
	go run cli/app/app.go

views:
	gorazor views views
