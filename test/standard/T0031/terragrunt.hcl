include "root" {
  path = find_in_parent_folders()
}

locals {
  run_id = get_env("MPI_RUN_ID")
  global_var_file = get_env("MPI_GLOBAL_VAR_FILE")
  global_vars = yamldecode(file(find_in_parent_folders("${local.global_var_file}")))
}

inputs = {
  // General configuration
  dr_service_class = "None"
  backup_plan      = "none"
  image_template   = "EC2MutableRedhat8Oracle"
  os_type          = "Linux"
  subnet_exposure  = local.global_vars.subnet_exposure
  subnet_routable  = local.global_vars.subnet_routable

  // Custom parameters
  parameter_group = {
    db_size                   = 20
    instance_name             = "T0060-${local.run_id}",
    provisioned_product_name  = "T0060-${local.run_id}",
    vpc_id                    = local.global_vars.vpc_id,
    oracle_disk_1_volume_type = "gp2"
    oracle_disk_2_volume_type = "gp2"
    oracle_disk_1_size        = 30
    oracle_disk_2_size        = 20
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
    "local.purpose"           = "mpi"
    "local.vmsource"          = "new"
    "axa_db_create.db_name"   = "NETEST1"
    "axa_db_create.charset"   = "UTF8"
    "axa_db_create.release"   = 19
    "axa_db_create.blk_size"  = 8192
    "axa_db_create.lang"      = "AMERICAN"
    "axa_db_create.territory" = "AMERICA"
  }
}
