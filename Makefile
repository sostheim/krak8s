NAME      := krak8s
VERSION   := 0.2.85
TYPE      := alpha
COMMIT    := $(shell git rev-parse HEAD)
IMAGE     := quay.io/samsung_cnct/krak8s
TAG       ?= latest
godep=GOPATH=$(shell godep path):${GOPATH}

TIMESTAMP    := $(shell date +"%s")
# GOAGEN_FILES := $(shell find . -maxdepth 1 -name "goa_*.go")
GOAGEN_FILES := openapi.go swagger.go project.go namespace.go application.go cluster.go health.go

design: # Execute goagen to rebuild API, save existing impl files
	@echo "Backing up previous generated soruces with timestamp: $(TIMESTAMP)"
	@for f in $(GOAGEN_FILES); do \
		if [ -f ./$$f ] ; \
		then \
			 mv -f ./$$f ./goagen_backups/$$f.$(TIMESTAMP); \
		fi; \
	done
	@goagen bootstrap -d krak8s/design

build:
	go build -ldflags "-X main.MajorMinorPatch=$(VERSION) \
		-X main.ReleaseType=$(TYPE) \
		-X main.GitCommit=$(COMMIT)"

compile: deps
	@rm -rf build/
	$(GODEP) gox -ldflags "-X main.MajorMinorPatch=$(VERSION) \
		-X main.ReleaseType=$(TYPE) \
		-X main.GitCommit=$(COMMIT) -w" \
	-osarch="linux/386" \
	-osarch="linux/amd64" \
	-osarch="darwin/amd64" \
	-output "build/{{.OS}}_{{.Arch}}/$(NAME)" \
	./...

install:
	godep go install -ldflags "-X main.MajorMinorPatch=$(VERSION) \
		-X main.ReleaseType=$(TYPE) \
		-X main.GitCommit=$(COMMIT) -w"

deps:
	go get github.com/mitchellh/gox
	go get github.com/tools/godep

dist: compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

container:
	$(GODEP) gox -ldflags "-X main.MajorMinorPatch=$(VERSION) \
		-X main.ReleaseType=$(TYPE) \
		-X main.GitCommit=$(COMMIT) -w" \
	-osarch="linux/amd64" \
	-output "build/{{.OS}}_{{.Arch}}/$(NAME)" \
	./...
	docker build --rm --pull --tag $(IMAGE):$(TAG) .

containerprep: deps
	go get github.com/spf13/cobra
	go get
	$(GODEP) gox -ldflags "-X main.MajorMinorPatch=$(VERSION) \
		-X main.ReleaseType=$(TYPE) \
		-X main.GitCommit=$(COMMIT) -w \
		-linkmode external -extldflags -static" \
	-osarch="linux/amd64" \
	-output "build/{{.OS}}_{{.Arch}}/$(NAME)" \
	./...


tag: container
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(COMMIT)

push: tag
	docker push $(IMAGE):$(COMMIT)
	docker push $(IMAGE):$(TAG)

release: dist push
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release samsung-cnct/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*'; \
	git pull

.PHONY: design build compile install deps dist release push tag container
