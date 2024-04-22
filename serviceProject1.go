package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func projectService1(ctx *pulumi.Context, projectId pulumi.StringOutput) error {
	// Create the first service project and attach it to the host project
	serviceProject1, err := organizations.NewProject(ctx, "serviceProject1", &organizations.ProjectArgs{
		Name:      pulumi.String("serviceProject1"),
		ProjectId: pulumi.String("your-project-id"),
		OrgId:     pulumi.String("1234567"),
	})
	if err != nil {
		return err
	}

	_, err = compute.NewSharedVPCServiceProject(ctx, "sharedVPCServiceProject1", &compute.SharedVPCServiceProjectArgs{
		HostProject:    projectId,
		ServiceProject: serviceProject1.ProjectId,
	})
	if err != nil {
		return err
	}

	// Create a virtual machine in service project 1
	_, err = compute.NewInstance(ctx, "serviceVm1", &compute.InstanceArgs{
		Project:     serviceProject1.ProjectId,
		Name:        pulumi.String("service-vm-1"),
		MachineType: pulumi.String("e2-medium"),
		Zone:        pulumi.String("us-central1-a"),
		BootDisk: &compute.InstanceBootDiskArgs{
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String("debian-cloud/debian-9"),
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
