#!/usr/bin/env bash

source "${PROJECT_ABSOLUTE_PATH}"/env
export OS_PASSWORD_INPUT="$OPENSTACK_PASSWORD_INPUT"
export TF_VAR_public_key="$TERRAFORM_RSA_PUBLIC_KEY"
source "${PROJECT_ABSOLUTE_PATH}"/deployment/terraform/cluster.sh

start=$(date +%s)

TFVARS="terraform_4_orgs_8_clients.tfvars"

cp "${PROJECT_ABSOLUTE_PATH}"/contractsbenchmarks/networks/"${TFVARS}" "${PROJECT_ABSOLUTE_PATH}"/deployment/terraform/terraform.tfvars

pushd "${PROJECT_ABSOLUTE_PATH}"/deployment/terraform/ || exit

terraform apply --auto-approve

popd || exit

end=$(date +%s)
echo Terrraform done in $(expr $end - $start) seconds.
