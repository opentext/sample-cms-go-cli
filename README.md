# Batch updating CMS metadata using a Go CLI

## Overview

This tool demonstrates how you can perform basic Content Metadata Service (CMS) operations using a Go CLI. The intended use of this approach is to be able to perform batch updates on an OCP app.

## Technologies

### Language

Go

### Dependencies

**godotenv** ([github.com/joho/godotenv](https://github.com/joho/godotenv))
Reads in environment variables with overrides in .env.

**cobra** ([github.com/spf13/cobra](https://github.com/spf13/cobra))
Provides cli command management.

**gjson** ([github.com/tidwall/gjson](https://github.com/tidwall/gjson))
Parses and retrieves data from JSON responses.

**retryablehttp** ([github.com/hashicorp/go-retryable](https://github.com/hashicorp/go-retryablehttp))
Handles temporary HTTP errors, especially those caused by rate limiting.

## Prerequisites

* Install Go from [https://go.dev/doc/install](https://go.dev/doc/install).
* Sign up for a trial at [https://developer.opentext.com/imservices/trial](https://developer.opentext.com/imservices/trial)
* Create a tenant in Admin Center.
* Install [VSCode](https://code.visualstudio.com) along with the [OpenText Cloud Developer Tools extension pack](https://marketplace.visualstudio.com/items?itemName=OpenText.ot2-vscode-extension-pack).

## Configuration

* Build the project using `go build`.
* Add an organization profile using the OpenText Cloud Developer Tools and add the tenant you created earlier.
* Deploy the models included in this package to your tenant using the OpenText Cloud Developer Tools. This will add one namespace and one type to CMS.
* You should have been provided with the client id and secret in the VS Code console when deploying the app. Populate the CMS_DEMO_TENANT_ID, CMS_DEMO_CONF_CLIENT_ID and CMS_DEMO_CLIENT_SECRET environment variables now you have this information. This can be done in either the .env file in the root of this project, or in the environment variables configuration of your IDE.

## Usage

Note for Mac and Linux users: `go build` generates an executable file called `planets`. The following commands might require `./planets` to execute. To avoid needing the ./ prefix you should be able to either move the file to an allowed location like `/usr/local/bin` or add the current location to your path variable.

* Run the command `planets info`. This should fetch the token for your app and report `No instances of type un_planet found`
* Run the command `planets create`. This should create planets using the data in `data/planet-data.json`.
* Run the command `planets info` again. This should print the information you just added to CMS. You will notice the `Number of moons` and `Mean temperature` fields are not currently populated.
* Run the command `planets update`. This should update the missing metadata fields for each instance.
* Run the command `planets info` again. This should print the information from CMS and should now include the data for the `Number of moons` and `Mean temperature` fields.
* Run the command `planets delete`. This should delete all the planet instances from CMS.

## Background

### Authentication

Uses the confidential client id and secret with the `client_credentials` grant type. See [authutil](internal/util/auth/authutil.go) for the implementation.

### APIs

This example calls the [Content Metadata Service](https://developer.opentext.com/imservices/products/contentmetadataservice) and uses the following endpoints:

#### Business Object

* GET List object instances
* POST Create new instance
* PUT Update instance details
* DELETE Delete object instance

### Processing CMS responses

For simplicity we have used the gjson library to parse the CMS JSON responses. This could be changed to use Go types. Two approaches that might be helpful in saving time here are:

* Use [Swagger CodeGen](https://swagger.io/tools/swagger-codegen/) to auto-generate type models.
* A third party library like [Marshmallow](https://github.com/PerimeterX/marshmallow) might help with simpler cases where only a small subset of the CMS API is being used.

### Concurrency and Rate Limiting

When running the `planets delete` command, you may notice debug messages in the logs referring to HTTP errors with a `429` status code. For example:

> 2023/10/19 14:05:53 [DEBUG] DELETE https://na-1-dev.api.opentext.com/cms/instances/object/un_planet/722c57e3-2384-45df-bdcf-373f5eb02c23 (status: 429): retrying in 1s (4 left)

This is caused by the rate limiting applied to the CMS API which is currently set to a maximum of 5 requests per second. This is an industry standard practice designed to prevent bursts of request activity, malicious or accidental, from causing system instability.

The info, create and update commands send requests serially. The delete command sends requests in parallel using goroutines. Even without concurrent requests, the delete command is sometimes quick enough to trigger the rate limiting. The problem is solved using the `retryablehttp` module. This will detect the `429` errors and activate a retry strategy using the exponential backoff algorithm. By default this attempts 5 retries per failed HTTP request and delays the wait period for repeated failures. This behaviour is configurable but the sample app just uses the defaults. See the [ioutil](internal/util/io/ioutil.go) `Do` function in this sample for the current implementation.

#### Appendix

Sample planet data obtained from the [Nasa Planetary Fact Sheet](https://nssdc.gsfc.nasa.gov/planetary/factsheet/).
