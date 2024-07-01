#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { BudgetingAppStack } from '../lib/infra-stack';
import { CognitoStack } from '../lib/cognito-stack';

const app = new cdk.App();
new BudgetingAppStack(app, 'BudgetingAppStack', {
  env: { 
    account: process.env.CDK_DEFAULT_ACCOUNT, 
    region: process.env.CDK_DEFAULT_REGION 
  },
});

new CognitoStack(app, 'CognitoStack', {
  env: { 
    account: process.env.CDK_DEFAULT_ACCOUNT, 
    region: process.env.CDK_DEFAULT_REGION 
  },
});

