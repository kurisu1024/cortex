package commands

import (
	"fmt"

	"github.com/kutidu2048/cortex/internal/job"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputDir      string
	voiceName      string
	background     string
	highVoicesOnly bool
	duration       int
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

		// Get background from flag or config
		if !cmd.Flags().Changed("background") {
			configBackground := viper.GetString("output.video.background")
			if configBackground != "" {
				background = configBackground
			}
		}

		// Get duration from flag or config
		if duration == 0 {
			duration = viper.GetInt("output.duration")
			if duration == 0 {
				duration = 10 // Fallback to 10 minutes
			}
		}

		// Validate duration
		if duration < 1 || duration > 60 {
			fmt.Printf("❌ Invalid duration: %d minutes. Must be between 1 and 60.\n", duration)
			return
		}

		manager := job.NewManager()
		jobID, err := manager.CreateJob(topic, outputDir, voiceName, background, highVoicesOnly, duration)
		if err != nil {
			fmt.Printf("❌ Error creating job: %v\n", err)
			return
		}

		if highVoicesOnly {
			fmt.Println("🎯 Using high-quality voices only")
		}

		fmt.Printf("⏱️  Target duration: %d minutes (~%d words)\n", duration, duration*150)

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
	generateCmd.Flags().StringVarP(&background, "background", "b", "", "video background style (ai-generated, gradient, solid, image)")
	generateCmd.Flags().BoolVarP(&highVoicesOnly, "high-voices-only", "H", false, "use only high-quality voices for generation")
	generateCmd.Flags().IntVarP(&duration, "duration", "d", 0, "target video duration in minutes (default: 10, max: 60)")
}
