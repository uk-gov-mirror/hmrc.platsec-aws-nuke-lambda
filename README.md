# platsec-aws-nuke-lambda

This is a lambda written in golang.  It runs in a docker container built from a standard alpine image and built with the aws-nuke go binary crafted from https://github.com/rebuy-de/aws-nuke.

There is a ```Makefile``` with targets to build the image, test, tag and push.

TODO: Build in a PR creator for platsec-terraform to deploy the new image.

## Develop

The lambda code can be developed and run locally by:

* Setup up RIE - https://docs.aws.amazon.com/lambda/latest/dg/images-test.html

* ```make clean-build-run```

(then calling the lambda with):
* ```curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"ConfigFilename": "platsec_sandbox_config.yaml", "DryRun": "true"}'```

## Test

To format and run the go tests:

```make test```

## Push

To tag and push the image to the sandbox account go-nuke ECR: 

```make push```