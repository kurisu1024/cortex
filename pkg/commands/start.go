package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/kutidu2048/cortex/internal/models"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local AI models",
	Long:  `Starts the required local AI models (LLM and TTS) for Cortex to function.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting Cortex models...")

		manager := models.NewManager()
		if err := manager.Start(); err != nil {
			fmt.Printf("❌ Error starting models: %v\n", err)
			return
		}

		fmt.Println("✅ Models started successfully")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
