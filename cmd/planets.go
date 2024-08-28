package cmd

import (
	"ocp/sample/planets/internal/cms"
	"os"

	"github.com/spf13/cobra"
)

var PlanetsCmd = &cobra.Command{
	Use:   "planets",
	Short: "A cli to manage CMS data",
	Long:  `This cli creates, updates, deletes and prints CMS instance data.`,
}

var cmsCreatePlanetsCmd = &cobra.Command{
	Use:   "create",
	Short: "Create planet CMS instances based on sample data.",
	Run: func(cmd *cobra.Command, args []string) {
		cms.CreatePlanets()
	},
}

var cmsUpdatePlanetsCmd = &cobra.Command{
	Use:   "update",
	Short: "Update planet CMS instances.",
	Run: func(cmd *cobra.Command, args []string) {
		cms.UpdatePlanets()
	},
}

var cmsDeletePlanetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete planet CMS instances.",
	Run: func(cmd *cobra.Command, args []string) {
		cms.DeletePlanets()
	},
}

var cmsInfoPlanetsCmd = &cobra.Command{
	Use:   "info",
	Short: "Print planet CMS instance info.",
	Run: func(cmd *cobra.Command, args []string) {
		cms.PlanetInfo()
	},
}

func Execute() {
	err := PlanetsCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	PlanetsCmd.AddCommand(cmsInfoPlanetsCmd)
	PlanetsCmd.AddCommand(cmsCreatePlanetsCmd)
	PlanetsCmd.AddCommand(cmsUpdatePlanetsCmd)
	PlanetsCmd.AddCommand(cmsDeletePlanetsCmd)
}
