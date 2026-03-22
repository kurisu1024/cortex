package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/topher/cortex/internal/job"
)

var (
	outputDir      string
	voiceName      string
	background     string
	highVoicesOnly bool
)

var generateCmd = &cobra.Command{
	Use:   "generate [topic]",
	Short: "Generate script, audio, and video",
	Long:  `Generates a complete video from a topic using local AI models.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]

		fmt.Printf("🎬 Generating content for: %s\n\n", topic)

		// Get config values with fallbacks
		if outputDir == "" {
			outputDir = viper.GetString("output.directory")
			if outputDir == "" {
				outputDir = "./output"
			}
		}

		manager := job.NewManager()
		jobID, err := manager.CreateJob(topic, outputDir, voiceName, background, highVoicesOnly)
		if err != nil {
			fmt.Printf("❌ Error creating job: %v\n", err)
			return
		}

		if highVoicesOnly {
			fmt.Println("🎯 Using high-quality voices only")
		}

		if err := manager.RunJob(jobID); err != nil {
			fmt.Printf("❌ Job failed: %v\n", err)
			return
		}

		fmt.Println("\n✅ Video generated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory for generated files")
	generateCmd.Flags().StringVarP(&voiceName, "voice", "v", "", "TTS voice to use")
	generateCmd.Flags().StringVarP(&background, "background", "b", "gradient", "video background style (gradient, solid, image)")
	generateCmd.Flags().BoolVarP(&highVoicesOnly, "high-voices-only", "H", false, "use only high-quality voices for generation")
}
