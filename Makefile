all:wiki

wiki:wiki.go
	go build $<

clean:
	rm -f wiki *~

run:wiki
	./wiki
