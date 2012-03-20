cw: cw.go
	go build .

test: clean cw
	./cw https://raw.github.com/simonz05/cw/master/cw.go

fmt: 
	gofmt -s=true -tabs=false -tabwidth=4 -w .

clean:
	rm -f cw
