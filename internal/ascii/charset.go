package ascii

// Preset character sets ordered from light (low density) to dark (high density).
// The ordering matters: index 0 represents empty/white space, and the last
// character represents the darkest/most filled pixel.

var presets = map[string]string{
	"simple":   " .,:;i1tfLCG08@",
	"detailed": " .'`^\",:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$",
	"block":    " \u2591\u2592\u2593\u2588",
	"braille":  "\u2800\u2801\u2803\u2807\u2847\u28C7\u28E7\u28F7\u28FF",
	"dots":     "\u2800\u2804\u2806\u2816\u2836\u2876\u28B6\u28FF",
	"ascii":    "  .:-=+*#%@",
}

// GetPreset returns the character set for a given preset name.
// Falls back to "simple" if the name is unknown.
func GetPreset(name string) string {
	if cs, ok := presets[name]; ok {
		return cs
	}
	return presets["simple"]
}

// ListPresets returns all available preset names.
func ListPresets() []string {
	return []string{"simple", "detailed", "block", "braille", "dots", "ascii"}
}
