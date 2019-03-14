#!/usr/bin/env bash

cd scripts/deploy/prod/
openssl aes-256-cbc -K $encrypted_xxxxxxx_key -iv $encrypted_xxxxxxx_iv -in gcloud-prod-key.json.enc -out gcloud-prod-key.json -d
gcloud auth activate-service-account --key-file gcloud-prod-key.json
cd ../../../cmd
cat ../scripts/deploy/prod/env.yaml >> app.yaml
cat app.yaml
gcloud -q app deploy