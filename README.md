OpTISAM
======

__OPTISAM__ (Optimized tool for inventive Software Asset Management) is a tool for the Software Asset Management Compliance Audit and Optimization Tool. This monorepo contains all the backend services namely:

- auth-service
- account-service
- license-service
- import-service

## Quick start
### Download

```
$ git clone https://github.com/Orange-OpenSource/optisam-backend.git
```

### Pre-Build

* Ensure all Go modules are valid, in the root directory run commands:
```
go clean -modcache && go clean -cache
go mod tidy
```

### Build

##### - Change configuration file
<em>Update values of config files **${service}/cmd/server/config-local.toml** as per your requirement</em>

* Building docker images for all micro-services

```
cd ${service-name}/cmd/server
export CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && go build
docker build --pull -t optisam/${service-name}-service:latest -f Dockerfile .
docker push optisam/${service-name}-service:latest
```

* Building docker image for postgres database having required schema for optisam

```
cd account-service\pkg\repository\v1\postgres\scripts
docker build --pull -t optisam/postgres:latest -f Dockerfile .
docker push optisam/postgres:latest
```

### Run

##### - Run using Docker-Compose

you can create and start all the services from your configuration (docker-compose.yml) -

```
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up
```
##### - Play with factory super admin user

1) Once docker-compose is up and running, open optisam dashboard at http://localhost:8090
2) login with below superadmin credentials
    * username - admin@test.com
    * password - admin

<!-- ### Install and Usage
## Contribute
Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.
## Versions  -->

## License

Copyright (c) 2019 Orange

This software is distributed under the terms and conditions of the 'Apache License 2.0'
license which can be found in the file 'License.txt' in this package distribution 
or at 'http://www.apache.org/licenses/LICENSE-2.0'. 

## Contact
* Homepage: [opensource.orange.com](http://opensource.orange.com/)
