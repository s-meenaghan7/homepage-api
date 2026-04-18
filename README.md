# Homepage API
Backend API and infrastructure for my website, <a href="https://seanmeenaghan.com" target="_blank">seanmeenaghan.com</a>.
- APIs written in **Go** and deployed as AWS **Lambda** functions.
- API routing and configuration via AWS **API Gateway**.
- Infrastructure provisioned with **Terraform HCL**.

## Architecture
### Visitor API
This API records non-unique page loads for a given webpage in a DynamoDB table. This is my first API built as part of the [Cloud Resume Challenge](https://cloudresumechallenge.dev/docs/the-challenge/aws/).

<img alt="AWS API hexagonal architecture pattern" src="https://docs.aws.amazon.com/images/prescriptive-guidance/latest/cloud-design-patterns/images/hexagonal-2.png">
