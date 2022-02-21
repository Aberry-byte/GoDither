compiler=go
flags=build
src=GoDither.go

build:
	$(compiler) $(flags) $(src)

run: build
	./GoDither Doge.png
