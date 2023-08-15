package components

import (
	"fmt"
	"math"
	"pscreenapp/bridge"
	"pscreenapp/bridge/encoder"
	"pscreenapp/config"
	"pscreenapp/constants"
	"pscreenapp/ui/styling"
	"strings"

	"github.com/muesli/termenv"
)

func VersionNumber() string {
	return "pScreen " + styling.ColorFg(constants.AppVersion, styling.HighlightedColor)
}

func KeybindsHints(keybinds []string) string {
	s := ""
	for index, keybind := range keybinds {
		if index != 0 {
			s += styling.Dot
		}
		s += styling.Subtle(keybind)
	}
	return styling.Indent(s, config.Config.UIPaddingIndentAmount)
}

func Checkbox(label string, checked bool, selected bool) string {
	s := fmt.Sprintf("[ ] %s", label)
	if checked {
		s = "[x] " + label
	}
	if selected {
		return styling.ColorFg(s, styling.HighlightedColor)
	}
	return s
}

func Progressbar(width int, percent float64) string {
	w := float64(styling.ProgressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(styling.ProgressFullChar).Foreground(styling.Term.Color(styling.Ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(styling.ProgressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

func Header(header string) string {
	return styling.Bold(styling.Indent(header, config.Config.UIPaddingIndentAmount))
}

func StatsFooter() string {
	s := styling.Subtle(fmt.Sprintf("%02.4f fps (%04.1fms)", 10e8/float64(bridge.FrameDeltaTime), float64(bridge.FrameDeltaTime)/10e5))
	s += styling.Dot
	s += styling.Subtle(fmt.Sprintf("%d bytes (%0.4f%%)", encoder.CompressedBytesN, 100*float64(encoder.CompressedBytesN)/float64(encoder.UncompressedBytesN)))
	return styling.Indent(s, config.Config.UIPaddingIndentAmount)
}
