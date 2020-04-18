#/bin/bash

set -e

case $(uname | tr '[:upper:]' '[:lower:]') in
  linux*)
    echo "linux is currently unsupported"
    ;;
  darwin*)
    OS_NAME=osx
    ./scripts/install-app-osx.sh
    ;;
  msys*)
    echo "windows is currently unsupported"
    ;;
  *)
    echo "os is unsupported"
    ;;
esac

