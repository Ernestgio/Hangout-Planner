#!/bin/bash

# LocalStack init hook to create S3 bucket on startup
echo "Creating S3 bucket: hangout-files"
awslocal s3 mb s3://hangout-files --region ap-southeast-1 2>/dev/null || echo "Bucket already exists"
echo "S3 initialization complete"
