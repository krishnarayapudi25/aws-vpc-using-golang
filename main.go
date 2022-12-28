package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"github.com/aws/aws-sdk-go/service/eks"
)

func main() {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	// Create a new EC2 client
	svc := ec2.New(sess)

	// Define the VPC CIDR block
	vpcCIDR := "192.168.0.0/16"

	// Define the public and private subnet CIDR blocks
	publicSubnet01CIDR := "192.168.0.0/18"
	publicSubnet02CIDR := "192.168.64.0/18"
	privateSubnet01CIDR := "192.168.128.0/18"
	privateSubnet02CIDR := "192.168.192.0/18"

	// Create the VPC
	vpcInput := &ec2.CreateVpcInput{
		CidrBlock:                   aws.String(vpcCIDR),
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
	}
	vpcOutput, err := svc.CreateVpc(vpcInput)
	if err != nil {
		fmt.Println("Error creating VPC:", err)
		return
	}
	vpcID := *vpcOutput.Vpc.VpcId

	// Create internet gateway
	igInput := &ec2.CreateInternetGatewayInput{}
	igOutput, err := svc.CreateInternetGateway(igInput)
	if err != nil {
		fmt.Println("Error creating Internet Gateway:", err)
		return
	}
	igID := *igOutput.InternetGateway.InternetGatewayId

	// Attach the Internet Gateway to the VPC
	_, err = svc.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igID),
		VpcId:             aws.String(vpcID),
	})
	if err != nil {
		fmt.Println("Error attaching Internet Gateway:", err)
	}

	// Creating EIP1
	Elasticipinput := &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	}

	Eip1, err := svc.AllocateAddress(Elasticipinput)

	if err != nil {
		fmt.Println("error creating eip1", err)
	}

	Elasticip1 := Eip1.AllocationId

	// Creating EIP2
	Elasticipinput2 := &ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	}

	Eip2, err := svc.AllocateAddress(Elasticipinput2)

	if err != nil {
		fmt.Println("error creating eip1", err)
	}

	Elasticip2 := Eip2.AllocationId

	// Create the public subnets
	publicSubnet01Input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(publicSubnet01CIDR),
		VpcId:     aws.String(vpcID),
	}
	publicSubnet01Output, err := svc.CreateSubnet(publicSubnet01Input)
	if err != nil {
		fmt.Println("Error creating public subnet 01:", err)
		return
	}
	publicSubnet01ID := publicSubnet01Output.Subnet.SubnetId

	publicSubnet02Input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(publicSubnet02CIDR),
		VpcId:     aws.String(vpcID),
	}
	publicSubnet02Output, err := svc.CreateSubnet(publicSubnet02Input)
	if err != nil {
		fmt.Println("Error creating public subnet 02:", err)
		return
	}
	publicSubnet02ID := publicSubnet02Output.Subnet.SubnetId

	// Create the private subnets
	privateSubnet01Input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(privateSubnet01CIDR),
		VpcId:     aws.String(vpcID),
	}
	privateSubnet01Output, err := svc.CreateSubnet(privateSubnet01Input)
	if err != nil {
		fmt.Println("Error creating private subnet 01:", err)
		return
	}
	privateSubnet01ID := privateSubnet01Output.Subnet.SubnetId

	privateSubnet02Input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(privateSubnet02CIDR),
		VpcId:     aws.String(vpcID),
	}
	privateSubnet02Output, err := svc.CreateSubnet(privateSubnet02Input)
	if err != nil {
		fmt.Println("Error creating private subnet 02:", err)
		return
	}

	privateSubnet02ID := privateSubnet02Output.Subnet.SubnetId

	// Create the public route table
	publicRouteTableInput := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcID),
	}

	publicRouteTableOutput, err := svc.CreateRouteTable(publicRouteTableInput)
	if err != nil {
		fmt.Println("Error creating public route table:", err)
		return
	}

	publicRouteTableID := *publicRouteTableOutput.RouteTable.RouteTableId

	// Create the private route tables
	// 1
	privateRouteTable01Input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcID),
	}
	privateRouteTable01Output, err := svc.CreateRouteTable(privateRouteTable01Input)
	if err != nil {
		fmt.Println("Error creating private route table 01:", err)
		return
	}
	privateRouteTable01ID := *privateRouteTable01Output.RouteTable.RouteTableId

	//  2
	privateRouteTable02Input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcID),
	}
	privateRouteTable02Output, err := svc.CreateRouteTable(privateRouteTable02Input)
	if err != nil {
		fmt.Println("Error creating private route table 02:", err)
		return
	}
	privateRouteTable02ID := *privateRouteTable02Output.RouteTable.RouteTableId

	// Create the public route to the Internet Gateway
	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		RouteTableId:         aws.String(publicRouteTableID),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igID),
	})
	if err != nil {
		fmt.Println("Error creating public route:", err)
		return
	}

	// Create the NAT gateways
	natGateway01Input := &ec2.CreateNatGatewayInput{
		AllocationId: Elasticip1,
		//AllocationId: natGateway01Output.AllocationId,
		SubnetId: privateSubnet01ID,
	}
	natGateway01Output, err := svc.CreateNatGateway(natGateway01Input)
	if err != nil {
		fmt.Println("Error creating NAT Gateway 01:", err)
		return
	}
	natGateway01ID := *natGateway01Output.NatGateway.NatGatewayId

	natGateway02Input := &ec2.CreateNatGatewayInput{
		AllocationId: Elasticip2,
		//AllocationId: natGateway02Output.AllocationId,
		SubnetId: publicSubnet02ID,
	}
	natGateway02Output, err := svc.CreateNatGateway(natGateway02Input)
	if err != nil {
		fmt.Println("Error creating NAT Gateway 02:", err)
		return
	}
	natGateway02ID := *natGateway02Output.NatGateway.NatGatewayId

	// Wait for NAT gateways to be available
	natGateway01Ready := false
	natGateway02Ready := false
	for !natGateway01Ready || !natGateway02Ready {
		time.Sleep(time.Second * 10)

		if !natGateway01Ready {
			natGateway01Output, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{
					aws.String(natGateway01ID),
				},
			})
			if err != nil {
				fmt.Println("Error describing NAT Gateway 01:", err)
				return
			}
			if len(natGateway01Output.NatGateways) > 0 && *natGateway01Output.NatGateways[0].State == "available" {
				natGateway01Ready = true
			}
		}

		if !natGateway02Ready {
			natGateway02Output, err := svc.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{
					aws.String(natGateway02ID),
				},
			})
			if err != nil {
				fmt.Println("Error describing NAT Gateway 02:", err)
				return
			}
			if len(natGateway02Output.NatGateways) > 0 && *natGateway02Output.NatGateways[0].State == "available" {
				natGateway02Ready = true
			}
		}
	}

	// Associate the public subnets with the public route table
	_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(publicRouteTableID),
		SubnetId:     publicSubnet01ID,
	})
	if err != nil {
		fmt.Println("Error associating public subnet 01 with public route table:", err)
	}

	_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(publicRouteTableID),
		SubnetId:     publicSubnet02ID,
	})
	if err != nil {
		fmt.Println("Error associating public subnet 02 with public route table:", err)
		return
	}

	// Associate the private subnets with the private route tables
	_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(privateRouteTable01ID),
		SubnetId:     privateSubnet01ID,
	})
	if err != nil {
		fmt.Println("Error associating private subnet 01 with private route table 01:", err)
		return
	}

	_, err = svc.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(privateRouteTable02ID),
		SubnetId:     privateSubnet02ID,
	})

	if err != nil {
		fmt.Println("Error associating private subnet 02 with private route table 02:", err)
		return
	}

	// Create the private route to the NAT gateways
	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		RouteTableId:         aws.String(privateRouteTable01ID),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		NatGatewayId:         aws.String(natGateway01ID),
	})
	if err != nil {
		fmt.Println("Error creating private route 01:", err)
		return
	}

	_, err = svc.CreateRoute(&ec2.CreateRouteInput{
		RouteTableId:         aws.String(privateRouteTable02ID),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		NatGatewayId:         aws.String(natGateway02ID),
	})
	if err != nil {
		fmt.Println("Error creating private route 02:", err)
		return
	}

	// security group

	inputsecuritygroup := &ec2.CreateSecurityGroupInput{
		Description: aws.String("My security group"),
		GroupName:   aws.String("my-security-group"),
		VpcId:       aws.String(vpcID),
	}

	securityGroupId1, err := svc.CreateSecurityGroup(inputsecuritygroup)
	if err != nil {
		fmt.Println("Error creating securitygroup", err)
		return
	}

	securitygrp := securityGroupId1.GroupId

	fmt.Println("security group created", securitygrp)

	fmt.Println("VPC creation complete!")
