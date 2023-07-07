SHELL := /bin/bash

.PHONY: install generate fmt vet test starttestcontainer removetestcontainer testacc testacccomplete devcheck golangcilint terrafmt terrafmtcheck tflint providerlint importfmtlint

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

generateresource:
	./scripts/generate-resource-source-files.sh

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -parallel=4 ./...

starttestcontainer:
# Start a PingDirectory instance locally to test against
	docker run --name pingdirectory_terraform_provider_container \
		-d -p 1443:1443 \
		-d -p 1389:1389 \
		-e TAIL_LOG_FILES= \
		--env-file "${HOME}/.pingidentity/config" \
		pingidentity/pingdirectory:$${PINGDIRECTORY_PROVIDER_PRODUCT_VERSION:-9.3.0.0}-latest
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingdirectory_terraform_provider_container 2>&1 | grep -q "Setting Server to Available"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
	docker logs pingdirectory_terraform_provider_container 2>&1 | grep -q "Setting Server to Available" || \
		{ echo "PingDirectory container did not become ready in time. Logs:"; docker logs pingdirectory_terraform_provider_container; exit 1; }

removetestcontainer:
	docker rm -f pingdirectory_terraform_provider_container    

testacc:
	PINGDIRECTORY_PROVIDER_HTTPS_HOST=https://localhost:1443 \
	PINGDIRECTORY_PROVIDER_USERNAME=cn=administrator \
	PINGDIRECTORY_PROVIDER_PASSWORD=2FederateM0re \
	PINGDIRECTORY_PROVIDER_INSECURE_TRUST_ALL_TLS=true \
	PINGDIRECTORY_PROVIDER_PRODUCT_VERSION=$${PINGDIRECTORY_PROVIDER_PRODUCT_VERSION:-9.3.0.0} \
	TF_ACC=1 go test -timeout 20m -v ./... -p 4

testacccomplete: removetestcontainer starttestcontainer testacc

devchecknotest: generate install golangcilint tfproviderlint tflint terrafmtlint importfmtlint

devcheck: devchecknotest test testacc

golangcilint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 10m ./...

tfproviderlint: 
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
									-c 1 \
									-AT001.ignored-filename-suffixes=_test.go \
									-AT003=false \
									-XAT001=false \
									-XR004=false \
									-XS002=false ./...

tflint:
	go run github.com/terraform-linters/tflint --recursive

terrafmtlint:
	find ./internal/acctest -type f -name '*_test.go' \
		| sort -u \
		| xargs -I {} go run github.com/katbyte/terrafmt -f fmt {} -v

importfmtlint:
	go run github.com/pavius/impi/cmd/impi --local . --scheme stdThirdPartyLocal ./...