#!/bin/bash

# Script which initializes AWS resources for LocalStack. It gets mounted at /etc/localstack/init/ready.d/ of the
# LocalStack container and automatically runs once the container is healthy.

# Create S3 buckets (one for app, one for testing)
awslocal s3 mb s3://hike-log
awslocal s3 mb s3://hike-log-test-bucket 

# Apply CORS and permissions policy to app bucket
awslocal s3api put-bucket-cors --bucket hike-log --cors-configuration file:///etc/localstack/init/ready.d/cors.json
awslocal s3api put-bucket-policy --bucket hike-log --policy file:///etc/localstack/init/ready.d/public-read-policy.json
