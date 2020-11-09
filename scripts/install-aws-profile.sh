#!/usr/bin/env bash
set -e

###############################################################################
# Script usage
###############################################################################
function print_usage() {
    echo "USAGE  : ${0} <aws_region> <environment>"
    echo "EXAMPLE: ${0} us-east-2 dev"
    echo "EXAMPLE: ${0} us-east-2 test"
    echo "EXAMPLE: ${0} us-east-2 staging"
    echo "EXAMPLE: ${0} us-east-2 prod"
}

###############################################################################
# Helper function to validate the environment value
###############################################################################
function is_valid_environment() {
    case ${environment} in
      dev|test|staging|prod)
        true
        ;;
      *)
        false
        ;;
    esac
}

###############################################################################
# Helper function to test if AWS keys were set to a value
###############################################################################
function is_aws_keys_set() {
    if [[ ${aws_id} == "" ]] || [[ ${aws_secret} == "" ]]; then
        echo "Unable to set profile '${profile}' - AWS_ACCESS_KEY_ID_${stage} or AWS_SECRET_ACCESS_KEY_${stage} environment variable was not set or empty empty."
        false
    else
        true
    fi
}

###############################################################################
# Command line should have two arguments - region and environment
###############################################################################
if [[ -z "${CI}" ]]; then
  echo "Install AWS profile should only be run in a containerized CI environment"
  exit 1
fi

if [[ "$#" -eq 2 ]]; then
    region=${1}
    environment=${2}

    if [[ ! is_valid_environment ]]; then
        echo "Invalid environment value: '${environment}'"
        print_usage
        exit 1
    fi

    # convert the environment value to upper case and call it the stage {dev => DEV, test => TEST, etc.}
    declare -r stage=`echo ${environment} | tr '[:lower:]' '[:upper:]'`
    declare -n aws_id="AWS_ACCESS_KEY_ID_${stage}"
    declare -n aws_secret="AWS_SECRET_ACCESS_KEY_${stage}"

    if [[ ! is_aws_keys_set ]]; then
        echo "Invalid environment values set for environment: '${environment}'"
        print_usage
        exit 1
    fi
else
    echo "Incorrect number of arguments"
    print_usage
    exit 1
fi

###############################################################################
# Set some local variables
###############################################################################
declare -r profile="proxy-affiliation-service"
declare -r aws_config="${HOME}/.aws/config"
declare -r aws_creds="${HOME}/.aws/credentials"

###############################################################################
# Setup
###############################################################################
function setup() {
    echo "Installing AWS profile: '${profile}' for region: '${region}' in the '${environment}' environment..."
    mkdir -p ${HOME}/.aws && echo "Created ${HOME}/.aws folder..."
}

###############################################################################
# Write the AWS config file
###############################################################################
function write_aws_config() {
    printf "[profile ${profile}]\nregion = ${region}\noutput = json\n" > ${aws_config}
    chmod 600 ${aws_config}
    echo "Wrote AWS config to: ${aws_config}"
    cat ${aws_config}
}

###############################################################################
# Write the AWS credentials file
###############################################################################
function write_aws_creds() {
    printf "[${profile}]\naws_access_key_id=${aws_id}\naws_secret_access_key=${aws_secret}\n" > ${aws_creds}
    chmod 600 ${aws_creds}
    echo "Wrote AWS credentials to: ${aws_creds}"
}

###############################################################################
# Main routine
###############################################################################
function main() {
    setup
    write_aws_config
    write_aws_creds
}

main
