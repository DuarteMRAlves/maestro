#!/bin/sh -l

cd "${GITHUB_WORKSPACE}" || exit 1

GOFMT_OUTPUT="$(/usr/local/go/bin/gofmt -l -s "${1}")"

if [ -n "${GOFMT_OUTPUT}" ]; then
  echo "gofmt errors in:"
  echo "${GOFMT_OUTPUT}"

  echo "::set-output name=report::gofmt errors in at least on file"
  exit 1
fi

echo "::set-output name=report::No errors found"