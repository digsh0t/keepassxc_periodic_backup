package terraform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func ApplyS3Bucket(bucketName string) error {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	tfOutPath := "./tf_out"

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return err
	}

	workingDir := "./terraform/live/services/"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return err
	}

	planConfig := []tfexec.PlanOption{
		tfexec.Out(tfOutPath),
		tfexec.Var(fmt.Sprintf("bucket_name=%s", bucketName)),
	}
	_, err = tf.Plan(context.Background(), planConfig...)
	if err != nil {
		return err
	}

	planStr, err := tf.ShowPlanFile(context.Background(), tfOutPath)
	if err != nil {
		return err
	}

	_, err = json.Marshal(planStr)
	if err != nil {
		return err
	}
	//fmt.Println(string(b))

	applyConfig := []tfexec.ApplyOption{
		// tfexec.Var(fmt.Sprintf("bucket_name=%s", bucketName)),
		tfexec.DirOrPlan(tfOutPath),
	}
	err = tf.Apply(context.Background(), applyConfig...)
	if err != nil {
		return err
	}
	return nil
}
