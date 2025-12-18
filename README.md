# Telegram bot on AWS

Telegram bot implementation along with serverless infrastructure on AWS Cloud.

> ⚠️ Status: WIP

**Note:**
This project has been developed with several objectives in mind:
- To utilize and demonstrate the integration of AWS, Go, GitHub Actions, and Terraform technologies.
- To serve as a testbed for evaluating recent changes in Terraform.
- To explore and assess the AWS DevOps Agent service.
- To gain familiarity with the development of Telegram bot functionalities.

## Requirements

See, [.github/actions/dependencies/action.yml](.github/actions/dependencies/action.yml) for Go and Terraform versions.

### Bootstrap

Before proceeding with deployment first we have to setup:
1. GitHub Actions OIDC integration with AWS, including creation of IAM Role which will be used for deployment. 
2. Terraform state bucket.

This has to be applied manually, see [bootstrap](./bootstrap/).

### GitHub Actions

See, [.github/README.md](./.github/README.md).

## Deployment

GitHub Actions workflow(s) are used to deploy the changes.
