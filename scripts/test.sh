#!/bin/bash
# source ./scripts/lib.trap.sh

set -e
set -x

TEST_OUTPUT=brnchswppr_test.out

cd ./tests

# test
echo "Test: list with empty file"
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 0 ]; then
    echo "Error: expected 0 branches in list"
    exit 1
fi

echo "Test: swap with empty file"
brnchswppr first
brnchswppr -ls > $TEST_OUTPUT
if ! [ "$(cat $TEST_OUTPUT | wc -l)" -eq 1 ]; then
    echo "Error: expected 1 branch in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == "0: main" ]; then
    echo "Error: expected main branch in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "first" ]; then
    echo "Error: expected first branch to be checked out"
    exit 1
fi

echo "Test: swap with branches in file"
brnchswppr second
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 2 ]; then
    echo "Error: expected 2 branches in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == $'0: main\n1: first' ]; then
    echo "Error: expected main and first branches in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "second" ]; then
    echo "Error: expected second branch to be checked out"
    exit 1
fi

echo "Test: swap the current branch"
# add swap the current branch with empty file
brnchswppr second
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 2 ]; then
    echo "Error: expected 2 branches in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == $'0: main\n1: first' ]; then
    echo "Error: expected main and first branches in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "second" ]; then
    echo "Error: expected second branch to be checked out"
    exit 1
fi

echo "Test: skip swap, stash current"
brnchswppr
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 3 ]; then
    echo "Error: expected 3 branches in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == $'0: main\n1: first\n2: second' ]; then
    echo "Error: expected main, first, and second branches in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "second" ]; then
    echo "Error: expected second branch to be checked out"
    exit 1
fi

echo "Test: skip swap, skip stash"
brnchswppr
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 3 ]; then
    echo "Error: expected 3 branches in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == $'0: main\n1: first\n2: second' ]; then
    echo "Error: expected main, first, and second branches in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "second" ]; then
    echo "Error: expected second branch to be checked out"
    exit 1
fi

echo "Test: swap to existing branch and unstash it"
brnchswppr first
brnchswppr -ls > $TEST_OUTPUT
if ! [ $(cat $TEST_OUTPUT | wc -l) -eq 2 ]; then
    echo "Error: expected 2 branches in list"
    exit 1
fi
if ! [ "$(cat $TEST_OUTPUT)" == $'0: main\n1: second' ]; then
    echo "Error: expected main and second branches in list"
    exit 1
fi
if ! [ "$(git rev-parse --abbrev-ref HEAD)" == "first" ]; then
    echo "Error: expected first branch to be checked out"
    exit 1
fi

exit 0
