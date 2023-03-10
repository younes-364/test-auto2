include "root" {
  path = find_in_parent_folders()
}

locals {
  product_name = get_env("MPI_PRODUCT_NAME")
  global_var_file = get_env("MPI_GLOBAL_VAR_FILE")
  global_vars = yamldecode(file(find_in_parent_folders("${local.global_var_file}")))
}

inputs = {
  // General configuration
  dr_service_class = "None"
  backup_plan      = "medium"
  image_template   = "EC2MutableRedhat8Base"
  os_type          = "Linux"
  subnet_exposure  = local.global_vars.subnet_exposure
  subnet_routable  = local.global_vars.subnet_routable

  // Custom parameters
  parameter_group = {
    data_disk_size            = 10
    instance_name             = local.product_name
    provisioned_product_name  = local.product_name
    vpc_id                    = local.global_vars.vpc_id
  }

  // Mandatory tags
  global_app          = local.global_vars.global_app
  global_appserviceid = local.global_vars.global_appserviceid
  global_broker       = local.global_vars.global_broker
  global_cbp          = local.global_vars.global_cbp
  global_dataclass    = local.global_vars.global_dataclass
  global_dcs          = local.global_vars.global_dcs
  global_env          = local.global_vars.global_env
  global_opco         = local.global_vars.global_opco
  global_project      = local.global_vars.global_project

  // Custom tags
  custom_tags         = {
    "local.purpose"  = "mpi"
    "local.vmsource" = "new"
  }
}