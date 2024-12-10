build:
	go build -o cmd/out/ssh-automator cmd/ssh-automator/main.go

clean:
	rm -rf cmd/out
