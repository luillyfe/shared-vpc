package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func projectService2(ctx *pulumi.Context, projectId pulumi.StringOutput) error {
	// Config
	gcpConfig := config.New(ctx, "gcp")

	// Create the second service project and attach it to the host project
	serviceProject2, err := organizations.NewProject(ctx, "serviceProject2", &organizations.ProjectArgs{
		Name:      pulumi.String("serviceProject2"),
		ProjectId: pulumi.String(gcpConfig.Require("projectId")),
		OrgId:     pulumi.String(gcpConfig.Require("organizationId")),
	})
	if err != nil {
		return err
	}

	_, err = compute.NewSharedVPCServiceProject(ctx, "sharedVPCServiceProject2", &compute.SharedVPCServiceProjectArgs{
		HostProject:    projectId,
		ServiceProject: serviceProject2.ProjectId,
	})
	if err != nil {
		return err
	}

	// Create a virtual machine in service project 2
	_, err = compute.NewInstance(ctx, "serviceVm2", &compute.InstanceArgs{
		Project:     serviceProject2.ProjectId,
		Name:        pulumi.String("service-vm-2"),
		MachineType: pulumi.String(gcpConfig.Require("serviceMachineType")),
		Zone:        pulumi.String(gcpConfig.Require("serviceMachineZone")),
		BootDisk: &compute.InstanceBootDiskArgs{
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String(gcpConfig.Require("serviceMachineImage")),
			},
		},
		NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
			&compute.InstanceNetworkInterfaceArgs{
				Network: projectId.ApplyT(func(id string) string {
					return "projects/" + id + "/global/networks/default"
				}).(pulumi.StringOutput),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
