/*
 * Copyright (c) 2023 AudD, LLC. All rights reserved.
 * Copyright (c) 2023 Mikhail Samin. All rights reserved.
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
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/AudDMusic/audd-go"
	"log"
	"os"
	"strconv"
)

func CreateCSV(songs []audd.RecognitionEnterpriseResult, path string) {
	records := make([][]string, 0)
	i := 0
	for _, result := range songs {
		for _, song := range result.Songs {
			i++
			records = append(records, []string{strconv.Itoa(i), result.Offset,
				song.Title, song.Album, song.Label, song.Artist, song.ReleaseDate, song.Timecode})
		}
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

func main() {
	UrlFlag := flag.String("url", "https://www.youtube.com/watch?v=ZgRuRcQqzTo", "Link to the YouTube video")
	apiTokenFlag := flag.String("api_token", "test", "AudD API token")
	pathToCSVFlag := flag.String("csv", "audd.csv", "Path to the .csv which will be created")
	skipFlag := flag.Int("skip", 0, "Skip audio files")
	everyFlag := flag.Int("every", 1, "Audio files between skips")
	flag.Parse()
	Url := *UrlFlag
	apiToken := *apiTokenFlag
	pathToCSV := *pathToCSVFlag
	skip := *skipFlag
	every := *everyFlag
	fmt.Println("Sending a request to the AudD API...")
	client := audd.NewClient(apiToken)
	client.SetEndpoint(audd.EnterpriseAPIEndpoint)
	parameters := map[string]string{"skip": strconv.Itoa(skip), "every": strconv.Itoa(every)}
	songs, err := client.RecognizeLongAudio(Url, parameters)
	if err != nil {
		panic(err)
	}
	fmt.Println("Creating csv...")
	CreateCSV(songs, pathToCSV)
}
