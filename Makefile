# Makefile
SHELL:=/bin/bash
.PHONY: dep

default: dep

dep:
	glide up -u -v -s
	
test-local-pre:
	@mv vendor _vendor 2> /dev/null || exit 0
	@mv $(GOPATH)/src/github.com/bitgaming/go-common/vendor \
		$(GOPATH)/src/github.com/bitgaming/go-common/_vendor 2> /dev/null || exit 0

test-local-post:
	@mv _vendor vendor 2> /dev/null || exit 0
	@mv $(GOPATH)/src/github.com/bitgaming/go-common/_vendor \
		$(GOPATH)/src/github.com/bitgaming/go-common/vendor 2> /dev/null || exit 0