package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create the host project for the shared VPC
		hostProject, err := organizations.NewProject(ctx, "hostProject", &organizations.ProjectArgs{
			Name:      pulumi.String("hostProject"),
			ProjectId: pulumi.String("your-project-id"),
			OrgId:     pulumi.String("1234567"),
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
			MachineType: pulumi.String("e2-medium"),
			Zone:        pulumi.String("us-central1-a"),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String("debian-cloud/debian-9"),
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

		// Create project service1 and resources
		err = projectService1(ctx, hostProject.ProjectId)
		if err != nil {
			return err
		}

		// Create project service2 and resources
		err = projectService2(ctx, hostProject.ProjectId)
		if err != nil {
			return err
		}

		return nil
	})
}
