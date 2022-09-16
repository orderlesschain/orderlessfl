#!/usr/bin/env bash

source "${PROJECT_ABSOLUTE_PATH}"/env
export OS_PASSWORD_INPUT="$OPENSTACK_PASSWORD_INPUT"
export TF_VAR_public_key="$TERRAFORM_RSA_PUBLIC_KEY"
source "${PROJECT_ABSOLUTE_PATH}"/deployment/terraform/cluster.sh

start=$(date +%s)

pushd "${PROJECT_ABSOLUTE_PATH}"/deployment/terraform/ || exit

terraform destroy --auto-approve

popd || exit

end=$(date +%s)
echo Terrraform done in $(expr $end - $start) seconds.
