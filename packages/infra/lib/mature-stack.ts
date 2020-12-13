import * as path from "path";
import * as events from "@aws-cdk/aws-events";
import * as targets from "@aws-cdk/aws-events-targets";
import * as iam from "@aws-cdk/aws-iam";
import * as lambda from "@aws-cdk/aws-lambda";
import * as cdk from "@aws-cdk/core";

export class MatureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const matureExecutionRole = new iam.Role(this, "mature-execution-role", {
      roleName: "mature-execution-role",
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
      ],
    });

    const stack = cdk.Stack.of(this);
    matureExecutionRole.addToPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["ssm:GetParameters"],
        resources: [
          `arn:aws:ssm:${stack.region}:${stack.account}:parameter/mature/production/*`,
        ],
      })
    );

    matureExecutionRole.addToPolicy(
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

    const matureFn = new lambda.Function(this, "mature-function", {
      functionName: "mature",
      code: lambda.Code.fromAsset(path.resolve(__dirname, "../../mature/dist")),
      handler: "mature",
      runtime: lambda.Runtime.GO_1_X,
      memorySize: 256,
      timeout: cdk.Duration.seconds(30),
      role: matureExecutionRole,
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

    matureFn.currentVersion.addAlias("development");

    const matureFnProdVersion = lambda.Version.fromVersionArn(
      this,
      "mature-function-version-production",
      `${matureFn.functionArn}:7`
    );
    const matureFnProdAlias = matureFnProdVersion.addAlias("production");

    const matureRule = new events.Rule(this, "mature-rule", {
      ruleName: "mature-rule",
      schedule: events.Schedule.expression("rate(1 minute)"),
    });

    matureRule.addTarget(new targets.LambdaFunction(matureFnProdAlias));
  }
}
