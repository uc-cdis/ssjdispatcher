# SSJDISPATCHER

The SQS S3 Job Dispatcher is designed for centralizing all gen3 jobs in 
the [Gen3 stack](https://gen3.org/). It monitors a SQS queue receiveing CRUD messages
from S3 buckets and determine an action basing on the object url pattern.

For example, an url with a pattern of `s3://bucketname/user.yaml` of the object uploaded to S3 will 
trigger an even in S3 to send a message to a configured SQS. The dispatcher service pulls the message from 
the queue and dispatches an job that pulls `fence` image and run `usersync` job with `fence-create` command.
For the other url with the pattern of `s3://data_upload_bucket/000ed0fb-d1f4-4b80-8d77-0d134bb4c0d6/TARGET-10-PAREBA-09A-01D_GAGTGG_L003.bam`, the service will dispatch an job to compute hashes, and size and register to `indexd` with url.

JobConfig is an json-base string to register an job. An example might be something like this:
```
{
    JOBS": [
          {
            "name": "indexing",
            "pattern": "s3://bucket/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}/.*",
            "deadline": 3600,
            "image": "quay.io/cdis/indexs3client:master",
            "imageConfig": {
              "url": "http://indexd-service/",
              "username": "test",
              "password": "test"
            }
          },
          {
            "name": "usersync",
            "pattern": "s3://bucket/user.yaml",
            "image": "quay.io/cdis/fence:master",
            "imageConfig" :{}
          }
}
```
## API Documentation

[OpenAPI documentation available here.](http://petstore.swagger.io/?url=https://raw.githubusercontent.com/uc-cdis/ssjdispatcher/master/openapis/openapi.yaml)

YAML file for the OpenAPI documentation is found in the `openapis` folder (in the root directory); see the README in that folder for more details.

## Setup

ssjdispatcher is generally deployed in a Gen3 environment and makes use of AWS resources.

For instructions on how to set up ssjdispatcher and the relevant AWS resources with Gen3 cloud-automation, see [here](https://github.com/uc-cdis/cloud-automation/blob/master/doc/kube-setup-ssjdispatcher.md).

For a high-level view of the context in which ssjdispatcher is used with Gen3 see [here](https://github.com/uc-cdis/cloud-automation/blob/master/doc/data_upload/README.md).

### Building From Source

Build the go code with:
```bash
go build -o bin/ssjdispatcher
```

### Building and Running a Docker Image

Build the docker image for ssjdispatcher:
```bash
# Run from root directory
docker build -t ssjdispatcher .
```

Run the docker image:
```bash
docker run -p 8080:8080 ssjdispatcher --port 8080
```
(This command exposes ssjdispatcher on port 8080 in the docker image, and maps port
8080 from the docker image onto 8080 on the host machine.)

## Development

### Tests

Run all the tests:
```bash
go test ./...
```
