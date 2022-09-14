# MPI Test Stack

This stack is aimed at testing MPI modules, at scale, and using automation. It leverages Terragrunt to wrap the MPI modules and Terratest to run programmatic tests written in Go. The latter is also used by Hashicorp to test Terraform stacks, making it the leading tool for testing Terraform code.

Using Terragrunt allows us to test the module produced by the MPI product team directly, without having to wrap it in another Terraform module (which would possibly introduce error cases). Terragrunt also allows us to dynamically alter the module versions on test runs, which would be impossible in Terraform because the source of a module must always be hard-coded.

Using Terratest adds a programmatic layer to our testing stack which can be used to experiment with test cases that go beyond simply invoking Terraform. For example, in Terratest we can create a VM with a set of parameters, update one of the parameters using Terraform, and monitor that the VM always remains running during the update (i.e. testing zero-downtime updates). Terratest also allows us to pull log files from various sources (Terraform, Cloudformation, CloudWatch Logs) on every test run, which can then be forwarded to the product teams when we spot a failure.
