package terraform

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func ApplyS3Bucket(bucketName string) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := "./terraform/live/services/"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	// tf.SetEnv(map[string]string{"TF_VAR_bucket_name": bucketName})

	fmt.Println(state.FormatVersion) // "0.1"
	planConfig := []tfexec.PlanOption{
		tfexec.Out("./out.txt"),
		tfexec.Var(fmt.Sprintf("bucket_name=%s", bucketName)),
	}
	plan, err := tf.Plan(context.Background(), planConfig...)
	if err != nil {
		log.Fatalf("error running Plan: %s", err)
	}

	fmt.Println(plan)
}
