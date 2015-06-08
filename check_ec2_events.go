package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	instanceIdResp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{&ec2.Filter{
			Name:   aws.String("tag-value"),
			Values: []*string{aws.String("at-drupal-dev.atrust.com")},
		}},
	})
	instanceId := *instanceIdResp.Reservations[0].Instances[0].InstanceID
	fmt.Println(instanceId, err)

	instanceStatusResp, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIDs: []*string{aws.String(instanceId)},
	})

	instanceStatus := instanceStatusResp.InstanceStatuses[0].Events[0]

	fmt.Println(instanceStatus.Description, err)
	fmt.Println(*instanceStatus.Description, err)
}
