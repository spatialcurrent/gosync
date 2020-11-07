#!/bin/bash

# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

testHelp() {
  "${DIR}/../bin/gosync" --help
}

testLocalToLocal() {
  "${DIR}/../bin/gosync" --parents testdata "${SHUNIT_TMPDIR}/testdata"
  assertEquals "unexpected output for tmp/doc.1.txt" 'hello' "$(cat "${SHUNIT_TMPDIR}/testdata/doc.1.txt")"
  assertEquals "unexpected output for tmp/a/doc.1.txt" 'hello' "$(cat "${SHUNIT_TMPDIR}/testdata/a/doc.1.txt")"
  assertEquals "unexpected output for tmp/a/doc.2.txt" 'world' "$(cat "${SHUNIT_TMPDIR}/testdata/a/doc.2.txt")"
  assertEquals "unexpected output for tmp/a/b/doc.1.txt" 'hello' "$(cat "${SHUNIT_TMPDIR}/testdata/a/b/doc.1.txt")"
  assertEquals "unexpected output for tmp/a/b/doc.2.txt" 'world' "$(cat "${SHUNIT_TMPDIR}/testdata/a/b/doc.2.txt")"
}

oneTimeSetUp() {
  echo "Using temporary directory at ${SHUNIT_TMPDIR}"
}

oneTimeTearDown() {
  echo "Tearing Down"
}

# Load shUnit2.
. "${DIR}/shunit2"
