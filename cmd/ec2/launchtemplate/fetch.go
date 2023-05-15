/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fetch called")
		fetchAllTemplatesInRegion()
		// fetchTemplateById()
	},
}

func init() {
	LaunchtemplateCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchTemplateByName() {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if err != nil {
		panic("failed to load SDK configuration, " + err.Error())
	}

	// Create a new Amazon EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Specify the name of the Launch Template
	templateName := "taylor-workstation-windows"

	// Create the input parameters for the DescribeLaunchTemplates method
	input := &ec2.DescribeLaunchTemplatesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("launch-template-name"),
				Values: []string{templateName},
			},
		},
	}

	// Call the DescribeLaunchTemplates method to retrieve information about the Launch Template
	output, err := svc.DescribeLaunchTemplates(context.Background(), input)
	if err != nil {
		panic("failed to describe launch template, " + err.Error())
	}

	// Print the Launch Template information
	if len(output.LaunchTemplates) > 0 {
		launchTemplate := output.LaunchTemplates[0]
		fmt.Printf("Launch Template Name: %s\n", *launchTemplate.LaunchTemplateName)
		fmt.Printf("Launch Template ID: %s\n", *launchTemplate.LaunchTemplateId)
		fmt.Printf("Launch Template Version: %d\n", *launchTemplate.LatestVersionNumber)
	} else {
		fmt.Println("Launch Template not found")
	}
}

func fetchTemplateById() {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if err != nil {
		panic("failed to load SDK configuration, " + err.Error())
	}

	// Create a new Amazon EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Specify the launch template ID and version number
	templateID := "lt-0628c39b01b4d281b"
	versionNumber := "10"

	// Create the input parameters for the DescribeLaunchTemplateVersions method
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: &templateID,
		Versions:         []string{versionNumber},
	}

	// Call the DescribeLaunchTemplateVersions method to retrieve information about the launch template version
	output, err := svc.DescribeLaunchTemplateVersions(context.Background(), input)
	if err != nil {
		panic("failed to describe launch template version, " + err.Error())
	}

	// Print the launch template version information
	fmt.Printf("Launch template version %s:\n", versionNumber)
	fmt.Printf("AMI ID: %s\n", *output.LaunchTemplateVersions[0].LaunchTemplateData.ImageId)
	fmt.Printf("Instance Type: %s\n", string(output.LaunchTemplateVersions[0].LaunchTemplateData.InstanceType))
}

func fetchAllTemplatesInRegion() {
	saveAllLaunchTemplatesToFile()
}

func fetchAllLaunchTemplates() ([]types.LaunchTemplate, error) {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK configuration: %v", err)
	}

	// Create a new Amazon EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Create the input parameters for the DescribeLaunchTemplates method
	input := &ec2.DescribeLaunchTemplatesInput{}

	// Call the DescribeLaunchTemplates method to retrieve information about all Launch Templates
	output, err := svc.DescribeLaunchTemplates(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe Launch Templates: %v", err)
	}

	// Convert []LaunchTemplate to []*LaunchTemplate
	launchTemplates := make([]types.LaunchTemplate, len(output.LaunchTemplates))
	for i, lt := range output.LaunchTemplates {
		fmt.Println(*lt.LaunchTemplateName)
		launchTemplates[i] = lt
	}

	return launchTemplates, nil
}

func saveLaunchTemplateToFile(launchTemplate *types.LaunchTemplate) error {
	filename := fmt.Sprintf("%s_%d.json", *launchTemplate.LaunchTemplateName,
		*launchTemplate.LatestVersionNumber)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	jsonBytes, err := json.MarshalIndent(launchTemplate, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal Launch Template %s to JSON: %v",
			*launchTemplate.LaunchTemplateName, err)
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to write Launch Template %s to file %s: %v",
			*launchTemplate.LaunchTemplateName, filename, err)
	}

	fmt.Printf("Saved Launch Template %s version %d to file %s\n",
		*launchTemplate.LaunchTemplateName, *launchTemplate.LatestVersionNumber, filename)

	return nil
}

func saveAllLaunchTemplatesToFile() error {
	launchTemplates, err := fetchAllLaunchTemplates()
	if err != nil {
		return fmt.Errorf("failed to fetch Launch Templates: %v", err)
	}

	for _, launchTemplate := range launchTemplates {
		fmt.Printf("Launch Template %s\n", *launchTemplate.LaunchTemplateName)
		err := saveLaunchTemplateToFile(&launchTemplate)
		if err != nil {
			return fmt.Errorf("failed to save Launch Template %s: %v",
				*launchTemplate.LaunchTemplateName, err)
		}
	}

	return nil
}
