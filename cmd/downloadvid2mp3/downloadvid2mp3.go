// Copyright (c) 2026 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiagomelo/vid2mp3"
	"github.com/tiagomelo/ytdld"
)

func run(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: downloadvid2mp3 <video-url> [output-file]")
	}

	ctx := context.Background()

	videoURL := strings.TrimSpace(args[0])
	if videoURL == "" {
		return errors.New("video url is required")
	}

	outputFile := strings.TrimSpace(args[1])
	if outputFile == "" {
		return errors.New("output file is required")
	}

	fmt.Printf("downloading video from %s to %s...\n", videoURL, outputFile)

	outputFilePath, err := ytdld.DownloadVideo(ctx, videoURL, removeExt(outputFile))
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	fmt.Printf("video downloaded successfully to %s\n", outputFilePath)
	defer deleteVideoFile(outputFilePath)

	fmt.Printf("converting video to mp3 and saving to %s...\n", outputFile)

	if err := vid2mp3.ExtractAudioFromVideo(ctx, outputFilePath, outputFile); err != nil {
		return fmt.Errorf("failed to convert to mp3: %w", err)
	}

	fmt.Printf("mp3 file created successfully at %s\n", outputFile)

	return nil
}

// deleteVideoFile deletes the video file at the given path.
func deleteVideoFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete video file %s: %w", path, err)
	}
	return nil
}

// removeExt removes the file extension from the given path.
func removeExt(path string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)]
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
