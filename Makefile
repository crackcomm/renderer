
REPO ?= bitbucket.org/moovie/renderer

install:
	go install $(REPO)/cmd/renderer
