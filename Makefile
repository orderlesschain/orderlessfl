include env

UNAME := $(shell uname)
GO_PATH := $(shell which go)

ifeq ($(UNAME), Linux)
NODE_BUILD_PATH = ${NODE_BUILD_PATH_LINUX}
CLIENT_BUILD_PATH = ${CLIENT_BUILD_PATH_LINUX}
SOURCE = .
PYTHON_PATH = ${PYTHON_PATH_LINUX}
PYTHON_ENV = ${PYTHON_ENV_PATH_LINUX}
endif

ifeq ($(UNAME), Darwin)
NODE_BUILD_PATH = ${NODE_BUILD_PATH_DARWIN}
CLIENT_BUILD_PATH = ${CLIEƒNT_BUILD_PATH_DARWIN}
SOURCE = source
PYTHON_PATH = ${PYTHON_PATH_DARWIN}
PYTHON_ENV = ${PYTHON_ENV_PATH_DARWIN}
endif

ifeq ($(GO_PATH),)
	GO_PATH = ${GO_REMOTE_PATH}
endif

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: tidy
## tidy: tidy up go modules
tidy:
	@${GO_PATH} mod tidy

# helper rule for deployment
check-environment:
ifndef OPENSTACK_PASSWORD_INPUT
    $(error environmet values not set)
endif

.PHONY: build
## build: build the node and the client
build: clean build-node build-client

.PHONY: build-node
## build-node: build the node component (OS-dependent)
build-node: check-environment
	@echo "Building Node ..."
	@${GO_PATH} build -o ${NODE_BUILD_PATH}${NODE_BINARY} ./cmd/node

.PHONY: build-client
## build-client: build the client component (OS-dependent)
build-client:check-environment
	@echo "Building Client ..."
	@${GO_PATH} build -o ${CLIENT_BUILD_PATH}${CLIENT_BINARY} ./cmd/client

.PHONY: build-local-linux
## build-local-linux: build for local Linux on Darwin
build-local-linux: check-environment
	@echo "Building the components on the build remote system..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" ./scripts/build_on_linux.sh

.PHONY: run-playground
## run-playground: run some experimental codes in ./cmd/playground
run-playground:
	@echo "Running Playground ..."
#	@${GO_PATH} run -race ./cmd/playground
	@${GO_PATH} run  ./cmd/playground

.PHONY: clean
## clean: clean for Darwin
clean: check-environment
	@echo "Cleaning"
	@${GO_PATH} clean
	@rm -rf ${NODE_BUILD_PATH}
	@rm -rf ${CLIENT_BUILD_PATH}

.PHONY: clean-all
## clean-all: clean everything
clean-all: check-environment
	@echo "Cleaning"
	@${GO_PATH} clean
	@rm -rf ${NODE_BUILD_PATH_LINUX}
	@rm -rf ${NODE_BUILD_PATH_DARWIN}
	@rm -rf ${CLIENT_BUILD_PATH_LINUX}
	@rm -rf ${CLIENT_BUILD_PATH_DARWIN}

.PHONY: git-commit
## git-commit: pull from master and push to master
git-commit:
	@echo "Commit"
	@git add . ; git commit -m 'auto push';

.PHONY: git-push
## git-push: push to master
git-push:
	@echo "Pushing to git master"
	@git add . ; git commit -m 'auto push'; \
		git push origin main

.PHONY: git-pull
## git-pull: pull from master
git-pull:
	@echo "Pulling from git master"
	@git pull origin main

