package test

import (
    "testing"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func Test0100(t *testing.T) {
    t.Parallel()

    workingDir := "./"
    runId := random.UniqueId()

    // defer test_structure.RunTestStage(t, "cleanup_terraform", func() {
    //     terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
    //     terraform.Destroy(t, terraformOptions)
    // })

    test_structure.RunTestStage(t, "deploy_terraform", func() {
        terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
            TerraformDir: workingDir,
            TerraformBinary: "terragrunt",
            EnvVars: map[string]string{
                "MPI_RUN_ID": runId,
            },
            NoColor: true,
        })
        test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)
        terraform.InitAndApply(t, terraformOptions)
    })
}