GO := /usr/local/go/bin/go
SWAG := /home/jens/go/bin/swag
AIR := /home/jens/go/bin/air
GOCACHE ?= $(CURDIR)/.cache
SWAG_DIRS := cmd/server,internal/modules/auth/api,internal/modules/pokemon/api

.PHONY: dev swag build test tidy fmt

dev:
	GOCACHE=$(GOCACHE) PATH=/usr/local/go/bin:$$PATH $(AIR)

swag:
	PATH=/usr/local/go/bin:$$PATH GOCACHE=$(GOCACHE) $(SWAG) init --generalInfo cmd/server/main.go --dir $(SWAG_DIRS) --useStructName --output internal/docs

build:
	GOCACHE=$(GOCACHE) $(GO) build ./...

test:
	GOCACHE=$(GOCACHE) $(GO) test ./...

tidy:
	GOCACHE=$(GOCACHE) $(GO) mod tidy

fmt:
	$(GO) fmt ./...
