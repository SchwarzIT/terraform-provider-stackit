HOSTNAME=github.com
NAMESPACE=schwarzit
NAME=stackit
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_arm64
OS_ARCH_DOCKER=linux_amd64
TEST?=$$(go list ./... | grep -v 'vendor')

default: install

build:
	go build -o ${BINARY} 

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

docs:
	@tfplugindocs generate --rendered-provider-name "STACKIT" --provider-name stackit
	@find . -name 'index.md' -exec sed -i '' 's/page_title: "stackit Provider"/page_title: "STACKIT Provider"/g' {} \;
	@find . -name 'index.md' -exec sed -i '' 's/# stackit Provider/# STACKIT Provider/g' {} \;

ci-docs:
	@${GITHUB_WORKSPACE}/tfplugindocs generate --rendered-provider-name "STACKIT" --provider-name stackit
	@find . -name 'index.md' -exec sed -i 's/page_title: "stackit Provider"/page_title: "STACKIT Provider"/g' {} \;
	@find . -name 'index.md' -exec sed -i 's/# stackit Provider/# STACKIT Provider/g' {} \;
	@cat docs/index.md

preview-docs: docs
	@tfplugindocs serve 	

test: 
	@go test $(TEST) || exit 1                                                   
	@echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
	
testacc:
	@TF_ACC=1 TF_ACC_LOG=INFO TF_LOG=INFO \
		ACC_TEST_BILLING_REF="$(ACC_TEST_BILLING_REF)" \
		ACC_TEST_USER_EMAIL="$(ACC_TEST_USER_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_EMAIL="$(STACKIT_SERVICE_ACCOUNT_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_TOKEN="$(STACKIT_SERVICE_ACCOUNT_TOKEN)" \
		STACKIT_CUSTOMER_ACCOUNT_ID="$(STACKIT_CUSTOMER_ACCOUNT_ID)" \
		go test -p 1 -timeout 99999s -v $(TEST)

quality:
	@goreportcard-cli -v .

pre-commit: docs quality
	@find docs -type f -exec md5 {} \; | sort -k 2 | md5 > .github/files/pre-commit-check/checksum
	@cat .github/workflows/acceptance_test.yml | md5 >> .github/files/pre-commit-check/checksum

ci-verify: ci-docs
	@find docs -type f -exec md5sum {} \; | sort -k 2 | md5sum > .github/files/pre-commit-check/checksum-check
	@cat .github/workflows/acceptance_test.yml | md5sum >> .github/files/pre-commit-check/checksum-check
	@flag=$(false)
	@if cmp -s ".github/files/pre-commit-check/checksum-check" ".github/files/pre-commit-check/checksum"; then \
		rm .github/files/pre-commit-check/checksum-check; \
		echo "files are identical";  \
	else \
		rm .github/files/pre-commit-check/checksum-check; \
		echo "incorrect checksum. please run 'make pre-commit'"; flag=$(true); \
		exit 1; \
	fi

.PHONY: all docs testacc ci-verify pre-commit
