mv .git/hooks/pre-commit.sample .git/hooks/pre-commit
nano .git/hooks/pre-commit

#!/bin/sh
make fmt && git add . && git rm -rf --cached . > /dev/null && git add --all
