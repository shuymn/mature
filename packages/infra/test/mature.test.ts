import { SynthUtils } from "@aws-cdk/assert";
import * as cdk from "@aws-cdk/core";
import * as Mature from "../lib/mature-stack";

test("Mature Stack", () => {
  const app = new cdk.App();
  const stack = new Mature.MatureStack(app, "TestMatureStack");
  expect(SynthUtils.toCloudFormation(stack)).toMatchSnapshot();
});
