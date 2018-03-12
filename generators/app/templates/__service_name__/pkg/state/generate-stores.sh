#!/bin/bash
set -e

(

# Add any required imports here, separated by commas
CUSTOM_IMPORTS=""

LOCAL=$(dirname $PWD)
echo $LOCAL

GENERATOR_PATH="$LOCAL/../cmd/igenerator"
ROOT=$LOCAL
GEN="go run $GENERATOR_PATH/main.go"
TEMPLATE_DIR="$GENERATOR_PATH"

echo "Generating metrics store"
$GEN -target=Store \
     -input="$ROOT/state/store.go" \
     -ignore="Initialize" \
     -template="$TEMPLATE_DIR/metrics_store.go.tpl" \
     -output="$ROOT/state/metrics_store.go" \
     -imports="$CUSTOM_IMPORTS"

echo "Generating trace store"
$GEN -target=Store \
     -input="$ROOT/state/store.go" \
     -ignore="Initialize" \
     -template="$TEMPLATE_DIR/trace_store.go.tpl" \
     -output="$ROOT/state/trace_store.go" \
     -imports="$CUSTOM_IMPORTS"
     
)
