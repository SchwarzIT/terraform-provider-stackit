HOSTNAME=registry.terraform.io
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

preview-docs: docs
	@tfplugindocs serve 	

test: 
	@go test $(TEST) || exit 1                                                   
	@echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
	
testacc-token-flow:
	@TF_ACC=1 TF_ACC_LOG=INFO TF_LOG=INFO \
		ACC_TEST_BILLING_REF="$(ACC_TEST_BILLING_REF)" \
		ACC_TEST_USER_EMAIL="$(ACC_TEST_USER_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_EMAIL="$(STACKIT_SERVICE_ACCOUNT_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_TOKEN="$(STACKIT_SERVICE_ACCOUNT_TOKEN)" \
		STACKIT_CUSTOMER_ACCOUNT_ID="$(STACKIT_CUSTOMER_ACCOUNT_ID)" \
		OS_AUTH_URL="$(OS_AUTH_URL)" \
		OS_PASSWORD="$(OS_PASSWORD)" \
		OS_PROJECT_DOMAIN_ID="$(OS_PROJECT_DOMAIN_ID)" \
		OS_PROJECT_NAME="$(OS_PROJECT_NAME)" \
		OS_REGION_NAME="$(OS_REGION_NAME)" \
		OS_TENANT_ID="$(OS_TENANT_ID)" \
		OS_TENANT_NAME="$(OS_TENANT_NAME)" \
		OS_USERNAME="$(OS_USERNAME)" \
		OS_USER_DOMAIN_NAME="$(OS_USER_DOMAIN_NAME)" \
		go test -p 1 -timeout 99999s -v $(TEST)

testacc-key-flow:
	@TF_ACC=1 TF_ACC_LOG=INFO TF_LOG=INFO \
		ACC_TEST_BILLING_REF="$(ACC_TEST_BILLING_REF)" \
		ACC_TEST_USER_EMAIL="$(ACC_TEST_USER_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_KEY_PATH="$(STACKIT_SERVICE_ACCOUNT_KEY_PATH)" \
		STACKIT_PRIVATE_KEY_PATH="$(STACKIT_PRIVATE_KEY_PATH)" \
		STACKIT_CUSTOMER_ACCOUNT_ID="$(STACKIT_CUSTOMER_ACCOUNT_ID)" \
		OS_AUTH_URL="$(OS_AUTH_URL)" \
		OS_PASSWORD="$(OS_PASSWORD)" \
		OS_PROJECT_DOMAIN_ID="$(OS_PROJECT_DOMAIN_ID)" \
		OS_PROJECT_NAME="$(OS_PROJECT_NAME)" \
		OS_REGION_NAME="$(OS_REGION_NAME)" \
		OS_TENANT_ID="$(OS_TENANT_ID)" \
		OS_TENANT_NAME="$(OS_TENANT_NAME)" \
		OS_USERNAME="$(OS_USERNAME)" \
		OS_USER_DOMAIN_NAME="$(OS_USER_DOMAIN_NAME)" \
		go test -p 1 -timeout 99999s -v $(TEST)

ci-testacc:
	@TF_ACC=1 TF_ACC_LOG=INFO TF_LOG=INFO \
		ACC_TEST_BILLING_REF="$(ACC_TEST_BILLING_REF)" \
		ACC_TEST_USER_EMAIL="$(ACC_TEST_USER_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_EMAIL="$(STACKIT_SERVICE_ACCOUNT_EMAIL)" \
		STACKIT_SERVICE_ACCOUNT_TOKEN="$(STACKIT_SERVICE_ACCOUNT_TOKEN)" \
		STACKIT_CUSTOMER_ACCOUNT_ID="$(STACKIT_CUSTOMER_ACCOUNT_ID)" \
		OS_AUTH_URL="$(OS_AUTH_URL)" \
		OS_PASSWORD="$(OS_PASSWORD)" \
		OS_PROJECT_DOMAIN_ID="$(OS_PROJECT_DOMAIN_ID)" \
		OS_PROJECT_NAME="$(OS_PROJECT_NAME)" \
		OS_REGION_NAME="$(OS_REGION_NAME)" \
		OS_TENANT_ID="$(OS_TENANT_ID)" \
		OS_TENANT_NAME="$(OS_TENANT_NAME)" \
		OS_USERNAME="$(OS_USERNAME)" \
		OS_USER_DOMAIN_NAME="$(OS_USER_DOMAIN_NAME)" \
		go test -json -p 1 -timeout 99999s -v $(TEST) > .github/files/analyze-test-output/testoutput; \
		cd .github/files/analyze-test-output && go run analyze.go

dummy:
	@echo $(TEST)

quality:
	@goreportcard-cli -v .

pre-commit: quality
	@go run .github/files/generate-acceptance-tests/main.go  
	@find docs -type f | sort | cat | md5 > .github/files/pre-commit-check/checksum
	@cat .github/workflows/acceptance_test.yml | md5 >> .github/files/pre-commit-check/checksum

ci-pre-commit: ci-docs
	@find docs -type f | sort | cat | md5sum  | cut -d' ' -f1 > .github/files/pre-commit-check/checksum

ci-verify: ci-docs
	@find docs -type f | sort | cat | md5sum  | cut -d' ' -f1 > .github/files/pre-commit-check/checksum-check
	@cat .github/workflows/acceptance_test.yml | md5sum  | cut -d' ' -f1 >> .github/files/pre-commit-check/checksum-check
	@flag=$(false)
	@if cmp -s ".github/files/pre-commit-check/checksum-check" ".github/files/pre-commit-check/checksum"; then \
		rm .github/files/pre-commit-check/checksum-check; \
		echo "files are identical";  \
	else \
		echo "expected:"; \
		cat .github/files/pre-commit-check/checksum; \
		echo "got:"; \
		cat .github/files/pre-commit-check/checksum-check; \
		rm .github/files/pre-commit-check/checksum-check; \
		echo "incorrect checksum. please run 'make pre-commit'"; flag=$(true); \
		exit 1; \
	fi

ci-process-results:
	@go run .github/files/process-test-results/process.go


.PHONY: all ci-docs docs testacc ci-testacc ci-verify pre-commit dummy test quality preview-docs install build ci-process-results ci-pre-commit
