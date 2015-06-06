tail -n +2 DIRS | xargs -L1 -i sh -c "cd {} && go test -i && go test -v"
