# ExternalDNS Plugin AWS Provider

**NOTE**: this was only used for development and it isn't in any way a final implementation. Do not use this code, will be deleted in the future.

This reimplements the AWS provider of ExternalDNS as a separate process. If it can be done for a very complicated provider, it can be done for any of them.

## How to test it

- Create a cluster.
- Run ExternalDNS locally with the plugin provider.
- Run the plugin locally.
- Test that we make the right calls to AWS.
