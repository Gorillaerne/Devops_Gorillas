#!/bin/bash

PYTHON_SCRIPT_PATH=$1

while true
do
    python "$PYTHON_SCRIPT_PATH"
    if [ "$EXIT_CODE" -ne 0 ]; then
        echo "Script crashed with exit code $EXIT_CODE. Restarting..." >&2
        sleep 1
    fi
done