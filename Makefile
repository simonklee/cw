cw: cw.go main.go debug.go store.go
	go build .

test: clean cw
	./cw http://simonklee.org/ http://simonklee.org/article/redis-protocol/

fmt: 
	gofmt -s=true -tabs=false -tabwidth=4 -w .

clean:
	rm -f cw
