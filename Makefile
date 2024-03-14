NAME = lute
MAIN = cmd/$(NAME).go

all: build

build:
	go build $(MAIN)

install:
	go install $(MAIN)

uninstall:
	rm $(GOPATH)/bin/$(NAME)

clean:
	rm ./$(NAME)
