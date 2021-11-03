import * as path from "path";
import * as events from "@aws-cdk/aws-events";
import * as targets from "@aws-cdk/aws-events-targets";
import * as iam from "@aws-cdk/aws-iam";
import * as lambda from "@aws-cdk/aws-lambda";
import * as cdk from "@aws-cdk/core";

export class NotifyStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const role = new iam.Role(this, "notify-execution-role", {
      roleName: "mature-notify-execution-role",
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
          `arn:aws:ssm:${stack.region}:${stack.account}:parameter/mature/production/slack-token`,
          `arn:aws:ssm:${stack.region}:${stack.account}:parameter/mature/production/channel-id`,
        ],
      })
    );

    role.addToPolicy(
      new iam.PolicyStatement({
        effect: iam.Effect.ALLOW,
        actions: ["cloudwatch:GetMetricWidgetImage"],
        resources: ["*"],
      })
    );

    const fn = new lambda.Function(this, "notify-function", {
      functionName: "mature-notify",
      code: lambda.Code.fromAsset(path.resolve(__dirname, "../../notify/dist")),
      handler: "notify",
      runtime: lambda.Runtime.GO_1_X,
      memorySize: 128,
      timeout: cdk.Duration.seconds(30),
      role: role,
      tracing: lambda.Tracing.ACTIVE,
      currentVersionOptions: {
        removalPolicy: cdk.RemovalPolicy.RETAIN,
        retryAttempts: 0,
      },
      environment: {
        MATURE_SLACK_TOKEN_KEY: "/mature/production/slack-token",
        MATURE_SLACK_CHANNEL_ID_KEY: "/mature/production/channel-id",
      },
    });

    fn.currentVersion.addAlias("development");

    const prodVersion = lambda.Version.fromVersionArn(
      this,
      "notify-function-version-production",
      `${fn.functionArn}:${fn.currentVersion.version}`
    );
    const prodAlias = prodVersion.addAlias("production");

    const rule = new events.Rule(this, "notify-rule", {
      ruleName: "mature-notify-rule",
      schedule: events.Schedule.expression("cron(0 15 * * ? *)"),
    });

    rule.addTarget(new targets.LambdaFunction(prodAlias));
  }
}
