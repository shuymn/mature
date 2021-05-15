import * as path from "path";
import * as targets from "@aws-cdk/aws-events-targets";
import * as events from "@aws-cdk/aws-events";
import * as iam from "@aws-cdk/aws-iam";
import * as lambda from "@aws-cdk/aws-lambda";
import * as cdk from "@aws-cdk/core";

export class MatureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const role = new iam.Role(this, "mature-execution-role", {
      roleName: "mature-execution-role",
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
      ],
    });

    const stack = cdk.Stack.of(this);
    role.addToPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["ssm:GetParameters"],
        resources: [
          `arn:aws:ssm:${stack.region}:${stack.account}:parameter/mature/production/access-token`,
          `arn:aws:ssm:${stack.region}:${stack.account}:parameter/mature/production/device-id/*`,
        ],
      })
    );

    role.addToPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["cloudwatch:PutMetricData"],
        resources: ["*"],
        conditions: {
          StringEquals: {
            "cloudwatch:namespace": "NatureRemo/RoomMetrics",
          },
        },
      })
    );

    const fn = new lambda.Function(this, "mature-function", {
      functionName: "mature",
      code: lambda.Code.fromAsset(path.resolve(__dirname, "../../mature/dist")),
      handler: "mature",
      runtime: lambda.Runtime.GO_1_X,
      memorySize: 256,
      timeout: cdk.Duration.seconds(30),
      role: role,
      tracing: lambda.Tracing.ACTIVE,
      currentVersionOptions: {
        removalPolicy: cdk.RemovalPolicy.RETAIN,
        retryAttempts: 0,
      },
      environment: {
        MATURE_ACCESS_TOKEN_KEY: "/mature/production/access-token",
        MATURE_DEVICE_ID_KEY: "/mature/production/device-id/main-room",
      },
    });

    fn.currentVersion.addAlias("development");

    const prodVersion = lambda.Version.fromVersionArn(
      this,
      "mature-function-version-production",
      `${fn.functionArn}:${fn.currentVersion.version}`
    );
    const prodAlias = prodVersion.addAlias("production");

    const rule = new events.Rule(this, "mature-rule", {
      ruleName: "mature-rule",
      schedule: events.Schedule.expression("rate(1 minute)"),
    });

    rule.addTarget(new targets.LambdaFunction(prodAlias));
  }
}
