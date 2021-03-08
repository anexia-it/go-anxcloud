# Testing

`Ginkgo` is used as test suite to implement integration tests for each of the sdk packages. The integration tests are 
located under `./tests`.

To execute the complete test suite just run  `ginkgo ./tests`.  
To specifiy a subset of testcases use  the `-focus` flag and provide a regex to match the description of the tests 
used with `Describe(...)`. For example to execute all `core` tests run `ginkgo -focus="Core API endpoint tests"`.