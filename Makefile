TEST = ./rancher
TESTARGS = -v
PROJ = terraform-provider-rancher

ARCHS = amd64 386
OS = windows darwin linux

OUT_DIR = target
BIN_DIR = $(OUT_DIR)/bin
PKG_DIR = $(OUT_DIR)/pkg

PKGS = $(foreach arch,$(ARCHS),$(foreach os,$(OS),$(PKG_DIR)/$(PROJ)_$(os)_$(arch)$(PKG_SUFFIX)))
BINS = $(foreach arch,$(ARCHS),$(foreach os,$(OS),$(BIN_DIR)/$(os)_$(arch)/$(PROJ)))

default: bin

build: fmt get-deps
	@go build ./...

bin: test
	@gox -help >/dev/null 2>&1 ; if [ $$? -ne 2 ]; then \
		go get github.com/mitchellh/gox; \
	fi
	@gox -output="$(BIN_DIR)/{{.OS}}_{{.Arch}}/{{.Dir}}" -arch="$(ARCHS)" -os="$(OS)"

dist:
	@mkdir -p $(PKG_DIR) 2>/dev/null
	@for arch in $(ARCHS); do \
		for os in $(OS); do \
			echo "$(PKG_DIR)/$(PROJ)_$${os}_$${arch}.tar.gz"; \
			tar czf $(PKG_DIR)/$(PROJ)_$${os}_$${arch}.tar.gz -C $(BIN_DIR)/$${os}_$${arch} .; \
		done \
	done

get-deps:
	@go get -t -v ./...

fmt:
	@gofmt -l -w . rancher/

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

test: build
	@TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: fmt build
	@if [[ "$(CATTLE_URL)" == "" || "$(CATTLE_ACCESS_KEY)" == "" || "-z $(CATTLE_SECRET_KEY)" == "" ]]; then \
		echo "ERROR: CATTLE_URL, CATTLE_SECRET_KEY and CATTLE_ACCESS_KEY must be set."; \
		exit 1; \
	fi
	@TF_ACC=1 go test $(TEST) $(TESTARGS) -timeout 120m

clean:
	@go clean
	@rm -rf target/
