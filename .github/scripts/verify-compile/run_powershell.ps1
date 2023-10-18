$ErrorActionPreference = "Stop"
$SALAMI_REPO_DIRECTORY = $args[0]

New-Item -ItemType Directory -Force -Path .\salami
New-Item -ItemType Directory -Force -Path .\terraform
Copy-Item "${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\s3_bucket.sami" -Destination .\salami\s3_bucket.sami
Copy-Item "${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\salami.yaml" -Destination .\salami.yaml

salami compile

if (!(Test-Path -Path .\terraform\s3_bucket.tf)) {
  Write-Host "Error: terraform\s3_bucket.tf does not exist"
  exit 1
}

if (!(Compare-Object (Get-Content .\terraform\s3_bucket.tf) (Get-Content "${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\s3_bucket.tf") -Quiet)) {
  Write-Host "Error: terraform\s3_bucket.tf and ${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\s3_bucket.tf are not the same"
  
  Write-Host "Contents of terraform\s3_bucket.tf:"
  Get-Content .\terraform\s3_bucket.tf

  Write-Host "Contents of ${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\s3_bucket.tf:"
  Get-Content "${SALAMI_REPO_DIRECTORY}\.github\scripts\verify-compile\s3_bucket.tf"

  exit 1
}