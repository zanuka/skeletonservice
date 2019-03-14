# Go kit skeleton Micro Service
This [GoKit Microservice](https://gokit.io/) skeleton designed to be copied and used as a starting point Google AppEngine based MicroServices. However it is setup to be easily be adapted to run anywhere by modifying only skelconfig/config.go 

Architecture was taken and adapted from [Anton Klimenko's](https://github.com/antklim) Medium Article Microservices in Go: https://medium.com/seek-blog/microservices-in-go-2fc1570f6800

Please feel free to use and remix.

## Skeleton implementation details
- http transport layer
- grpc transport layer
- separation of concerns
- easy config pattern
- HealthCheck implementation example
- Login implementation example
- travis, test and deploy support
- deploy scripts
- App Engine app.yaml with evn vars on deploy

## Running with Google AppEngine - Flexible env

### Perquisites
 - gcloud cli installed
 - service account - json file
 - kms key ring

### Local Dev
Running the app in its current form a google service-account an encrypted config.yml file. 

You will need to set the evn var - this will be used to decrypt the config.yml at runtime and expose it to the app
`GOOGLE_APPLICATION_CREDENTIALS=gcloud-service-account.json`


### How to Encrypting config.yml
```
gcloud kms encrypt --location=global --keyring=[key-ring-name] --key=l[key-name] --plaintext-file=config.yaml --ciphertext-file=local-config.yaml.enc
```

## Running with AWS or other platforms
The skeleton  implementation is agnostics apart from the skelconfig.config.go file. Changing this or deleting most of its content will remove the platform deps.

Or conversely implement your own config that allows you to easily access your platforms specific features.

