package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"os/exec"

	opts "github.com/dieple/tfrunner/options"
)


func TFRun(args opts.Options, bd *opts.BuildData) {
	// fmt.Printf("AutoApprove: %v, Bucket: %v, Modules: %v\n", bd.AutoApprove, bd.Bucket, bd.Modules)
	for _, m := range bd.Modules {
	    runModules (args, m, bd)
	}
}

func runModules(args opts.Options, module_path string, bd *opts.BuildData) {

	curr_path, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }

    module_name := strings.TrimRight(module_path[strings.LastIndex(module_path, "/")+1:], "\r\n")
    fmt.Println("Running terraform action: **" + bd.TFaction +"** with module: ******************************* " + module_name + " *******************************")

    // soft linking files and run terraform command...
    softLinkingFiles(curr_path, module_path, bd, args)

    // generate and run terraform command based on current module
    runTerraformCommand(curr_path, module_path, bd, args)
}

func runTerraformCommand(curr_path string, module_path string, bd *opts.BuildData, args opts.Options){

    err := os.Chdir(strings.TrimRight(module_path, "\r\n"))
    if err != nil {
        log.Fatalf("os.Chdir Error:- %s\n", err)
    }

    module_name := strings.TrimRight(module_path[strings.LastIndex(module_path, "/")+1:], "\r\n")
    key_config := fmt.Sprintf("\"key=%s/terraform.tfstate\"", module_name)
    bucket_region_config := fmt.Sprintf("\"region=%s\"", bd.BucketRegion)
    bucket_config := fmt.Sprintf("\"bucket=%s\"", bd.Bucket)
    dynamodb_config := fmt.Sprintf("\"dynamodb_table=%s\"", bd.Dynamodb)

    tf_init_cmd := fmt.Sprintf("terraform init --backend-config=%s --backend-config=%s --backend-config=%s --backend-config=%s && terraform workspace new %s || terraform workspace select %s",
                                key_config, bucket_region_config, dynamodb_config, bucket_config, bd.Workspace, bd.Workspace)

    // fmt.Println("tf_init_cmd:\n" + tf_init_cmd)

    tf_plan_cmd := fmt.Sprintf("terraform workspace select %s || terraform workspace new %s && terraform plan -out %s",
                                bd.Workspace, bd.Workspace, opts.PlanOutputFile)
    // fmt.Println("tf_plan_cmd:\n" + tf_plan_cmd)

    tf_plan_destroy_cmd := fmt.Sprintf("terraform workspace select %s || terraform workspace new %s && terraform plan -destroy -out %s",
                                        bd.Workspace, bd.Workspace, opts.PlanOutputFile)
    // fmt.Println("tf_plan_destroy_cmd:\n" + tf_plan_destroy_cmd)

    tf_apply_cmd := fmt.Sprintf("terraform workspace select %s || terraform workspace new %s && terraform apply %s",
                                bd.Workspace, bd.Workspace, opts.PlanOutputFile)
    // fmt.Println("tf_apply_cmd:\n" + tf_apply_cmd)

    if err := executeCommand(tf_init_cmd); err != nil {
        log.Fatalf("executeCommand() failed with %s\n", err)
    }

    if bd.TFaction == "plan" {
        if err = executeCommand(tf_plan_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
    } else if bd.TFaction == "plan-destroy" {
        if err = executeCommand(tf_plan_destroy_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
    } else if bd.TFaction == "apply" {
        if err = executeCommand(tf_plan_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
        if err = executeCommand(tf_apply_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
    } else if bd.TFaction == "apply-destroy" {
        if err = executeCommand(tf_plan_destroy_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
        if err = executeCommand(tf_apply_cmd); err != nil {
            log.Fatalf("executeCommand() failed with %s\n", err)
        }
    }

    fmt.Println("Module %s ran successfully...", module_name)
}

func executeCommand(cmd string) (err error){
    tf_cmd := exec.Command(cmd)
    err = tf_cmd.Run()
    if err != nil {
        log.Fatalf("executeCommand: tf_cmd.Run() failed with %s\n", err)
    }
    return err
}