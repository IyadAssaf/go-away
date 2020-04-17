
install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh

install-taskbar:
	./scripts/install-taskbar.sh

run: install
	go-away