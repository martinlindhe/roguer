.PHONY: views

bench:
	go test -bench=.

run:
	go run cli/app/app.go

views:
	gorazor views views

shrink-png:
	find ./public/img -name '*.png' -print0 | xargs -0 -n1 shrink-png.sh
