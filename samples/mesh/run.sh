# !/bin/bash
set -x 

cmd=$1
case $cmd in
"r"|"run")
    nohup ./main -t model -p 8001 > ./log/model.1.log 2>&1  &
    nohup ./main -t model -p 8002 > ./log/model.2.log 2>&1  &
    nohup ./main -t entry -p 8026 > ./log/entry.log 2>&1  &
    nohup ./main -t middle -p 8015 > ./log/middle.log 2>&1 &
    ;;
"k"|"kill")
    eval "ps -ef | grep main | awk '{print(\$2)}' | xargs -I {} kill {}"
    ;;
"vet")
    GOPATH="" go vet ./...
    ;;
"lint")
    GOPATH="" golint ./...
    ;;
*)
    echo "invalid arg: $cmd"
    ;;
esac