MODIFIED_FILES=$(git diff --name-status HEAD | grep -E '^(A|M|R)' | awk '{ print $NF }' | grep '\.go$' | xargs -n1 dirname | sort -u)

exit_code=0

for FILE in $MODIFIED_FILES; do
    echo "# Checking $FILE"

    cover_dir=.coverage/$FILE
    cover_file=$coverage_dir/coverage.out
    mkdir -p "$cover_dir"
    timeout=60s

    echo "Running tests..."
    go test -gcflags=all=-d=checkptr -cover -coverprofile=$cover_file -timeout=$timeout -short "./$FILE"

    if [ $? -ne 0 ]; then
        exit_code=1
    fi
    
    echo "Working on $FILE"

    coverage_code=0
    output=$(go tool cover -func="$cover_file")
    while read line; do
        if [[ $(echo "$line" | grep -c -E '\b([0-7][0-9]|[0-9])\.[0-9]+%') -gt 0 ]] && [[ $(echo "$line" | grep -cE '\s(init|main|\(statements\))\s') -eq 0 ]] && [[ $(echo "$line" | grep -cE '\s\S+SC\s') -eq 0 ]]; then
            echo "$line"
            coverage_code=1
        else
            echo "$line"
        fi
    done <<<"$output"

    if [[ $coverage_code -eq 1 ]]; then
        exit_code=1
    fi
done

exit $exit_code
