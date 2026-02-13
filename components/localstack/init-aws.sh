#!/bin/bash

# LocalStack init hook to create S3 buckets on startup
echo "Creating S3 bucket: hangout-files"
if awslocal s3 mb s3://hangout-files --region ap-southeast-1 2>&1; then
  echo "Bucket hangout-files created successfully"
else
  echo "Bucket hangout-files already exists or creation failed (this is OK if bucket exists)"
fi

echo "Creating S3 bucket: hangout-traces"
if awslocal s3 mb s3://hangout-traces --region ap-southeast-1 2>&1; then
  echo "Bucket hangout-traces created successfully"
else
  echo "Bucket hangout-traces already exists or creation failed (this is OK if bucket exists)"
fi

echo "S3 initialization complete"
