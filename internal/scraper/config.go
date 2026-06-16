package scraper

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

//	func saveConfig(c ConfigFile, filename string) error {
//		// 1. Convert struct to JSON bytes
//		// "" is the prefix (none), "  " is the indentation
//		data, err := json.MarshalIndent(c, "", "  ")
//		if err != nil {
//			return err
//		}
//
//		// 2. Write to disk (0644 is standard file permission)
//		return os.WriteFile(filename, data, 0o644)
//	}

func (c *ConfigFile) ToConfig() (config, error) {
	var cfg config

	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return cfg, err
	}
	seriesPath, err := url.Parse(c.SeriesPath)
	if err != nil {
		return cfg, err
	}
	cfg = config{
		baseURL:           baseURL,
		seriesPath:        seriesPath,
		volumePattern:     c.VolumePattern,
		volumePatternStop: c.VolumePatternStop,
		chapterPattern:    c.ChapterPattern,
	}
	return cfg, nil
}

func getConfig() (config, error) {
	var cfgFile ConfigFile
	var cfg config
	filePath := filepath.Join("internal", "scraper", "config.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfgFile)
	if err != nil {
		return cfg, err
	}
	cfg, err = cfgFile.ToConfig()
	if err != nil {
		return cfg, err
	}
	return cfg, err
}

type config struct {
	baseURL           *url.URL
	seriesPath        *url.URL
	volumePattern     string
	volumePatternStop string
	chapterPattern    string
}
type ConfigFile struct {
	BaseURL           string `json:"base_url"`
	SeriesPath        string `json:"series_path"`
	VolumePattern     string `json:"volume_pattern"`
	VolumePatternStop string `json:"volume_pattern_stop"`
	ChapterPattern    string `json:"chapter_pattern"`
}
