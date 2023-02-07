# Parseable AWS Lambda extension

[![goreportcard](https://goreportcard.com/badge/github.com/parseablehq/aws-lambda-extension)](https://goreportcard.com/report/github.com/parseablehq/aws-lambda-extension)
[![godoc](https://img.shields.io/badge/godoc-reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/parseablehq/aws-lambda-extension)
[![license](https://img.shields.io/github/license/parseablehq/aws-lambda-extension.svg)](https://raw.githubusercontent.com/parseablehq/aws-lambda-extension/master/LICENSE)

[Parseable](https://parseable.io/) is a lightweight, cloud native log observability engine. It can use either a local drive or S3 (and compatible stores) for backend data storage. Parseable is written in Rust and uses Apache Arrow and Parquet as underlying data structures. Parseable consumes up to ~80% lower memory and ~50% lower CPU than Elastic for similar ingestion throughput.

You can deploy Parseable on AWS, GCP, Azure, and on-premises. Refer the [documentation](https://parseable.io/docs) for more details.

## Usage

To use the parseable-lambda-extension with a lambda function, it must be configured as a layer. There are two variants of the extension available: one for `x86_64` architecture and one for `arm64` architecture.

You can add the extension as a layer with the AWS CLI tool:

```sh
$ aws lambda update-code-configuration \
  --function-name MyAwesomeFunction
  --layers "<layer version ARN>"
```

The extension's layer version ARN follows the pattern below.

```sh
# Layer Version ARN Pattern
arn:aws:lambda:<AWS_REGION>:724973952305:layer:parseable-lambda-extension-<ARCH>-<VERSION>:1
```

* `AWS_REGION` - This must match the region of the Lambda function to which you are adding the extension.
* `ARCH` - x86_64 or arm64.
* `VERSION` - The version of the extension you want to use. Current version is v1.0. For current latest release `v1.0`, use the value `v1-0`.

### Configuration

The extension is configurable via environment variables set for your lambda function.

* **PARSEABLE_LOG_URL** - Parseable endpoint URL. It should be set to `https://<parseable-url>/api/v1/ingest`. Change `<parseable-url>` to your Parseable instance URL. (required)
* **PARSEABLE_USERNAME** - Username set for your Parseable instance. (required)
* **PARSEABLE_PASSWORD** - Password set for your Parseable instance. (required)
* **PARSEABLE_LOG_STREAM** - Parseable stream name where you want to ingest logs. (default: ``Lambda Function Name``).

Refer Parseable [installation documentation](https://www.parseable.io/docs/category/installation) for more details.

## Container image lambda

In case if you deploy your lambda as container image, to inject extension as part of your function just copy it to your image:

```Dockerfile
FROM parseable/aws-lambda-extension:latest AS parseable-extension
FROM public.ecr.aws/lambda/python:3.8
# Layer code
WORKDIR /opt
COPY --from=parseable-extension /opt/ .
# Function code
WORKDIR /var/task
COPY app.py .
CMD ["app.lambda_handler"]
```

More details you can find [here](https://aws.amazon.com/blogs/compute/working-with-lambda-layers-and-extensions-in-container-images/).
