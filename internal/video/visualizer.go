package video

// Visualizer handles audio visualization generation
type Visualizer struct{}

// NewVisualizer creates a new visualizer
func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

// VisualizationStyle represents different visualization styles
type VisualizationStyle string

const (
	StyleWaveform   VisualizationStyle = "waveform"
	StyleSpectrum   VisualizationStyle = "spectrum"
	StyleBars       VisualizationStyle = "bars"
	StyleCircular   VisualizationStyle = "circular"
)

// GetFFmpegFilter returns the ffmpeg filter for a visualization style
func (v *Visualizer) GetFFmpegFilter(style VisualizationStyle, width, height int) string {
	switch style {
	case StyleWaveform:
		return "showwaves=s=1920x1080:mode=line:rate=25:colors=0x00FF00"

	case StyleSpectrum:
		return "showfreqs=s=1920x1080:mode=line:colors=0x00FF00|0x0080FF"

	case StyleBars:
		return "showwaves=s=1920x1080:mode=cline:rate=25:colors=0x00FF00|0x0080FF"

	case StyleCircular:
		return "showcqt=s=1920x1080:fps=25"

	default:
		return "showwaves=s=1920x1080:mode=line:rate=25:colors=0x00FF00"
	}
}
