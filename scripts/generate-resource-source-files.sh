#!/bin/bash

set -e

if test -z "${PINGDIRECTORY_GENERATOR_CONFIG_FILE}"; then
	echo "No config file specified with PINGDIRECTORY_GENERATOR_CONFIG_FILE environment variable. Exiting."
	exit 0
fi

echo "Generating resource files"

# Run the generator
java -jar ./bin/pingdirectory-openapi-generator.jar \
	--configFile "${PINGDIRECTORY_GENERATOR_CONFIG_FILE}" \
	--printParsedConfig
