# Exoscale Account Wiper

This is a simple utility to delete all resources from an Exoscale account. This is useful for testing scenarios when you want to have a clean Exoscale account.

## Usage

You can download the utility from the [releases section]() for your platform. You can use it from the command line:

```
./exoscale-account-wiper --api-key API-KEY-HERE --api-secret API-SECRET-HERE [OPTIONS]
```

Optionally, you can pass the following parameters:

- `--nodelete` do not delete by default
- `--[no]instances` to delete or not to delete instances
- `--[no]templates` to delete or not to delete templates
- `--[no]pools` to delete or not to delete instance pools
- `--[no]sg` to delete or not to delete security groups
- `--[no]aa` to delete or not to delete anti-affinity groups
- `--[no]eip` to delete or not to delete elastic IPs
- `--[no]sshkeys` to delete or not to delete SSH keys
- `--[no]nlb` to delete or not to delete network load balancers
- `--[no]privnet` to delete or not to delete private networks
- `--[no]sos` to delete or not to delete object storage buckets and objects
- `--[no]dns` to delete or not to delete DNS zones
- `--[no]iam` to delete or not to delete IAM API keys
- `--iam-exclude-self` exclude current API key when deleting IAM API keys
- `--[no]runstatus` to delete or not to delete runstatus pages

You can also pass the parameters as environment variables:

```
DELETE=0
```
