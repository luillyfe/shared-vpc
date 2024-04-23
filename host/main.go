package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Config
		gcpConfig := config.New(ctx, "host")

		// Create the host project for the shared VPC
		hostProject, err := organizations.NewProject(ctx, "hostProject", &organizations.ProjectArgs{
			Name:      pulumi.String("hostProject"),
			ProjectId: pulumi.String(gcpConfig.Require("projectId")),
			OrgId:     pulumi.String(gcpConfig.Require("organizationId")),
		})
		if err != nil {
			return err
		}

		// Enable the Shared VPC feature for host project
		_, err = compute.NewSharedVPCHostProject(ctx, "sharedVPCHostProject", &compute.SharedVPCHostProjectArgs{
			Project: hostProject.ProjectId,
		})
		if err != nil {
			return err
		}

		// Create the virtual machine in host project
		_, err = compute.NewInstance(ctx, "hostVm", &compute.InstanceArgs{
			Name:        pulumi.String("host-vm"),
			MachineType: pulumi.String(gcpConfig.Require("machineType")),
			Zone:        pulumi.String(gcpConfig.Require("zone")),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String(gcpConfig.Require("machineImage")),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: hostProject.ProjectId.ApplyT(func(id string) string {
						return "projects/" + id + "/global/networks/default"
					}).(pulumi.StringOutput),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the Project ID
		ctx.Export("hostProjectId", hostProject.ProjectId)
		ctx.Export("organizationId", hostProject.OrgId)
		return nil
	})
}
