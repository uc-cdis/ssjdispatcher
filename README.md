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

### Quickstart with Helm

You can now deploy individual services via Helm! 

If you are looking to deploy all Gen3 services, that can be done via the Gen3 Helm chart. 
Instructions for deploying all Gen3 services with Helm can be found [here](https://github.com/uc-cdis/gen3-helm#readme).

To deploy the ssjdispatcher service:
```bash
helm repo add gen3 https://helm.gen3.org
helm repo update
helm upgrade --install gen3/ssjdispatcher
```
These commands will add the Gen3 helm chart repo and install the ssjdispatcher service to your Kubernetes cluster. 

Deploying ssjdispatcher this way will use the defaults that are defined in this [values.yaml file](https://github.com/uc-cdis/gen3-helm/blob/master/helm/ssjdispatcher/values.yaml)
You can learn more about these values by accessing the ssjdispatcher [README.md](https://github.com/uc-cdis/gen3-helm/blob/master/helm/ssjdispatcher/README.md)

If you would like to override any of the default values, simply copy the above values.yaml file into a local file and make any changes needed. 

To deploy the service independant of other services (for testing purposes), you can set the .postgres.separate value to "true". This will deploy the service with its own instance of Postgres:
```bash
  postgres:
    separate: true
```

You can then supply your new values file with the following command: 
```bash
helm upgrade --install gen3/ssjdispatcher -f values.yaml
```

If you are using Docker Build to create new images for testing, you can deploy them via Helm by replacing the .image.repository value with the name of your local image. 
You will also want to set the .image.pullPolicy to "never" so kubernetes will look locally for your image. 
Here is an example:
```bash
image:
  repository: <image name from docker image ls>
  pullPolicy: Never
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""
```

Re-run the following command to update your helm deployment to use the new image: 
```bash
helm upgrade --install gen3/ssjdispatcher
```

You can also store your images in a local registry. Kind and Minikube are popular for their local registries:
- https://kind.sigs.k8s.io/docs/user/local-registry/
- https://minikube.sigs.k8s.io/docs/handbook/registry/#enabling-insecure-registries
