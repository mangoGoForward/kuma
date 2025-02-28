#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(dirname -- "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/../common.sh"

GOARCH=(amd64 arm64)

# first component is the system - must map to valid $GOOS values
# if present, second is the distribution and third is the envoy distribution
# without a distribution we package only kumactl as a static binary
DISTRIBUTIONS=(linux:debian:alpine linux:ubuntu:alpine linux:rhel:centos linux:centos:centos darwin:darwin:darwin linux)

PULP_HOST="https://api.pulp.konnect-prod.konghq.com"
PULP_PACKAGE_TYPE="mesh"
PULP_DIST_NAME="alpine"
[ -z "$RELEASE_NAME" ] && RELEASE_NAME="kuma"
ENVOY_VERSION=1.21.1
[ -z "$KUMA_CONFIG_PATH" ] && KUMA_CONFIG_PATH=pkg/config/app/kuma-cp/kuma-cp.defaults.yaml
CTL_NAME="kumactl"

function get_envoy() {
  local distro=$1
  local envoy_distro=$2

  local status
  status=$(curl -L -o build/envoy-"$distro" \
    --write-out '%{http_code}' --silent --output /dev/null \
    "https://download.konghq.com/mesh-alpine/envoy-$ENVOY_VERSION-$envoy_distro")

  if [ "$status" -ne "200" ]; then msg_err "Error: failed downloading Envoy"; fi
}

# create_kumactl_tarball packages only kumactl
function create_kumactl_tarball() {
  local arch=$1
  local system=$2

  msg ">>> Packaging ${RELEASE_NAME} static kumactl binary for $system-$arch..."
  msg

  make GOOS="$system" GOARCH="$arch" build/kumactl

  local dest_dir=build/$RELEASE_NAME-$arch
  local kuma_dir=$dest_dir/$RELEASE_NAME-$KUMA_VERSION

  rm -rf "$dest_dir"
  mkdir "$dest_dir"
  mkdir "$kuma_dir"
  mkdir "$kuma_dir/bin"

  artifact_dir=$(artifact_dir "$arch" "$system")
  cp -p "$artifact_dir/kumactl/kumactl" "$kuma_dir/bin"

  cp tools/releases/templates/LICENSE "$kuma_dir"
  cp tools/releases/templates/NOTICE-kumactl "$kuma_dir"

  archive_path=$(archive_path "$arch" "$system")

  tar -czf "${archive_path}" -C "$dest_dir" .
}

function create_tarball() {
  local arch=$1
  local system=$2
  local distro=$3
  local envoy_distro=$4

  msg ">>> Packaging ${RELEASE_NAME} for $distro ($system-$arch)..."
  msg

  make GOOS="$system" GOARCH="$arch" build

  local dest_dir=build/$RELEASE_NAME-$distro-$arch
  local kuma_dir=$dest_dir/$RELEASE_NAME-$KUMA_VERSION

  rm -rf "$dest_dir"
  mkdir "$dest_dir"
  mkdir "$kuma_dir"
  mkdir "$kuma_dir/bin"
  mkdir "$kuma_dir/conf"

  get_envoy "$distro" "$envoy_distro"
  chmod 755 build/envoy-"$distro"

  artifact_dir=$(artifact_dir "$arch" "$system")
  cp -p "build/envoy-$distro" "$kuma_dir"/bin/envoy
  cp -p "$artifact_dir/kuma-cp/kuma-cp" "$kuma_dir/bin"
  cp -p "$artifact_dir/kuma-dp/kuma-dp" "$kuma_dir/bin"
  cp -p "$artifact_dir/kumactl/kumactl" "$kuma_dir/bin"
  cp -p "$artifact_dir/coredns/coredns" "$kuma_dir/bin"
  cp -p "$artifact_dir/kuma-prometheus-sd/kuma-prometheus-sd" "$kuma_dir/bin"
  cp -p "$KUMA_CONFIG_PATH" "$kuma_dir/conf/kuma-cp.conf.yml"

  cp tools/releases/templates/* "$kuma_dir"

  archive_path=$(archive_path "$arch" "$system" "$distro")

  tar -czf "${archive_path}" -C "$dest_dir" .
}

function package() {
  for os in "${DISTRIBUTIONS[@]}"; do
    local system
    system=$(echo "$os" | awk '{split($0,parts,":"); print parts[1]}')
    local distro
    distro=$(echo "$os" | awk '{split($0,parts,":"); print parts[2]}')
    local envoy_distro
    envoy_distro=$(echo "$os" | awk '{split($0,parts,":"); print parts[3]}')

    for arch in "${GOARCH[@]}"; do

      if [[ -n $distro ]]; then
        create_tarball "$arch" "$system" "$distro" "$envoy_distro"
      else
        create_kumactl_tarball "$arch" "$system"
      fi

      msg
      msg_green "... success!"
      msg
    done
  done
}

function artifact_dir() {
  local arch=$1
  local system=$2

  echo "build/artifacts-$system-$arch"
}

function archive_path() {
  local arch=$1
  local system=$2
  local distro=$3

  if [[ -n $distro ]]; then
    echo "$(artifact_dir "$arch" "$system")/$RELEASE_NAME-$KUMA_VERSION-$distro-$arch.tar.gz"
  else
    echo "$(artifact_dir "$arch" "$system")/$RELEASE_NAME-$CTL_NAME-$KUMA_VERSION-$system-$arch.tar.gz"
  fi
}

function release() {
  for os in "${DISTRIBUTIONS[@]}"; do
    local system
    system=$(echo "$os" | awk '{split($0,parts,":"); print parts[1]}')
    local distro
    distro=$(echo "$os" | awk '{split($0,parts,":"); print parts[2]}')

    for arch in "${GOARCH[@]}"; do
      local artifact
      artifact="$(archive_path "$arch" "$system" "$distro")"
      [ ! -f "$artifact" ] && msg_yellow "Package '$artifact' not found, skipping..." && continue

      if [[ -n $distro ]]; then
        msg_green ">>> Releasing ${RELEASE_NAME} $KUMA_VERSION for $distro ($system-$arch)..."
      else
        msg_green ">>> Releasing ${RELEASE_NAME} $KUMA_VERSION static kumactl binary for $system-$arch..."
      fi

      docker run --rm \
        -e PULP_USERNAME="${PULP_USERNAME}" -e PULP_PASSWORD="${PULP_PASSWORD}" \
        -e PULP_HOST=${PULP_HOST} \
        -v "${PWD}":/files:ro -it kong/release-script \
        --file /files/"$artifact" \
        --package-type ${PULP_PACKAGE_TYPE} --dist-name ${PULP_DIST_NAME} --publish
    done
  done
}

function usage() {
  echo "Usage: $0 [--package|--release]"
  exit 0
}

function main() {
  KUMA_VERSION=$("${SCRIPT_DIR}/version.sh")

  while [[ $# -gt 0 ]]; do
    flag=$1
    case $flag in
    --help)
      usage
      ;;
    --package)
      op="package"
      ;;
    --release)
      op="release"
      ;;
    *)
      usage
      break
      ;;
    esac
    shift
  done

  case $op in
  package)
    package
    ;;
  release)
    [ -z "$PULP_USERNAME" ] && msg_err "PULP_USERNAME required"
    [ -z "$PULP_PASSWORD" ] && msg_err "PULP_PASSWORD required"

    release
    ;;
  esac
}

main "$@"
