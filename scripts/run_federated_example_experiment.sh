#!/usr/bin/env bash

source "${PROJECT_ABSOLUTE_PATH}"/env
start=$(date +%s)

rm -rf "${PROJECT_ABSOLUTE_PATH}"/configs/certs
mkdir "${PROJECT_ABSOLUTE_PATH}"/configs/certs
if [[ $BUILD_MODE == "local" ]]; then
  cp "${PROJECT_ABSOLUTE_PATH}"/configs/endpoints_local.yml "${PROJECT_ABSOLUTE_PATH}"/configs/endpoints.yml
  cp "${PROJECT_ABSOLUTE_PATH}"/configs/certs_local/* "${PROJECT_ABSOLUTE_PATH}"/configs/certs
else
  cp "${PROJECT_ABSOLUTE_PATH}"/configs/endpoints_remote.yml "${PROJECT_ABSOLUTE_PATH}"/configs/endpoints.yml
  cp "${PROJECT_ABSOLUTE_PATH}"/configs/certs_remote/* "${PROJECT_ABSOLUTE_PATH}"/configs/certs
fi

pushd "${PROJECT_ABSOLUTE_PATH}" || exit

go run ./cmd/client -coordinator=true -benchmark=mnist.federatedchain

popd || exit

end=$(date +%s)

echo Experiments executed in $(expr $end - $start) seconds.
