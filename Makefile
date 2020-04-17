
install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh

run: install
	go-away