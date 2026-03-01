package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const apiURL = "https://markdown.new/"

type request struct {
	URL          string `json:"url"`
	Method       string `json:"method,omitempty"`
	RetainImages bool   `json:"retain_images,omitempty"`
}

type response struct {
	Success    bool   `json:"success"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Method     string `json:"method"`
	DurationMs int    `json:"duration_ms"`
	Tokens     int    `json:"tokens"`
}

func main() {
	method := flag.String("method", "auto", "conversion method: auto, ai, or browser")
	images := flag.Bool("images", false, "retain images in output")
	jsonOut := flag.Bool("json", false, "output raw JSON response")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: getm [flags] <url>\n\nFetch a URL and return it as clean Markdown.\nPowered by markdown.new (Cloudflare)\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	body, err := json.Marshal(request{
		URL:          url,
		Method:       *method,
		RetainImages: *images,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: HTTP %d: %s\n", resp.StatusCode, strings.TrimSpace(string(rawBody)))
		os.Exit(1)
	}

	if *jsonOut {
		fmt.Println(string(rawBody))
		return
	}

	var result response
	if err := json.Unmarshal(rawBody, &result); err != nil {
		// Not JSON — print raw (might be plain markdown)
		fmt.Print(string(rawBody))
		return
	}

	if !result.Success {
		fmt.Fprintf(os.Stderr, "error: conversion failed\n")
		os.Exit(1)
	}

	fmt.Println(result.Content)
}
