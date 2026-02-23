## 🛡️ Go Linter Pre-commit Setup

We use `golangci-lint` to ensure code quality. This hook runs automatically every time you try to commit code. If it finds issues, the commit will be blocked.
``#!/bin/sh

echo "Running golangci-lint on all files..."

# 1. Run the linter on the entire project (./...)
# We use the --config flag to point to your specific path
cd ./Go_Refined_Code/
golangci-lint run --config=linters/.golangci.yml .


# 2. Capture the exit code
PASS=$?

# 3. If the linter finds any issues anywhere in the project, block the commit
if [ $PASS -ne 0 ]; then
    echo "-------------------------------------------------------"
    echo "Linter failed! Please fix the issues above before committing."
    echo "Note: This check scans ALL files, not just your changes."
    echo "-------------------------------------------------------"
    exit 1
fi

exit 0
``
