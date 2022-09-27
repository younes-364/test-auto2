package utils

import (
    "fmt"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudformation"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    "github.com/aws/aws-sdk-go/service/servicecatalog"
    "github.com/aws/aws-sdk-go/service/sts"
    "github.com/gruntwork-io/terratest/modules/logger"
    "github.com/gruntwork-io/terratest/modules/testing"
)

func CollectProductLogs(t testing.TestingT, productName string) {
    logger.Log(t, ">>>>>>>>>>>>>>> COLLECTING PRODUCT LOGS <<<<<<<<<<<<<<<")
    accountId := getAccountId()
    stackArn := getCloudformationStackArn(productName)
    
    stackFailures := getCloudformationStackFailures(stackArn)
    logger.Log(t, ">>>>>>>>>>>>>>> @@stack-failures <<<<<<<<<<<<<<<")
    for _, event := range stackFailures {
        logger.Logf(t, "%v (%v at %v): %v", *event.LogicalResourceId, *event.ResourceStatus, 
            *event.Timestamp, *event.ResourceStatusReason)
    }

    stackResources := getCloudformationStackResources(stackArn)
    ec2Instance := findCloudformationInstanceResource(stackResources)
    logGroup := findCloudformationLogGroupResource(stackResources)
    if len(stackFailures) > 0 && logGroup != "" && ec2Instance != "" {
        logger.Log(t, ">>>>>>>>>>>>>>> @@cloud-init-log <<<<<<<<<<<<<<<")
        logStream := fmt.Sprintf("%v-cfn-init.log", ec2Instance)
        logEvents := getCloudWatchLogEvents(logGroup, logStream)
        
        // This method of exporting the logs breaks the terratest log parser output for now
        // logFile := "\n"
        // for _, log := range logEvents {
        //     logFile = logFile + fmt.Sprintln("    ", *log.Message)
        // }
        // logger.Log(t, logFile)

        // We need to print lines individually which results in two timestamps, one coming from
        // the terratest logger, and one from the CloudWatch message
        for _, log := range logEvents {
            logger.Logf(t, " | %v", *log.Message)
        }
    }

    logger.Log(t, ">>>>>>>>>>>>>>> @@details <<<<<<<<<<<<<<<")
    logger.Logf(t, "Account ID: %s", accountId)
    logger.Logf(t, "Provisioned product name: %s", productName)
    logger.Logf(t, "Cloudformation stack ARN: %s", stackArn)
    if ec2Instance != "" {
        logger.Logf(t, "EC2 instance ID: %s", ec2Instance)
    }
    if logGroup != "" {
        logger.Logf(t, "Cloud Watch log group: %s", logGroup)
    }
    logger.Log(t, ">>>>>>>>>>>>>>> END OF PRODUCT LOGS <<<<<<<<<<<<<<<")

}

func getAccountId() string {
    sess, err := session.NewSession()
    svc := sts.New(sess)

    // Get the STS caller identity, which includes the account ID
    resCaller, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput {})

    if err != nil {
        panic(err)
    }

    accountId := resCaller.Account

    return *accountId
}

func getCloudformationStackArn(productName string) string {
    sess, err := session.NewSession()
    svc := servicecatalog.New(sess)

    // Get the provisioned product id from the product details
    resProduct, err := svc.DescribeProvisionedProduct(&servicecatalog.DescribeProvisionedProductInput {
        Name: aws.String(productName),
    })

    if err != nil {
        panic(err)
    }

    provisioningRecord := resProduct.ProvisionedProductDetail.LastProvisioningRecordId

    // Get the last provisioning record detail
    resRecord, err := svc.DescribeRecord(&servicecatalog.DescribeRecordInput {
        Id: aws.String(*provisioningRecord),
    })

    if err != nil {
        panic(err)
    }

    cfnStack := ""
    for _, output := range resRecord.RecordOutputs {
        if *output.OutputKey == "CloudformationStackARN" {
            cfnStack = *output.OutputValue
        }
    }

    return cfnStack
}

func getCloudformationStackFailures(stackArn string) []cloudformation.StackEvent {
    sess, err := session.NewSession()
    svc := cloudformation.New(sess)
    
    // Get the errors from the CFN stack
    resEvents, err := svc.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
        StackName: aws.String(stackArn),
    })

    if err != nil {
        panic(err)
    }

    events := []cloudformation.StackEvent{}

    for _, event := range resEvents.StackEvents {
        if strings.HasSuffix(*event.ResourceStatus, "FAILED") {
            events = append(events, *event)
        }
    }

    return events
}

func getCloudformationStackResources(stackArn string) []*cloudformation.StackResource {
    sess, err := session.NewSession()
    svc := cloudformation.New(sess)

    // Retrieve the resources from the CloudFormation stack
    resResources, err := svc.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{
        StackName: aws.String(stackArn),
    })

    if err != nil {
        panic(err)
    }

    return resResources.StackResources
}

func findCloudformationInstanceResource(resources []*cloudformation.StackResource) string {
    ec2Instance := ""
    for _, resource := range resources {
        if *resource.ResourceType == "AWS::EC2::Instance" && resource.PhysicalResourceId != nil {
            ec2Instance = *resource.PhysicalResourceId
        }
    }

    return ec2Instance
}

func findCloudformationLogGroupResource(resources []*cloudformation.StackResource) string {
    logGroup := ""
    for _, resource := range resources {
        if *resource.ResourceType == "AWS::Logs::LogGroup" && resource.PhysicalResourceId != nil {
            logGroup = *resource.PhysicalResourceId
        }
    }

    return logGroup
}

func getCloudWatchLogEvents(logGroup string, logStream string) []*cloudwatchlogs.OutputLogEvent {
    sess, err := session.NewSession()
    svc := cloudwatchlogs.New(sess)
    
    // Retrieve the logs from CloudWatch
    resLogs, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
        LogGroupName: aws.String(logGroup),
        LogStreamName: aws.String(logStream),
    })

    if err != nil {
        panic(err)
    }

    return resLogs.Events
}