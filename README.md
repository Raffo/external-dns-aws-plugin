# ExternalDNS Plugin AWS Provider

This reimplements the AWS provider of ExternalDNS as a separate process. If it can be done for a very complicated provider, it can be done for any of them.

## How to test it

- Create a cluster.
- Run ExternalDNS locally with the plugin provider.
- Run the plugin locally.
- Test that we make the right calls to AWS.
