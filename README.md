# efk-e2e-tests
E2E tests for Elasticsearch-Fluentbit-Kibana stack.  
Currently, it contains a little of test cases to check healthy or unhealty status. 

## Getting started

```console
$ curl -L https://github.com/jabbukka/efk-e2e-tests/releases/download/0.1/efk-e2e-runner-0.1-linux-amd64.tar.gz --output ./efk-e2e-runner.tar.gz
$ tar xzf efk-e2e-runner.tar.gz && chmod +x efk-e2e-runner
$ ./efk-e2e-runner -kibana-host "http://$KIBANA_HOST:$KIBANA_PORT" -es-host "https://$ELASTICSEARCH_HOST:$ELASTICSEARCH_PORT"
```

## Test cases
### Elasticsearch
1. Search an existing index
2. Check data from the above index
  
### Kibana
1. Create `test` index pattern.
2. Delete `test` index pattern.