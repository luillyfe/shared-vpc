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
		gcpConfig := config.New(ctx, "service2")

		// Stack reference to Host project
		hostProject, err := pulumi.NewStackReference(ctx, "luillyfe/shared-vpc-host/dev", nil)
		if err != nil {
			return err
		}
		hostProjectId := hostProject.GetStringOutput(pulumi.String("hostProjectId"))
		organizationId := hostProject.GetStringOutput(pulumi.String("organizationId"))

		// Create the second service project and attach it to the host project
		serviceProject2, err := organizations.NewProject(ctx, "serviceProject2", &organizations.ProjectArgs{
			Name:      pulumi.String("serviceProject2"),
			ProjectId: hostProjectId,
			OrgId:     organizationId,
		})
		if err != nil {
			return err
		}

		_, err = compute.NewSharedVPCServiceProject(ctx, "sharedVPCServiceProject2", &compute.SharedVPCServiceProjectArgs{
			HostProject:    hostProjectId,
			ServiceProject: serviceProject2.ProjectId,
		})
		if err != nil {
			return err
		}

		// Create a virtual machine in service project 2
		_, err = compute.NewInstance(ctx, "serviceVm2", &compute.InstanceArgs{
			Project:     serviceProject2.ProjectId,
			Name:        pulumi.String("service-vm-2"),
			MachineType: pulumi.String(gcpConfig.Require("machineType")),
			Zone:        pulumi.String(gcpConfig.Require("machineZone")),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String(gcpConfig.Require("machineImage")),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: hostProjectId.ApplyT(func(id string) string {
						return "projects/" + id + "/global/networks/default"
					}).(pulumi.StringOutput),
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
