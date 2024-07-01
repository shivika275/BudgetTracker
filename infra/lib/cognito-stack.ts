import * as cdk from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';

export class CognitoStack extends cdk.Stack {
  userPool: cognito.UserPool;
  userPoolClient: cognito.UserPoolClient;

  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create Cognito user pool
    const userPool = new cognito.UserPool(this, 'BudgetingAppUserPool', {
      userPoolName: 'budgeting-app-user-pool',
      selfSignUpEnabled: true,
      signInAliases: {
        username: true,
        email: true,
      },
      autoVerify: {
        email: true,
      },
      standardAttributes: {
        email: {
          required: true,
          mutable: true,
        },
      },
    });

    // Create Cognito user pool client
    const userPoolClient = userPool.addClient('BudgetingAppUserPoolClient', {
      userPoolClientName: 'budgeting-app-user-pool-client',
      authFlows: {
        adminUserPassword: true,
        custom: true,
        userPassword: true,
        userSrp: true,
      },
      generateSecret: false,
    });

    // Expose the user pool and client for other stacks
    this.userPool = userPool;
    this.userPoolClient = userPoolClient;
  }
}
