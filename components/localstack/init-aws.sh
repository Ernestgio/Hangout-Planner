#!/bin/bash

# LocalStack init hook to create S3 bucket on startup
echo "Creating S3 bucket: hangout-files"
if awslocal s3 mb s3://hangout-files --region ap-southeast-1 2>&1; then
  echo "Bucket created successfully"
else
  echo "Bucket already exists or creation failed (this is OK if bucket exists)"
fi
echo "S3 initialization complete"
