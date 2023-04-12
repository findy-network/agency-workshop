#!/bin/bash

set -e

export CRED_DEF_ID=${1:-"$CRED_DEF_ID"}

function createCredDef {
  if [ -z "$CRED_DEF_ID" ]; then
    echo "Create schema"
    sch_id=$(cli agent create-schema \
      --name="foobar" \
      --version=1.0 foo)

    # read schema - make sure it's found in ledger
    echo "Read schema $sch_id"
    schema=$(cli agent get-schema --schema-id $sch_id)

    # create cred def
    echo "*** Create cred def with schema id $sch_id"
    CRED_DEF_ID=$(cli agent create-cred-def \
      --id $sch_id --tag "TAG")

    # read cred def - make sure it's found in ledger
    echo "Read cred def $CRED_DEF_ID"
    cred_def=$(cli agent get-cred-def --id $CRED_DEF_ID)

    export CRED_DEF_ID="$CRED_DEF_ID"
  fi
}

current_dir=$(dirname "$BASH_SOURCE")

createCredDef

printf "\n\nHi there 👋 \n"
printf "\nIssuer bot started 🤖\n"

cli bot start $current_dir/issue-bot.yaml

printf "\n\nHi there 👋 \n"
printf "\nVerify bot started 🤖\n"

cli bot start $current_dir/verify-bot.yaml

printf "\n\nHi there 👋 All done!\n"
