
install-cli:
	./scripts/install-cli.sh

uninstall-cli:
	./scripts/uninstall-cli.sh

install-app:
	./scripts/install-app.sh

run: install-cli
	go-away