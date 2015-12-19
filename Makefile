bench:
	go test -bench=.

run:
	go run cli/app/app.go && imgcat roller.png && imgcat island.png && imgcat island_col.png
