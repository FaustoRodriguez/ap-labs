build:
	go get github.com/hajimehoshi/ebiten
	go get golang.org/x/image
	go build main.go

test: build
	@echo Test 1 - negative values - failed test
	./main -2 -3
	@echo Test 2 - 4 enemies 2 foods
	./main 4 2
	@echo Test 3 - 16 enemies 16 foods
	./main 16
	@echo Test 3 - 45 enemies 30 foods
	./main 45 30

clean:
	rm -rf main