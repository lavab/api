d=$(dirname $0)
xargs -L1 -i sh -c "go build $d/{}/*.go" < DIRS