.PHONY: protos
## protos: generate the protos
protos:
	@echo "Generating the Protos ..."
	@rm -rf ./protos/goprotos; rm -rf ./federatedpython/pythonprotos
	@mkdir ./protos/goprotos ./federatedpython/pythonprotos
	@cp ./federatedpython/__init__for_protos.py. ./federatedpython/pythonprotos/__init__.py
	@protoc -I ./protos ./protos/*.proto  --go_out=./protos/goprotos
	@protoc -I ./protos ./protos/*.proto  --go-grpc_out=require_unimplemented_servers=false:./protos/goprotos
	@${SOURCE} ${PYTHON_ENV}; ${PYTHON_PATH} -m grpc_tools.protoc -I ./protos ./protos/*.proto  --python_out=./federatedpython/pythonprotos --grpc_python_out=./federatedpython/pythonprotos

.PHONY: terraform-deploy-prepare-vms
## terraform-deploy-prepare-vms: deploy the current plan of terraform and prepare remote VMs
terraform-deploy-prepare-vms: check-environment terraform-deploy sleep-shortly prepare-remote-linux-env prepare-tensorflow-env-remote
	@echo "Terraform deploying and preparing remote VMS..."

.PHONY: sleep-shortly
sleep-shortly:
	@echo "Sleeping for 10 seconds ..."
	@sleep 10

.PHONY: terraform-deploy
## terraform-deploy: deploy the current plan of terraform
terraform-deploy: check-environment
	@echo "Terraform deploying ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} ./scripts/terraform_opennebula_deploy.sh

.PHONY: terraform-destroy
## terraform-destroy: destroy the current plan of terraform
terraform-destroy: check-environment
	@echo "Terraform destroying ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} ./scripts/terraform_opennbula_destroy.sh

.PHONY: build-deploy-all-remote
## build-deploy-all-remote: build and deploy remote components using Ansible
build-deploy-all-remote: check-environment
	@echo "Running the bash script for building and deploying components ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" ./scripts/build_on_linux.sh
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" WAN_DRIVE=${WAN_DRIVE_REMOTE} FEDERATED=${FEDERATED_BUILD} ./scripts/deploy_components.sh

.PHONY: build-deploy-all-local
## build-deploy-all-local: build and deploy local components using Ansible
build-deploy-all-local: check-environment
	@echo "Running the bash script for building and deploying components ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" ./scripts/build_on_linux.sh
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" WAN_DRIVE=${WAN_DRIVE_REMOTE} FEDERATED=${FEDERATED_BUILD} ./scripts/deploy_components.sh

.PHONY: destroy-all-remote
## destroy-all-remote: destroy remote components using Ansible
destroy-all-remote: check-environment
	@echo "Running the bash script for destroying components ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" WAN_DRIVE=${WAN_DRIVE_REMOTE} FEDERATED=${FEDERATED_BUILD} ./scripts/destroy_components.sh

.PHONY: destory-all-local
## destroy-all-local: destroy local components using Ansible
destroy-all-local: check-environment
	@echo "Running the bash script for destroying components ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" WAN_DRIVE=${WAN_DRIVE_REMOTE} FEDERATED=${FEDERATED_BUILD} ./scripts/destroy_components.sh

.PHONY: run-federated-experiment-remote
## run-federated-experiment-remote: Run the federated experiment remote
run-federated-experiment-remote:
	@echo "Running the federated experiment remote..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" ./scripts/run_federated_example_experiment.sh

.PHONY: run-federated-experiment-local
## run-federated-experiment-local: Run the federated experiment local
run-federated-experiment-local:
	@echo "Running the federated experiment local..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" ./scripts/run_federated_example_experiment.sh

.PHONY: build-deploy-run-federated-remote
## build-deploy-run-federated-remote: Build deploy federated experiment remote
build-deploy-run-federated-remote: build-deploy-all-remote run-federated-experiment-remote
	@echo "Build deploy and run the example federated remote..."

.PHONY: build-deploy-run-federated-local
## build-deploy-run-federated-remote: Build deploy federated experiment local
build-deploy-run-federated-local: build-deploy-all-local run-federated-experiment-local
	@echo "Build deploy and run the example federated remote..."

.PHONY: run-federated-server
## run-federated-server: Run federated python script
run-federated-server:
	@echo "Running federated Python server ..."
	@${SOURCE} ${PYTHON_ENV};${PYTHON_PATH} ./federatedpython/federated_service.py

.PHONY: freeze-python-requirements
## freeze-python-requirements: Freeze Python venv requirements to a the file
freeze-python-requirements:
	@echo "Freeze Python venv requirements ..."
	@${SOURCE} ${PYTHON_ENV};pip freeze > ./federatedpython/requirements_${UNAME}.txt

.PHONY: install-python-requirements
## install-python-requirements: Install Python venv requirements
install-python-requirements:
	@echo "Install Python venv requirements ..."
	@${SOURCE} ${PYTHON_ENV};pip install -r ./federatedpython/requirements_${UNAME}.txt

.PHONY: vagrant-deploy
## vagrant-deploy: deploy the vagrant locally
vagrant-deploy: check-environment
	@echo "Vagrant deploying locally ..."
	@cd ./deployment/vagrant/; vagrant up

.PHONY: vagrant-destroy
## vagrant-destroy: destroy the vagrant locally
vagrant-destroy: check-environment
	@echo "Vagrant destroy locally ..."
	@cd ./deployment/vagrant/; vagrant destroy -f

.PHONY: vagrant-suspend
## vagrant-suspend: suspend the vagrant locally
vagrant-suspend: check-environment
	@echo "Vagrant suspend locally ..."
	@cd ./deployment/vagrant/; vagrant suspend

.PHONY: vagrant-resume
## vagrant-resume: resume the vagrant locally
vagrant-resume: check-environment
	@echo "Vagrant resume locally ..."
	@cd ./deployment/vagrant/; vagrant resume

.PHONY: prepare-remote-linux-env
## prepare-remote-linux-env: Create the remote linux build env
prepare-remote-linux-env:
	@echo "Preparing the remote linux build env ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" ./scripts/prepare_linux_build_env.sh

.PHONY: prepare-tensorflow-env-remote
## prepare-tensorflow-env-remote: Create the tensorflow env
prepare-tensorflow-env-remote:
	@echo "Preparing the tensorflow env ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" ./scripts/prepare_tensorflow_env.sh

.PHONY: prepare-local-vms
## prepare-local-vms: Install dependencies on local VMs
prepare-local-vms:
	@echo "Preparing the local linux build env ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" ./scripts/prepare_linux_build_env.sh

.PHONY: prepare-tensorflow-env-local
## prepare-tensorflow-env-local: Create the tensorflow env
prepare-tensorflow-env-local:
	@echo "Preparing the tensorflow env ..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="local" ./scripts/prepare_tensorflow_env.sh

.PHONY: run-all-experiments-remote
## run-all-experiments: Run ALL experiments remote
run-all-experiments-remote:
	@echo "Running ALL experiments remote..."
	@PROJECT_ABSOLUTE_PATH=${PROJECT_ABSOLUTE_PATH} BUILD_MODE="remote" ./scripts/run_all_experiments.sh
