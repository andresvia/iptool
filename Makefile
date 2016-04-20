all:

docs: README.md

README.md: iptool.go
	go build
	echo '# iptool' > README.md
	./iptool --help | sed 's|^|   |' >> README.md

.PHONY: all docs
