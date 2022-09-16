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

commitResults(){
    pushd "${PROJECT_ABSOLUTE_PATH}"
    git add .
    git commit -m 'auto push'
    git pull origin main
    git push origin main
    popd
}

runScaleOne(){
    pushd "${PROJECT_ABSOLUTE_PATH}" || exit

    for BENCHMARK in mnist.scale1 mnist.scale2 mnist.scale3 mnist.scale4 mnist.scale5 mnist.scale6 mnist.scale7 mnist.scale8; do
        for i in {1..1}; do
          echo "Benchmark $BENCHMARK is executing for $i times"
          go run ./cmd/client -coordinator=true -benchmark=${BENCHMARK}
          commitResults
        done
    done

    popd || exit
}

runScaleTwo(){
    pushd "${PROJECT_ABSOLUTE_PATH}" || exit

    for BENCHMARK in mnist.scale9 mnist.scale10 mnist.scale11 mnist.scale12 mnist.scale13 mnist.scale14 mnist.scale15 mnist.scale16; do
        for i in {1..1}; do
          echo "Benchmark $BENCHMARK is executing for $i times"
          go run ./cmd/client -coordinator=true -benchmark=${BENCHMARK}
          commitResults
        done
    done

    popd || exit
}

runScaleThree(){
    pushd "${PROJECT_ABSOLUTE_PATH}" || exit

    for BENCHMARK in mnist.scale17 mnist.scale18 mnist.scale19 mnist.scale20 mnist.scale21 mnist.scale22 mnist.scale23 mnist.scale24; do
        for i in {1..1}; do
          echo "Benchmark $BENCHMARK is executing for $i times"
          go run ./cmd/client -coordinator=true -benchmark=${BENCHMARK}
          commitResults
        done
    done

    popd || exit
}

runScaleFour(){
    pushd "${PROJECT_ABSOLUTE_PATH}" || exit

    for BENCHMARK in mnist.scale25 mnist.scale26 mnist.scale27 mnist.scale28 mnist.scale29 mnist.scale30 mnist.scale31 mnist.scale32; do
        for i in {1..1}; do
          echo "Benchmark $BENCHMARK is executing for $i times"
          go run ./cmd/client -coordinator=true -benchmark=${BENCHMARK}
          commitResults
        done
    done

    popd || exit
}

runAllScale(){
  runScaleOne
  runScaleTwo
  runScaleThree
  runScaleFour
}

runAllScale

runAllScale

runAllScale

end=$(date +%s)
echo Experiments executed in $(expr $end - $start) seconds.
