#!/bin/bash

set -e
echo "Generating resource files"

# Start a PD container to read the config model from
echo "Starting PingDirectory container to read config model from"
docker run --name pingdirectory_terraform_generator \
	-d -p 1389:1389 \
	-e TAIL_LOG_FILES= \
	--env-file "${HOME}/.pingidentity/config" \
	"pingidentity/pingdirectory:${PINGDIRECTORY_TAG:-9.1.0.0-latest}"

# Wait for the instance to become ready, up to 3 minutes
echo "Waiting for PingDirectory to become ready..."
sleep 1
duration=0
while (( duration < 240 )) && ! docker logs pingdirectory_terraform_generator 2>&1 | grep -q "Setting Server to Available"; \
do \
    duration=$((duration+1)); \
	sleep 1; \
done

# Fail if the container didn't become ready in time
docker logs pingdirectory_terraform_generator 2>&1 | grep -q "Setting Server to Available"

# Run the generator, specifying the endpoints to be generated
java -jar ../../bin/pingdirectory-openapi-generator.jar \
    --generateMode terraform \
    --targetDirectory ../../ \
    --endpoint global-configuration \
    --endpoint location \
    --endpoint root-dn \
    --endpoint server-instance \
    --endpoint trust-manager-provider

# Remove the PD container
echo "Stopping and removing PingDirectory container"
docker rm -f pingdirectory_terraform_generator
