# logspout-cloudwatch

A [Logspout](https://github.com/gliderlabs/logspout) adapter for writing Docker container logs to AWS CloudWatch.

## Usage

Route logs from all containers to CloudWatch by starting logspout with the following options:

```
docker run --name=logspout --hostname=$HOSTNAME \
    -e AWS_REGION=us-west-2 \
    -v /var/run/docker.sock:/tmp/docker.sock
    my-logspout-container cloudwatch://my-log-group
```

Once started, this container will create a new Log Stream (`$HOSTNAME`) within the specified Log Group (`my-log-group`) and begin streaming Docker container logs into CloudWatch.

This example depends on a custom Logspout container (`my-logspout-container`) built with the logspout-cloudwatch module and an existing CloudWatch Log Group (`my-log-group`). See [gliderlabs/logspout#modules](https://github.com/gliderlabs/logspout#modules) for more information on building custom logspout containers.

## Configuration

logspout-cloudwatch accepts a number of environment variables that can be used to customize behavior.

### AWS_REGION

Determines the AWS region to which logs will be sent. This option is required.

### LOG_LEVEL

Determines the log level used for logspout-cloudwatch logs. This option defaults to `INFO`, and logspout-cloudwatch will only log startup information, information about failed uploads, and information about rejected events. Set this option to `DEBUG` for detailed information about each uploaded log batch.
