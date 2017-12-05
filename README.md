[![Build Status](https://travis-ci.org/nlindblad/confidential.svg?branch=master)](https://travis-ci.org/nlindblad/confidential)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/c5195a16de6f455986b13a5ff04388d3)](https://www.codacy.com/app/niklas/confidential?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nlindblad/confidential&amp;utm_campaign=Badge_Grade)
# confidential (working title)

Export parameters from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) as environment variables.

See [some examples](#examples) of common use cases.

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
cd $GOPATH/src/github.com/nlindblad/confidential/apps/confidential
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
     exec, e    retrieve environment variables and execute command with an updated environment
     output, o  retrieve and atomically output environment variables to a file
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --forwarded-profile value  AWS profile to forward credentials for in the created environment [$AWS_FORWARDED_PROFILE]
   --prefix value             Amazon SSM parameter prefix
   --profile value            AWS profile to use when calling Amazon SSM [$AWS_PROFILE]
   --region value             AWS region e.g. eu-west-1 [$AWS_REGION]
   --help, -h                 show help
   --version, -v              print the version
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

### :whale2: Use with Docker and systemd services:

A handy way of running Docker containers supervised by systemd is to create a unit (service) using the [`systemd-docker` wrapper](https://github.com/ibuildthecloud/systemd-docker):

```
[Unit]
Description=My service
Requires=docker.service
After=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=/usr/local/bin/confidential --region eu-west-1 --prefix /my-service/prod output --env-file /etc/my-service/prod.env
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

1. Uses confidential to get the latest environment variables from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) in the `eu-west-1` AWS Region and writes them to the file `/etc/my-service/prod.env` in a format that Docker understands

2. Pulls down the latest version of the `username/image-name` Docker image

This ensures that the service is always running using the latest published Docker image and that any configuration changes are picked up automatically.

Managing the environment variables for the service is now done within the `/my-service/prod` namespace in [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html).

### :horse: Use with generic systemd services:

The `EnvironmentFile` directive can be used to expose the retrieved environment variable to any kind of executable running as a systemd service:

```
[Unit]
Description=Service
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=forking
PIDFile=/run/my-service.pid
ExecStartPre=/usr/local/bin/confidential --region eu-west-1 --prefix /my-service/prod output --env-file /etc/my-service/prod.env
EnvironmentFile=-/etc/my-service/prod.env
ExecStart=/usr/local/bin/my-service --flag=something --foo=bar
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/bin/kill -s QUIT $MAINPID
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

Every time the service is started/restarted, it runs the `ExecStartPre` steps and populates `/etc/my-service/prod.env` and includes it using the `EnvironmentFile` directive (*note* the `-` before the filename).

Managing the environment variables for the service is now done within the `/my-service/prod` namespace in [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) in the `eu-west-1` AWS Region.

### :cloud: Give EC2 hosts permissions to access specific parameters:

See full CloudFormation template: [examples/cloudformation/example-3-cloudformation.yml](examples/cloudformation/example-3-cloudformation.yml)

### :telescope: Create an IAM role with permissions to access specific parameters:

See full CloudFormation template: [examples/cloudformation/example-4-cloudformation.yml](examples/cloudformation/example-4-cloudformation.yml)

Creates a dedicated IAM user and access keys that is allowed to decrypt and retrieve parameters with a specific prefix.

*Note*: Some other tools using [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) use a mix of `ssm:DescribeParameters` and `ssm:GetParameters`, which makes it hard to create fine grained acess control, especially when iterating parameters requires permissions to describe **all** parameters.

### :pencil2: Create an IAM role with permissions to set specific parameters:

See full CloudFormation template: [examples/cloudformation/example-5-cloudformation.yml](examples/cloudformation/example-5-cloudformation.yml)

Creates a dedicated IAM user and access keys that is allowed to encrypt and set parameters with a specific prefix, but not retrieve or decrypt.

*Example usage*:

```
aws --profile <PROFILE> ssm put-parameter --name '<PREFIX>/<PARAMETER NAME>' --type "SecureString" --value '<VALUE>'
```


### :shell: Run arbitrary executable with an environment populated by confidential:

The simplest example of this is running the `/usr/bin/env` utility and print out the environment variables that are accessible to the newly invoked process:

```
/usr/local/bin/confidential --region eu-west-1 --prefix /my-service/prod exec -- env
```

### :octopus: Use with supervisord:

```
[program:my-service]
command=/usr/local/bin/confidential --region eu-west-1 --prefix /my-service/prod exec -- /usr/local/bin/my-service --flag=something --foo=bar
directory=/tmp
autostart=true
autorestart=true
startretries=3
stdout_logfile=/tmp/my-service.log
stderr_logfile=/tmp/my-service.err.log
user=username
```

### :card_index: Use specific AWS profile from `~/.aws/credentials`

By default, the AWS SDK for Go will [automatically look for AWS credentials in a couple of pre-defined places](https://github.com/aws/aws-sdk-go#configuring-credentials).

If you are using the standard `~/.aws/credentials` (used by the standard [AWS CLI tool](https://aws.amazon.com/cli/)), you can specify multiple sections with different credentials:

```
[default]
aws_access_key_id = AKIAPEIPJKJSOJ267
aws_secret_access_key = XXXXXXXXXXXXXXXXXXX

[parameters-read]
aws_access_key_id = AKIABCDEFGH12345
aws_secret_access_key = XXXXXXXXXXXXXXXXXXX
```

Using the `--profile` flag, you can specify that you want to use the `parameters-read` profile instead of the `default` one (which would get picked up by the AWS SDK for Go):

```
/usr/local/bin/confidential --profile parameters-read --region eu-west-1 --prefix /my-service/prod output --env-file /etc/my-service/prod.env
```

### :fast_forward: Forward AWS credentials from `~/.aws/credentials` to new environment

It is possible to forward AWS credentials from `~/.aws/credentials` for a given profile to the new enviromment using the `--forwarded-profile` flag.

Given a `~/.aws/credentials` file:

```
[default]
aws_access_key_id = AKIAPEIPJKJSOJ267
aws_secret_access_key = XXXXXXXXXXXXXXXXXXX

[parameters-read]
aws_access_key_id = AKIABCDEFGH12345
aws_secret_access_key = XXXXXXXXXXXXXXXXXXX

[my-service]
aws_access_key_id = AKIAHIHIIW233445
aws_secret_access_key = XXXXXXXXXXXXXXXXXXX
```

You can use the the AWS credentials for the `parameters-read` profile to retrieve the parameters from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) and forward the AWS credentials for the `my-service` profile using:

```
/usr/local/bin/confidential --profile parameters-read --forwarded-profile my-service --region eu-west-1 --prefix /my-service/prod output --env-file /etc/my-service/prod.env
```

In the above example, `/etc/my-service/prod.env` would contain all parameters retrieved from [AWS Systems Manager Parameters](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html) in the `eu-west-1` AWS Region in addition to:

```
AWS_ACCESS_KEY_ID=AKIAHIHIIW233445
AWS_SECRET_ACCESS_KEY=XXXXXXXXXXXXXXXXXXX
AWS_SESSION_TOKEN=
```

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
* [gurusi/systemd-make-environment](https://github.com/gurusi/systemd-make-environment) for mentioning some common gotchas with `EnvironmentFile` and `ExecStartPre` with `systemd`
* [Storing Secrets with AWS ParameterStore](https://typicalrunt.me/2017/04/07/storing-secrets-with-aws-parameterstore/) provided a great starting point for the CloudFormation templates
* [segmentio/chamber](https://github.com/segmentio/chamber) for a nice way of implementing the `exec` command
