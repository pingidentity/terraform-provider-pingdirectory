#!/bin/bash

set -e

if test -z "${PINGDIRECTORY_ENDPOINT_TO_GENERATE}"; then
	echo "No endpoint specified with PINGDIRECTORY_ENDPOINT_TO_GENERATE environment variable. Exiting."
	exit 0
fi

echo "Generating resource files"

# Check (and start if needed) PD container to read the config model from
if (docker ps | grep -E "pingdirectory_terraform_provider_container" >> /dev/null); then
  echo "Existing PingDirectory Terraform container exists. No need to start a new one..."
else
  echo "No existing PingDirectory Terraform container. Starting a new one..."
	docker run --name pingdirectory_terraform_provider_container \
		-d -p 1389:1389 \
		-e TAIL_LOG_FILES= \
		--env-file "${HOME}/.pingidentity/config" \
		"pingidentity/pingdirectory:${PINGDIRECTORY_TAG:-9.1.0.0-latest}"
fi

# Wait for the instance to become ready, up to 4 minutes
echo "Waiting for PingDirectory to become ready..."
sleep 1
duration=0
while (( duration < 240 )) && ! docker logs pingdirectory_terraform_provider_container 2>&1 | grep -q "Setting Server to Available"; \
do \
		duration=$((duration+1)); \
	sleep 1; \
done

# Fail if the container didn't become ready in time
docker logs pingdirectory_terraform_provider_container 2>&1 | grep -q "Setting Server to Available"

# Run the generator, specifying the endpoints to be generated.
# --endpoint can be specified multiple times to generate multiple endpoints in one run.
java -jar ./bin/pingdirectory-openapi-generator.jar \
    --generateMode terraform \
    --targetDirectory ./ \
    --endpoint "${PINGDIRECTORY_ENDPOINT_TO_GENERATE}"

# Remove the PD container
echo "Stopping and removing PingDirectory container"
docker rm -f pingdirectory_terraform_provider_container
