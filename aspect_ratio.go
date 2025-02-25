package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func getVideoAspectRatio(path string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", path)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("issue running ffprobe: %v", err)
	}

	var data struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	err = json.Unmarshal(out.Bytes(), &data)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal json: %v", err)
	}
	if len(data.Streams) == 0 {
		return "", fmt.Errorf("video metadata not found")
	}

	width := data.Streams[0].Width
	height := data.Streams[0].Height
	if width == 16*height/9 {
		return "16:9", nil
	} else if width == 9*height/16 {
		return "9:16", nil
	}
	return "other", nil
}
