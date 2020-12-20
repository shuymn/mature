import "@aws-cdk/assert/jest";
import * as cdk from "@aws-cdk/core";
import { MatureStack } from "../lib/mature-stack";

describe("fine-grained assertions", () => {
  test("iam role", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::IAM::Role", {
      AssumeRolePolicyDocument: {
        Version: "2012-10-17",
        Statement: [
          {
            Action: "sts:AssumeRole",
            Effect: "Allow",
            Principal: {
              Service: "lambda.amazonaws.com",
            },
          },
        ],
      },
      ManagedPolicyArns: [
        {
          "Fn::Join": [
            "",
            [
              "arn:",
              {
                Ref: "AWS::Partition",
              },
              ":iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
            ],
          ],
        },
      ],
      RoleName: "mature-execution-role",
    });
  });

  test("iam policy", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::IAM::Policy", {
      PolicyDocument: {
        Statement: [
          {
            Action: "ssm:GetParameters",
            Effect: "Allow",
            Resource: [
              {
                "Fn::Join": [
                  "",
                  [
                    "arn:aws:ssm:",
                    {
                      Ref: "AWS::Region",
                    },
                    ":",
                    {
                      Ref: "AWS::AccountId",
                    },
                    ":parameter/mature/production/access-token",
                  ],
                ],
              },
              {
                "Fn::Join": [
                  "",
                  [
                    "arn:aws:ssm:",
                    {
                      Ref: "AWS::Region",
                    },
                    ":",
                    {
                      Ref: "AWS::AccountId",
                    },
                    ":parameter/mature/production/device-id/*",
                  ],
                ],
              },
            ],
          },
          {
            Action: "cloudwatch:PutMetricData",
            Condition: {
              StringEquals: {
                "cloudwatch:namespace": "NatureRemo/RoomMetrics",
              },
            },
            Effect: "Allow",
            Resource: "*",
          },
          {
            Action: ["xray:PutTraceSegments", "xray:PutTelemetryRecords"],
            Effect: "Allow",
            Resource: "*",
          },
        ],
        Version: "2012-10-17",
      },
    });
  });

  test("lambda function", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::Lambda::Function", {
      Handler: "mature",
      Runtime: "go1.x",
      Environment: {
        Variables: {
          MATURE_ACCESS_TOKEN_KEY: "/mature/production/access-token",
          MATURE_DEVICE_ID_KEY: "/mature/production/device-id/main-room",
        },
      },
      FunctionName: "mature",
      MemorySize: 256,
      Timeout: 30,
      TracingConfig: {
        Mode: "Active",
      },
    });
  });

  test("lambda event invoke config", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::Lambda::EventInvokeConfig", {
      MaximumRetryAttempts: 0,
    });
  });

  test("lambda alias", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::Lambda::Alias", {
      Name: "production",
    });
  });

  test("event rule", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::Events::Rule", {
      Name: "mature-rule",
      ScheduleExpression: "rate(1 minute)",
      State: "ENABLED",
    });
  });

  test("lambda permission", () => {
    const app = new cdk.App();
    const stack = new MatureStack(app, "TestMatureStack");

    expect(stack).toHaveResource("AWS::Lambda::Permission", {
      Action: "lambda:InvokeFunction",
      Principal: "events.amazonaws.com",
    });
  });
});
