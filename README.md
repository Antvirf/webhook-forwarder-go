# Webhook forwarder API in Go

Very much a work in progress. Replicating in Go the webhook forwarder I built in Python as a way to better understand the language.

## To-do

* ~~Forward webhook function: headers and body~~
* ~~Fetch acceptable IPs from GitHub meta api~~
* ~~Proper 'is in' logic for IPs - GitHub returns CIDR ranges~~
* ~~Parse IP fields as IPs for safety~~
* Env variables: TARGET_URL
* Env variables: WEBHOOK_TOKEN_SECRET

## Other

* Tests
* Coverage analysis for tests
* Go 'maintainability' analysis
