# check_ec2_events
A short Nagios/Icinga script for checking AWS EC2 maintenance events, written in Go

## Usage
```
$ ./check_ec2_events -h
Usage of ./check_ec2_events:
  -name="": Name of the EC2 instance to check for events
  -region="us-east-1": AWS region
```

## Credentials
See the "[Getting Started - Using Credentials](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials)" document from the AWS Go SDK documentation for more information about configuring credentials.

## Example Output
```
$ ./check_ec2_events -name="web1.example.com"
CRITICAL - instance is scheduled to be retired on 2015-06-10 18:30:00 UTC. Description: The instance will be terminated.
```
```
$ ./check_ec2_events -name="web2.example.com"
OK - no events for instance web2.example.com
```
