#!/bin/sh

COVERAGE_REPORT="$1"
test -r "$COVERAGE_REPORT" || exit 1
MODULE_PATH="$(go list -m)"

if [ -n "$CODECOV_TOKEN" ]; then
  echo "Upload Codecov Coverage"
  codecov -f "$COVERAGE_REPORT" &
fi

if [ -n "$CC_TEST_REPORTER_ID" ]; then
  echo "Upload Code Climate Coverage"
  cc-test-reporter format-coverage -t gocov -p "$MODULE_PATH" -o .coverage/codeclimate.json "$COVERAGE_REPORT"
  cc-test-reporter upload-coverage -r "$CC_TEST_REPORTER_ID" -i .coverage/codeclimate.json &
fi

wait || true
echo "Coverage Upload Done"
