terraform {
  source = local.module_uri
}

//
// It's always required to specify the MPI module version and a global variable file
locals {
  module_uri = get_env("MPI_MODULE_URI")
  global_var_file = get_env("MPI_GLOBAL_VAR_FILE")
  global_vars = yamldecode(file(find_in_parent_folders("${local.global_var_file}")))
}

//
// Indicate what region to deploy the resources into overriding any module defaults
generate "provider" {
  path = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents = <<EOF
provider "aws" {
  region = "${local.global_vars.aws_region}"
}
EOF
}
