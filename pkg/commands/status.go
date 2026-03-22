package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/topher/cortex/internal/models"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check model health status",
	Long:  `Checks the health and status of all local AI models.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Checking Cortex model status...\n")

		manager := models.NewManager()
		status := manager.Status()

		fmt.Println(status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
