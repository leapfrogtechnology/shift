# Shift

![GitHub release](https://img.shields.io/github/release/leapfrogtechnology/shift.svg?style=flat)
![CircleCI](https://img.shields.io/circleci/build/github/leapfrogtechnology/shift/master?style=flat)
![GitHub](https://img.shields.io/github/license/leapfrogtechnology/shift.svg?style=flat)

Automated infrastructure creation and deployment tool.

## About

Shift focuses on creating your infrastructure with ease and automating your deployments.

## Usage

### 1. Install Shift:

```
$ curl -sf https://raw.githubusercontent.com/leapfrogtechnology/shift/master/install.sh | sudo sh
```

### 2. Verify Installation:

```
$ shift-cli
```

### 3. AWS Credentials

Before using shift, you have to provide your AWS credentials. Also, make sure you have the correct access for your credentials.

Simply create `~/.aws/credentials` file for storing AWS credentials. <br/>
Example:

```
[example-profile]
aws_access_key_id = xxxxxx
aws_secret_access_key = xxxxxx
```

To know more about AWS configurations, view [Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)

### 4. Setup

Go to your project directory and run `setup` command to initiate the setup process.

```
$ shift-cli setup
```

- Enter the name of your project
- Choose the platform
- Choose your AWS profile that was setup earlier.
- Choose the AWS Region where you want your infrastructure to be created.
- Choose the type of deployment you want.
- Name your environment.
- If Backend:
  - Enter your application port
  - Enter healthcheck path
  - Enter dockerfile path
- Enter Slack webhook URL
-

```
 Example:

 Project Name: shiftbeta
 Choose Cloud Provider: AWS
 Choose AWS profile: default
 Region: US East (North Virginia)
 Choose Deployment Type: Backend
 Environment Name: dev
 Application port: 8080
 Heathcheck Path: /api
 Dockerfile Path: ./
 Slack Webhook URL: https://webhook.slack.com
```

`shift.json` file will be created in your project directory like below.

```json
{
  "name": "shiftbeta",
  "platform": "AWS",
  "profile": "default",
  "region": "us-east-1",
  "type": "Backend",
  "port": "80",
  "dockerFilePath": "./",
  "healthCheckPath": "/",
  "env": {
    "dev": {
      "cluster": "shiftbeta-dev",
      "service": "shiftbeta-dev",
      "image": "009409476372.dkr.ecr.us-east-1.amazonaws.com/shiftbeta/dev-backend",
      "taskDefinition": "shiftbeta-dev"
    }
  }
}
```

**If you do not want a project-specific config file, you can skip the above step.**

### 5. Deploy

```
$ shift-cli deploy [env]
$ shift-cli deploy dev
```

Here `dev` is the environment you specified in `shif.json`.

### 6. Destroy

```
$ shift-cli destroy [env]
```

Destroy infrastracture.

### 7. Usage with CI/CD:

Instead of setting up a `~/.aws/credentials` file. You can also use the following environment variables to set up your AWS credentials.

| Variable              | Description                            |
| --------------------- | -------------------------------------- |
| AWS_ACCESS_KEY_ID     | Your AWS access key                    |
| AWS_SECRET_ACCESS_KEY | Your AWS secret key                    |
| AWS_REGION            | AWS region where you added your secret |
