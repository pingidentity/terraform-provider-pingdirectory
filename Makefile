SHELL := /bin/bash

.PHONY: install generate fmt vet test starttestcontainer removetestcontainer testacc testacccomplete devcheck golangcilint tflint providerlint

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

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
		pingidentity/pingdirectory:$${PINGDIRECTORY_TAG:-9.1.0.0-latest}
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
	TF_ACC=1 go test -timeout 10m -v ./... -p 4

testacccomplete: removetestcontainer starttestcontainer testacc

devcheck: generate install golangcilint test testacccomplete

golangcilint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m ./...

providerlint: 
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
									-c 1 \
									-AT001.ignored-filename-suffixes=_test.go \
									-AT003=false \
									-R009=false \
									-XAT001=false \
									-XR004=false \
									-XS002=false ./internal/...

tflint:
	go run github.com/terraform-linters/tflint --recursive --disable-rule=terraform_required_providers