import * as cdk from 'aws-cdk-lib';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecs from 'aws-cdk-lib/aws-ecs';
import * as ecs_patterns from 'aws-cdk-lib/aws-ecs-patterns';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as ecr_assets from 'aws-cdk-lib/aws-ecr-assets';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';

export class BudgetingAppStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create DynamoDB tables
    const usersTable = new dynamodb.Table(this, 'UsersTable', {
      tableName: 'Users',
      partitionKey: {
        name: 'userName',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const incomeTable = new dynamodb.Table(this, 'IncomeTable', {
      tableName: 'Income',
      partitionKey: {
        name: 'userId#month',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'incomeItemName',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const budgetTable = new dynamodb.Table(this, 'BudgetTable', {
      tableName: 'Budget',
      partitionKey: {
        name: 'userId#month',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'budgetItemName',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const expensesTable = new dynamodb.Table(this, 'ExpensesTable', {
      tableName: 'Expenses',
      partitionKey: {
        name: 'userId#month',
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: 'expenseItemName',
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });


    // Create VPC
    const vpc = new ec2.Vpc(this, 'BudgetingVPC', {
      ipAddresses: ec2.IpAddresses.cidr('10.0.0.0/16'),
      enableDnsHostnames: true,
      enableDnsSupport: true,
      maxAzs: 2,
      natGateways: 1,
      subnetConfiguration: [
        {
          cidrMask: 18,
          name: 'private',
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
        },
        {
          cidrMask: 18,
          name: 'public',
          subnetType: ec2.SubnetType.PUBLIC,
        },
      ],
    });

    const existingZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
      domainName: 'shivikasingh.com',
    })
    
    const certificate = new acm.DnsValidatedCertificate(this, 'CrossRegionCertificate', {
      domainName: '*.shivikasingh.com',
      hostedZone: existingZone,
      region: 'us-west-2',
    });


     // Build Docker images
     

    const backendImageAsset = new ecr_assets.DockerImageAsset(this, 'BackendImage', {
      directory:  '../backend',
      assetName: "BudgetingBackendImage"
    });

    const frontendImageAsset = new ecr_assets.DockerImageAsset(this, 'FrontendImage', {
      directory: '../frontend/budgeting-app',
      assetName: "BudgetingFrontendImage"
    });

    // Create frontend service
    const frontendService = new ecs_patterns.ApplicationLoadBalancedFargateService(
      this,
      'BudgetingFrontendService',
      {
        cpu: 256,
        desiredCount: 1,
        taskImageOptions: {
          image: ecs.ContainerImage.fromDockerImageAsset(frontendImageAsset),
          containerPort: 80,
        },
        memoryLimitMiB: 512,
        publicLoadBalancer: true,
        vpc: vpc,
        certificate: certificate,
        domainName: 'budgeting.shivikasingh.com',
        domainZone: existingZone
      }
    );

       // Create the task execution role
       const taskExecutionRole = new iam.Role(this, 'TaskExecutionRole', {
        roleName: "BudgetBackendRole",
        assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      });
  
      // Attach the AmazonECSTaskExecutionRolePolicy
      taskExecutionRole.addManagedPolicy(
        iam.ManagedPolicy.fromAwsManagedPolicyName(
          'service-role/AmazonECSTaskExecutionRolePolicy'
        )
      );
  
      // Create the DynamoDB permissions policy
      const dynamoDBPermissionsPolicy = new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: [
          'dynamodb:BatchGetItem',
          'dynamodb:GetItem',
          'dynamodb:Query',
          'dynamodb:Scan',
          'dynamodb:BatchWriteItem',
          'dynamodb:PutItem',
          'dynamodb:UpdateItem',
          'dynamodb:DeleteItem',
        ],
        resources: ['*'],
      });
  
    // Create backend service
    const backendService = new ecs_patterns.ApplicationLoadBalancedFargateService(
      this,
      'BudgetingBackendService',
      {
        cpu: 256,
        desiredCount: 1,
        taskImageOptions: {
          image: ecs.ContainerImage.fromDockerImageAsset(backendImageAsset),
          containerPort: 8080,
          taskRole: taskExecutionRole,
        },
        memoryLimitMiB: 512,
        vpc: vpc,
        certificate: certificate,
        domainName: 'backend.shivikasingh.com',
        domainZone: existingZone,
        publicLoadBalancer: true,
        healthCheckGracePeriod: cdk.Duration.seconds(600), 
      }
    );
    backendService.taskDefinition.addToTaskRolePolicy(dynamoDBPermissionsPolicy)

    // Create Route53 records
    // new route53.CnameRecord(this, 'FrontendCnameRecord', {
    //   zone: existingZone,
    //   recordName: 'budgeting.shivikasingh.com',
    //   domainName: frontendService.loadBalancer.loadBalancerDnsName
    // });

    // new route53.CnameRecord(this, 'BackendCnameRecord', {
    //   zone: existingZone,
    //   recordName: 'backend.shivikasingh.com',
    //   domainName: backendService.loadBalancer.loadBalancerDnsName,
    // });
  }
}                                 
