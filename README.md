# acl-update
## About
This utility has been created as I have a need to update the list of CIDR values in the `loadBalancerSourceRanges:` section of service specs in Kubernetes.
This section, on cloud providers, is used to configure firewall rules for the load balancers created by the service when the service is of `type: LoadBalancer`

Tested on Azure
Not tested on GCP / AWS 

## Usage
Retrieve the current list (returns yaml style array)
```bash
acl-update get -service [servicename]
```

Example

```shell
./acl-update get -service fluentd-es
- 104.41.181.213/32
- 193.28.124.0/22
- 193.38.192.0/19
```


Add an entry (outputs new list in yaml)
Note: Always appends to end to avoid the cloud provider rewriting the entire security ruleset
```bash
acl-update add -service [servicename] -cidr [cidrstring]
```

Example
```shell
./acl-update add -service fluentd-es -cidr 204.212.100.34/32
- 104.41.181.213/32
- 193.28.124.0/22
- 193.38.192.0/19
- 204.212.100.34/32
```

Delete an entry (outputs new list in yaml)
Note: deletes from index - this will result in rules being rewritten in Azure.  Behaviour untested in AWS/GCP.
```bash
acl-update del -service [servicename] -cidr [cidrstring]
```

Example
```shell
./acl-update del -service fluentd-es -cidr 204.212.100.34/32
- 104.41.181.213/32
- 193.28.124.0/22
- 193.38.192.0/19
```

## Future features
- Some sort of reconcile option that will modify the list based on some known good list / source?
- Convert to an operator?