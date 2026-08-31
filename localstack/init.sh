#!/bin/bash

# Script which initializes AWS resources for LocalStack. It gets mounted at /etc/localstack/init/ready.d/ of the
# LocalStack container and automatically runs once the container is healthy.

awslocal s3 mb s3://hike-log
awslocal s3api put-bucket-cors --bucket hike-log --cors-configuration file:///etc/localstack/init/ready.d/cors.json
