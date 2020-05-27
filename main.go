/*
 * Copyright (c) 2019 AudD, LLC. All rights reserved.
 * Copyright (c) 2019 Mikhail Samin. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *    * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *    * Neither the name of Mikhail Samin, AudD, nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Mihonarium/ytdl"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

type AudDResponse struct {
	Status string `json:"status"`
	Error  struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
	Result []SongInfo `json:"result"`
}

type SongInfo struct {
	Artist      string `json:"artist"`
	Title       string `json:"title"`
	Album       string `json:"album"`
	ReleaseDate string `json:"release_date"`
	Label       string `json:"label"`
	Underground bool   `json:"underground"`
	TimeCode    string `json:"timecode"`
	Start       int    `json:"start,omitempty"`
	End         int    `json:"end,omitempty"`
}

func SecondsToTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func AudDAPIRecognizeAll(reader io.Reader, apiToken string) []SongInfo {
	apiResponse := AudDAPIUpload(map[string]string{"all": "true"}, reader, apiToken)
	var result AudDResponse
	json.Unmarshal(apiResponse, &result)
	//fmt.Println(result.Result)
	if result.Status != "success" {
		fmt.Println("Request failed:", result.Error.ErrorMessage)
		return nil
	}
	return result.Result
}

func AudDAPIUpload(params map[string]string, reader io.Reader, apiToken string) []byte {
	params["api_token"] = apiToken
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "file")
	if err != nil {
		panic(err)
	}
	io.Copy(part, reader)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()
	req, _ := http.NewRequest("POST", "https://api.audd.io/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody
}

func DownloadYoutubeVideo(Url string, file *os.File) error {
	c := ytdl.DefaultClient
	vid, err := c.GetVideoInfo(context.TODO(), Url)
	if err != nil {
		return err
	}
	if len(vid.Formats) == 0 {
		return fmt.Errorf("ytdl can't find available formats")
	}
	return c.Download(context.TODO(), vid, vid.Formats[0], file)
}

func RecognizeMultipleFiles(untilFirst bool, dir string, secondsPerFile int, apiToken string) []SongInfo {
	files, _ := ioutil.ReadDir(dir)
	if untilFirst {
		for _, fileInfo := range files {
			file, err := os.Open(dir + string(os.PathSeparator) + fileInfo.Name())
			if err != nil {
				panic(err)
			}
			result := AudDAPIRecognizeAll(file, apiToken)
			file.Close()
			if len(result) > 0 {
				Json, _ := json.Marshal(result[0])
				fmt.Printf("%v\n", string(Json))
				return []SongInfo{result[0]}
			}
		}
		return []SongInfo{}
	}
	results := make([][]SongInfo, len(files))
	var mu = &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for i, fileInfo := range files {
		wg.Add(1)
		go func(fileInfo os.FileInfo, i int, dir string) {
			defer wg.Done()
			file, err := os.Open(dir + string(os.PathSeparator) + fileInfo.Name())
			if err != nil {
				panic(err)
			}
			result := AudDAPIRecognizeAll(file, apiToken)
			file.Close()
			mu.Lock()
			results[i] = result
			mu.Unlock()
		}(fileInfo, i, dir)
	}
	wg.Wait()
	first := true
	result := make([]SongInfo, 0)
	for i, response := range results {
		l := len(response)
		for j, song := range response {
			fmt.Printf("%v: %v/%v: %v\n", i, j, l, song)
			if first {
				first = false
				song.Start = secondsPerFile * i
				song.End = song.Start + secondsPerFile/l
				result = append(result, song)
				continue
			}
			if result[len(result)-1].Title == song.Title {
				newR := result[len(result)-1]
				newR.End = secondsPerFile*i + secondsPerFile/l
				result[len(result)-1] = newR
				continue
			}
			song.Start = secondsPerFile*i + (secondsPerFile/l)*j
			song.End = song.Start + secondsPerFile/l
			result = append(result, song)
		}
	}
	return result
}

func SplitVideo(path, tmpDir string, secondsPerFile int) {
	ffmpeg := exec.Command("ffmpeg", "-i", path, "-q:a", "0", "-map", "a", "-f", "segment", "-segment_time", strconv.Itoa(secondsPerFile), "-c", "copy", tmpDir+"/out%03d.aac")
	ffmpeg.Run()
}

func CreateCSV(songs []SongInfo, path string) {
	records := make([][]string, 0)
	for i, song := range songs {
		records = append(records, []string{strconv.Itoa(i + 1), SecondsToTime(song.Start) + "-" + SecondsToTime(song.End),
			song.Title, song.Album, song.Label, song.Artist, song.ReleaseDate, song.TimeCode})
	}
	file, _ := os.Create(path)
	w := csv.NewWriter(file)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func getCurrentDir() string {
	currentFile, _ := os.Executable()
	return filepath.Dir(currentFile)
}

func remove(dir string) {
	os.RemoveAll(dir)
}

func main() {
	secondsPerFileFlag := flag.Int("s", 9, "Seconds per audio file")
	UrlFlag := flag.String("url", "https://www.youtube.com/watch?v=ANEOD16twxo", "Link to the YouTube video")
	untilFirstFlag := flag.Bool("first", false, "Send requests only until first result")
	apiTokenFlag := flag.String("api_token", "test", "AudD API token")
	pathToCSVFlag := flag.String("csv", "audd.csv", "Path to the .csv which will be created")
	flag.Parse()
	secondsPerFile := *secondsPerFileFlag
	Url := *UrlFlag
	untilFirst := *untilFirstFlag
	apiToken := *apiTokenFlag
	pathToCSV := *pathToCSVFlag
	videoFile, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer videoFile.Close()
	fmt.Println("Downloading video...")
	err = DownloadYoutubeVideo(Url, videoFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generating audio files...")
	dir, err := ioutil.TempDir(getCurrentDir(), "temp_")
	if err != nil {
		panic(err)
	}
	SplitVideo(videoFile.Name(), dir, secondsPerFile)
	fmt.Println("Sending files to AudD API...")
	songs := RecognizeMultipleFiles(untilFirst, dir, secondsPerFile, apiToken)
	remove(dir)
	if untilFirst {
		return
	}
	fmt.Println("Removing temp files...")
	fmt.Println("Creating csv...")
	CreateCSV(songs, pathToCSV)
}
