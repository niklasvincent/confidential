[![Build Status](https://travis-ci.org/nlindblad/confidential.svg?branch=master)](https://travis-ci.org/nlindblad/confidential)
# confidential (working title)

Export parameters from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) as environment variables.

## Why I wrote this?

Configuration management, specifically secrets management tends to get complicated. After having been through several projects, both in my spare time and at work, using solutions such as [Ansible Vault](https://docs.ansible.com/ansible/2.4/vault.html), private [AWS CodeCommit](https://aws.amazon.com/codecommit/) repositories or [Amazon KMS](https://aws.amazon.com/kms/) encrypted configuration files in [Amazon S3](https://aws.amazon.com/s3/), I was looking for something simpler, while still maintaining a high level of security.

I deemed self-hosted solution, such as [Hashicorp Vault](https://www.vaultproject.io/) (and the other solutions listed on [this Hashicorp Vault vs. Other Software](https://www.vaultproject.io/intro/vs/index.html) page) too time consuming to set up and maintain.

Luckily, Amazon Web Services have [been busy improving their Amazon EC2 Systems Manager Parameter Store](https://aws.amazon.com/blogs/mt/amazon-ec2-systems-manager-parameter-store-adds-support-for-parameter-versions/) in 2017 and it now supports both seamless [Amazon KMS](https://aws.amazon.com/kms/) encryption and versioning of parameters.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Confidential is written in [Go](https://golang.org) and can be run as a single binary, no language specific requirements are needed.

It is designed to run on the GNU/Linux, macOS, and Windows operating systems. Other operating systems will probably work as long as you can compile a Go binary on them.

### Installing

Make sure you have [Go installed](https://golang.org/doc/install) and that [the `$GOPATH` is set correctly](https://github.com/golang/go/wiki/SettingGOPATH).

### Build binary

```
go get github.com/nlindblad/confidential
cd $GOPATH/src/github.com/nlindblad/confidential
go build
```

### Run

```
./confidential --help
```

And you should see:

```
NAME:
   confidential - A new cli application

USAGE:
   confidential [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR:
   Niklas Lindblad <niklas@lindblad.info>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prefix value    parameter prefix
   --env-file value  output environment file
   --help, -h        show help
   --version, -v     print the version
```

## Running the tests

Tests use the excellent Go [stretchr/testify](https://github.com/stretchr/testify) package.

```
go test -v ./...
```

## Deployment

TODO

## Examples

The machine needs to have the following AWS IAM permissions:

- `kms:Decrypt` on the relevant [Amazon KMS](https://aws.amazon.com/kms/) key used to encrypt sensitive parameters.
- `ssm:GetParametersByPath` on the relevant resource: `arn:aws:ssm:${AWS::Region}:${AWS::AccountId}:parameter/<PREFIX>` (**note** there should be no trailing slash or wildcards)

### Example 1) Use with Docker and systemd services:

A handy way of running Docker containers supervised by systemd is to create a unit (service) using the [`systemd-docker` wrapper](https://github.com/ibuildthecloud/systemd-docker):

```
[Unit]
Description=My service
Requires=docker.service
After=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=/usr/local/bin/confidential --prefix /my-service/prod --env-file /etc/my-service/prod.env
ExecStartPre=/usr/bin/docker pull username/image-name:latest
ExecStart=/usr/local/bin/systemd-docker --cgroups name=systemd run \
    --name %n \
    --env-file /etc/my-service/prod.env \
    ... Add other Docker run flags here ...
    -d username/image-name:latest
ExecStop=/usr/bin/docker stop %n
ExecStopPost=/usr/bin/docker rm -f %n
Restart=always
RestartSec=10s
Type=notify
NotifyAccess=all

[Install]
WantedBy=default.target
```

The following service will run `username/image-name` as a service which will get restarted if it falls over.

Every time the service is started/restarted, it runs the two `ExecStartPre` steps:

1. Uses confidential to get the latest environment variables from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) and writes them to the file `/etc/my-service/prod.env` in a format that Docker understands

2. Pulls down the latest version of the `username/image-name` Docker image

This ensures that the service is always running using the latest published Docker image and that any configuration changes are picked up automatically.

Managing the environment variables for the service is now done within the `/my-service/prod` namespace in [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html).

### Example 2) Use with generic systemd services:

The `EnvironmentFile` directive can be used to expose the retrieved environment variable to any kind of executable running as a systemd service:

```
[Unit]
Description=Service
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=forking
PIDFile=/run/my-service.pid
ExecStartPre=/usr/local/bin/confidential --prefix /my-service/prod --env-file /etc/my-service/prod.env
EnvironmentFile=/etc/my-service/prod.env
ExecStart=/usr/local/bin/my-service
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/bin/kill -s QUIT $MAINPID
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Every time the service is started/restarted, it runs the `ExecStartPre` steps and populates `/etc/my-service/prod.env` and includes it using the `EnvironmentFile` directive.

Managing the environment variables for the service is now done within the `/my-service/prod` namespace in [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html).

### Example 3) Give EC2 hosts permissions to access specific parameters:

Based on "*[Storing Secrets with AWS ParameterStore](https://typicalrunt.me/2017/04/07/storing-secrets-with-aws-parameterstore/)*":

See full Cloudformation template: [examples/cloudformation/example-3-cloudformation.yml](examples/cloudformation/example-3-cloudformation.yml)

### Example 4) Create an IAM role with permissions to access specific parameters:

See full Cloudformation template: [examples/cloudformation/example-4-cloudformation.yml](examples/cloudformation/example-4-cloudformation.yml)

Creates a dedicated IAM user and access keys that is allowed to decrypt and retrieve parameters with a specific prefix.

### Example 5) Create an IAM role with permissions to set specific parameters:

See full Cloudformation template: [examples/cloudformation/example-5-cloudformation.yml](examples/cloudformation/example-5-cloudformation.yml)

Creates a dedicated IAM user and access keys that is allowed to encrypt and set parameters with a specific prefix, but not retrieve or decrypt.

*Example usage*:

```
aws --profile <PROFILE> ssm put-parameter --name '<PREFIX>/<PARAMETER NAME>' --type "SecureString" --value '<VALUE>'
```


### Example 6) Run arbitrary executable with an environment populated by confidential:

TODO + implement

### Example 7) Use with supervisord:

TODO + implement

## Built With

* [aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) - AWS SDK for Go
* [urfave/cli](https://github.com/urfave/cli) - Command line interface helpers
* [dchest/safefile](https://github.com/dchest/safefile) - Safe "atomic" saving of files

## Versioning

Uses [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nlindblad/confidential/tags). 

## Authors

* **Niklas Lindblad** - *Initial work*

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* [Sjeanpierre/param_api](https://github.com/Sjeanpierre/param_api) provided a great starting point for using the Amazon SSM API in Go
