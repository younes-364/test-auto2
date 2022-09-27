package test

import (
    "fmt"
    "testing"

    "mpiawstests/utils"

    "github.com/gruntwork-io/terratest/modules/random"
    "github.com/gruntwork-io/terratest/modules/terraform"
    test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestBronzeDrOracleRedHat8(t *testing.T) {
    t.Parallel()

    workingDir := "./"
    productName := "bronze-dr-oracle-redhat8"
    runId := random.UniqueId()
    uniqueProductName := fmt.Sprintf("%s-%s", productName, runId)

    defer test_structure.RunTestStage(t, "cleanup_terraform", func() {
        terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
        terraform.Destroy(t, terraformOptions)
    })

    defer test_structure.RunTestStage(t, "collect_logs", func() {
        utils.CollectProductLogs(t, uniqueProductName)
    })

    test_structure.RunTestStage(t, "deploy_terraform", func() {
        terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
            TerraformDir: workingDir,
            TerraformBinary: "terragrunt",
            EnvVars: map[string]string{
                "MPI_PRODUCT_NAME": uniqueProductName,
            },
            NoColor: true,
        })
        test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)
        terraform.InitAndApply(t, terraformOptions)
    })
}