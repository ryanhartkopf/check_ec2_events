package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

var instanceNamePtr = flag.String("name", "", "Name of the EC2 instance to check for events")
var awsRegionPtr = flag.String("region", "us-east-1", "AWS region")

func getInstanceId(instanceName string, svc *ec2.EC2) *string {
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{&ec2.Filter{
			Name:   aws.String("tag-value"),
			Values: []*string{aws.String(instanceName)},
		}},
	})
	if err != nil {
		fmt.Println("WARNING - check AWS credentials and try again")
		os.Exit(1)
	}
	if resp.Reservations != nil {
		return resp.Reservations[0].Instances[0].InstanceID
	} else {
		return nil
	}
}

func getInstanceStatus(instanceId *string, svc *ec2.EC2) *ec2.InstanceStatusEvent {
	resp, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIDs: []*string{aws.String(*instanceId)},
	})
	if err != nil {
		panic(fmt.Sprintln("Error retrieving instance status"))
	}
	if resp.InstanceStatuses[0].Events != nil {
		return resp.InstanceStatuses[0].Events[0]
	} else {
		return nil
	}
}

func main() {

	flag.Parse()
	svc := ec2.New(&aws.Config{Region: *awsRegionPtr})

	if *instanceNamePtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	instanceId := getInstanceId(*instanceNamePtr, svc)

	if instanceId == nil {
		fmt.Println("WARNING - no instance was found with instance name", *instanceNamePtr, "in region", *awsRegionPtr)
		os.Exit(1)
	}

	instanceStatus := getInstanceStatus(instanceId, svc)
	if instanceStatus != nil {
		switch code := instanceStatus.Code; *code {
		case "instance-reboot":
			fmt.Println("WARNING - instance", *instanceNamePtr, "will be rebooted at", instanceStatus.NotBefore, "- Description:", *instanceStatus.Description)
			os.Exit(1)
		case "system-reboot":
			fmt.Println("WARNING - physical host will be rebooted at", *instanceStatus.NotBefore, "- Description:", *instanceStatus.Description)
			os.Exit(1)
		case "system-maintenance":
			fmt.Println("WARNING - physical host maintenance will be performed at", instanceStatus.NotBefore, "- Description:", *instanceStatus.Description)
			os.Exit(1)
		case "instance-retirement", "instance-stop":
			fmt.Println("CRITICAL - instance", *instanceNamePtr, "will be retired at", instanceStatus.NotBefore, "- Description:", *instanceStatus.Description)
			os.Exit(2)
		default:
			fmt.Println("CRITICAL - unknown event type", instanceStatus.Code, "scheduled for", instanceStatus.NotBefore, "- Description:", *instanceStatus.Description)
			os.Exit(2)
		}
	} else {
		fmt.Println("OK - no events for instance", *instanceNamePtr)
		os.Exit(0)
	}
}
