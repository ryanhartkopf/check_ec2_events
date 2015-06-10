# check_ec2_events
A short Nagios/Icinga script for checking AWS EC2 maintenance events, written in Go

## Usage
``./check_ec2_events -h
Usage of ./check_ec2_events:
  -name="": Name of the EC2 instance to check for events
  -region="us-east-1": AWS region``

## Credentials
See the "[Getting Started - Using Credentials](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials)" document from the AWS Go SDK documentation for more information about configuring credentials.
