PACKAGE := $(shell go list -e)
APP_NAME = $(lastword $(subst /, ,$(PACKAGE)))

include gomakefiles/common.mk
include gomakefiles/metalinter.mk
include gomakefiles/upx.mk

SOURCES := $(shell find $(SOURCEDIR) -name '*.go' \
	-not -name '*.pb.go' \
	-not -path './vendor/*')

${APP_NAME}: $(SOURCES)
	go build -ldflags '-X main.Version=${TAG}' -o ${APP_NAME}

RELEASE_SOURCES := $(SOURCES)

include gomakefiles/semaphore.mk

.PHONY: metalinter
metalinter: ${APP_NAME}
	gometalinter --exclude=".*.pb.go" --disable=gotype --vendor --deadline=2m ./...

.PHONY: run
run: ${APP_NAME}
ifndef DAYS
	$(error DAYS parameter must be set - how many days to be processed)
endif
ifndef FLOWDOCK_API_TOKEN
	$(error FLOWDOCK_API_TOKEN parameter must be set - what is your personal API token for flowdock?)
endif
	${APP_NAME} -companyToAnalyze basware -days ${DAYS} -flowToAnalyze portalsos -flowdockApiToken ${FLOWDOCK_API_TOKEN}

.PHONY: prepare
prepare: prepare_metalinter prepare_upx prepare_github_release

.PHONY: clean
clean: clean_common
