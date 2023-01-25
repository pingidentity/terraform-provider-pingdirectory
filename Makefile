SHELL := /bin/bash

.PHONY: install test starttestcontainer removetestcontainer testacc testacccomplete

default: install

install: fmt
	go install .

fmt:
	go fmt ./...

test:
	go test -parallel=4 ./...

starttestcontainer:
# Start a PingDirectory instance locally to test against
	docker run --name pingdirectory_terraform_acceptance_test \
		-d -p 1443:1443 \
		-e TAIL_LOG_FILES= \
		--env-file "${HOME}/.pingidentity/config" \
		pingidentity/pingdirectory:$${PINGDIRECTORY_TAG:-9.1.0.0-latest}
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 180 )) && ! docker logs pingdirectory_terraform_acceptance_test 2>&1 | grep -q "Setting Server to Available"; \
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
	TF_ACC=1 go test -parallel=4 -timeout 10m -v ./...

testacccomplete:
# Ensure removetestcontainer runs even if an earlier target fails
	${MAKE} starttestcontainer testacc removetestcontainer || ${MAKE} removetestcontainer
