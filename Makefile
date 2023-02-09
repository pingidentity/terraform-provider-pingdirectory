SHELL := /bin/bash

.PHONY: install generate fmt vet lint test starttestcontainer removetestcontainer testacc testacccomplete devcheck

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

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

test:
	go test -parallel=4 ./...

starttestcontainer:
# Start a PingDirectory instance locally to test against
	docker run --name pingdirectory_terraform_acceptance_test \
		-d -p 1443:1443 \
		--env-file "${HOME}/.pingidentity/config" \
		-e SERVER_PROFILE_URL=https://github.com/henryrecker-pingidentity/pingidentity-server-profiles.git \
		-e SERVER_PROFILE_PATH=debugtrace \
		-e TAIL_LOG_FILES="/opt/out/instance/logs/config-audit.log /opt/out/instance/logs/debug" \
		pingidentity/pingdirectory:$${PINGDIRECTORY_TAG:-9.1.0.0-latest}
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingdirectory_terraform_acceptance_test 2>&1 | grep -q "Setting Server to Available"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
	docker logs pingdirectory_terraform_acceptance_test 2>&1 | grep -q "Setting Server to Available"

removetestcontainer:
	docker rm -f pingdirectory_terraform_acceptance_test    

testacc:
	PINGDIRECTORY_PROVIDER_HTTPS_HOST=https://localhost:1443 \
	PINGDIRECTORY_PROVIDER_USERNAME=cn=administrator \
	PINGDIRECTORY_PROVIDER_PASSWORD=2FederateM0re \
	TF_ACC=1 go test -p 4 -timeout 10m -v ./...

testacccomplete: removetestcontainer starttestcontainer testacc removetestcontainer

devcheck: generate install lint test testacccomplete
