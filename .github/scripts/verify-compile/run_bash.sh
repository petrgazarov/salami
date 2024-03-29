#!/bin/bash

set -euo pipefail

SALAMI_REPO_DIRECTORY=$1

mkdir salami
mkdir terraform
cp ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/s3_bucket.sami salami/s3_bucket.sami
cp ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/salami.yaml salami.yaml

salami compile

if [[ ! -f terraform/s3_bucket.tf ]]; then
  echo "Error: terraform/s3_bucket.tf does not exist"
  exit 1
fi

if ! cmp -s terraform/s3_bucket.tf ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/s3_bucket.tf; then
  echo "Error: terraform/s3_bucket.tf and ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/s3_bucket.tf are not the same"
  
  echo "Contents of terraform/s3_bucket.tf:"
  cat terraform/s3_bucket.tf

  echo "Contents of ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/s3_bucket.tf:"
  cat ${SALAMI_REPO_DIRECTORY}/.github/scripts/verify-compile/s3_bucket.tf

  exit 1
fi