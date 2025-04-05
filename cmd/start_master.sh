#!/bin/bash

PROJECT_ROOT=/home/yippee/Documents/fedoraWorkspace/go-imdg
PROJECT_LOG=$PROJECT_ROOT/var/log

nohup go run ./cmd -config="./.config/master.json" > $PROJECT_LOG/master_script.log 2>&1 &

