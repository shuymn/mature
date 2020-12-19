#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { MatureStack } from "../lib/mature-stack";
import { NotifyStack } from "../lib/notify-stack";

const app = new cdk.App();

new MatureStack(app, "mature-stack");
new NotifyStack(app, "notify-stack");
