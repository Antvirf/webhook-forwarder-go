# Webhook forwarder API in Go

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Antvirf_webhook-forwarder-go&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Antvirf_webhook-forwarder-go)
[![CodeQL](https://github.com/Antvirf/webhook-forwarder-go/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Antvirf/webhook-forwarder-go/actions/workflows/codeql-analysis.yml)

Very much a work in progress. Replicating in Go the webhook forwarder I built in Python as a way to better understand the language.

## To-do

* ~~Forward webhook function: headers and body~~
* ~~Fetch acceptable IPs from GitHub meta api~~
* ~~Proper 'is in' logic for IPs - GitHub returns CIDR ranges~~
* ~~Parse IP fields as IPs for safety~~
* ~~Env variables: TARGET_URL~~
* ~~Env variables: WEBHOOK_TOKEN_SECRET~~
* Tests
* Coverage analysis for tests
* Go 'maintainability' analysis
