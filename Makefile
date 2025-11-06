GOBIN ?= $$(go env GOPATH)/bin

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: install-go-test-coverage
install-go-test-coverage:
	@go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	@${GOBIN}/go-test-coverage --config=./.testcoverage.yml

deps: install-pre-commit install-commitlint install-govulncheck  ## Check dependencies
	$(info Checking and getting dependencies)

install-pre-commit:  ## Install pre-commit if not installed
	@if ! command -v pre-commit >/dev/null 2>&1; then \
		echo "pre-commit not found. Installing..."; \
		sudo apt install pre-commit; \
	else \
		echo "$$(pre-commit --version)"; \
	fi
	pre-commit install -t commit-msg -t pre-commit
	# pre-commit install --install-hooks --overwrite --overwrite-on-success

install-commitlint:  ## Install commitlint if not installed
	@if ! command -v commitlint >/dev/null 2>&1; then \
		echo "commitlint not found. Installing..."; \
		go install github.com/conventionalcommit/commitlint@latest; \
	else \
		echo "$$(commitlint --version)"; \
	fi
	@if ! command -v commitlint >/dev/null 2>&1; then \
		echo "commitlint not found or not accessible yet."; \
		exit 1; \
	fi
	@if [ ! -f .commitlint.yml ] && [ ! -f .commitlint.yaml ] && [ ! -f commitlint.yml ] && [ ! -f commitlint.yaml ]; then \
		echo "No commitlint config file found."; \
		read -p "Do you want to create commitlint config? (y/n): " answer && \
		if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
			echo "Creating commitlint config..."; \
			commitlint config create; \
		else \
			echo "Skipping commitlint config creation."; \
		fi; \
	fi

	@echo "Creating commitlint hook...";
	commitlint init --global;


	@if [ -d .git ]; then \
		echo "Creating commitlint hook..."; \
		commitlint hook create; \
	else \
		echo "Not a git repository."; \
	fi

install-govulncheck:  ## Install govulncheck if not installed
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "govulncheck not found. Installing..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	else \
		echo "$$(govulncheck --version)"; \
	fi
