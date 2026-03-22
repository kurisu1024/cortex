package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/topher/cortex/internal/models"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop local AI models",
	Long:  `Stops all running local AI models.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛑 Stopping Cortex models...")

		manager := models.NewManager()
		if err := manager.Stop(); err != nil {
			fmt.Printf("❌ Error stopping models: %v\n", err)
			return
		}

		fmt.Println("✅ Models stopped successfully")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
