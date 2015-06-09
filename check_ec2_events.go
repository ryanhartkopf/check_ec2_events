package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kr/pretty"
	"os"
)

var svc ec2.EC2

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

	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	instanceId := getInstanceId("test", svc)

	if instanceId == nil {
		fmt.Println("WARNING - no instance was found with that ID")
		os.Exit(2)
	}

	instanceStatus := getInstanceStatus(instanceId, svc)
	if instanceStatus != nil {
		fmt.Printf("DEBUG - Instance status output:\n%# v\n", pretty.Formatter(instanceStatus))

		switch code := instanceStatus.Code; *code {
		case "instance-reboot":
			fmt.Println("WARNING - instance reboot scheduled for", instanceStatus.NotBefore, ". Description:", instanceStatus.Description)
			os.Exit(1)
		case "system-reboot":
			fmt.Println("WARNING - host will be rebooted at ", instanceStatus.NotBefore, ". Description:", instanceStatus.Description)
			os.Exit(1)
		case "system-maintenance":
			fmt.Println("WARNING - host maintenance scheduled for ", instanceStatus.NotBefore, ". Description:", instanceStatus.Description)
			os.Exit(1)
		case "instance-retirement", "instance-stop":
			fmt.Println("CRITICAL - instance is scheduled to be retired on ", instanceStatus.NotBefore, ". Description:", instanceStatus.Description)
			os.Exit(2)
		}
	} else {
		fmt.Println("OK - no events for instance")
		os.Exit(0)
	}

}
