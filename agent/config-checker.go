// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// BEFORE RUNNING:
// ---------------
// 0. This script runs gapic-config-validator on all APIs.
//    Install gapic-config-validator with the following command:
//      go get -u github.com/googleapis/gapic-config-validator/cmd/protoc-gen-gapic-validator
// 1. If not already done, enable the Google Sheets API
//    and check the quota for your project at
//    https://console.developers.google.com/apis/api/sheets
// 2. Install and update the Go dependencies by running `go get -u` in the
//    project directory.
// 3. Set the SHEETID to the spreadsheet you wish to update.
// 4. Create a service account (or use an existing one) and share
//    your spreadsheet with that service account.
// 5. Be sure that the GOOGLE_APPLICATION_CREDENTIALS environment variable
//    points to the file containing your service account credentials.

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const SHEETID = ""

type CheckerResult struct {
	api      string
	command  string
	messages []string
}

// Run gapic-config-validator plugin for every API in googleapis.
// APIs are recognized by the presence of a gapic.yaml file.
// If a gapic.legacy.yaml file is present, prefer that.
func checkConfigurations() []CheckerResult {
	checkerResults := make([]CheckerResult, 0)
	// find all gapic.yaml files
	cmd0 := exec.Command("find", ".", "-name", "*_gapic.yaml")
	out0, err := cmd0.CombinedOutput()
	gapic_yaml_files := strings.Split(string(out0), "\n")
	if err != nil {
		log.Fatal(err)
	}
	// find all gapic.legacy.yaml files
	cmd1 := exec.Command("find", ".", "-name", "*_gapic.legacy.yaml")
	out1, err := cmd1.CombinedOutput()
	gapic_legacy_yaml_files := make(map[string]bool, 0)
	for _, f := range strings.Split(string(out1), "\n") {
		if f != "" {
			gapic_legacy_yaml_files[f] = true
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range gapic_yaml_files {
		// skip blank lines
		if len(f) == 0 {
			continue
		}
		// skip language-specific yaml files
		if strings.HasPrefix(f, "./gapic/lang") {
			continue
		}
		// If there is a gapic.legacy.yaml file, use that instead.
		// ALERT! This assumes there is always a gapic.yaml (v2) file after the conversion.
		legacy := strings.Replace(f, "gapic.yaml", "gapic.legacy.yaml", 1)
		if gapic_legacy_yaml_files[legacy] {
			f = legacy
		}
		// use the gapic yaml to run the config checker...
		api := strings.Replace(filepath.Dir(f), "./", "", -1)
		command := strings.Join([]string{
			"/usr/local/bin/protoc",
			filepath.Dir(f) + "/*.proto",
			"--gapic-validator_out=.",
			"--gapic-validator_opt='gapic-yaml=" + f + "'"},
			" ")
		out, err := exec.Command("sh", "-c", command).CombinedOutput()
		checkerResult := CheckerResult{api: api, command: command, messages: nil}
		if err != nil {
			s := string(out)
			lines := strings.Split(s, "\n")
			messages := make([]string, 0)
			for _, line := range lines {
				if strings.Contains(line, "warning: ") {
					continue
				}
				if line == "--gapic-validator_out: " {
					continue
				}
				messages = append(messages, line)
			}
			checkerResult.messages = messages
		}
		checkerResults = append(checkerResults, checkerResult)
	}
	return checkerResults
}

type StatusSheetConnection struct {
	spreadsheetId string
	sheetsService *sheets.Service
}

func NewStatusSheetConnection(id string) (*StatusSheetConnection, error) {
	ssc := &StatusSheetConnection{spreadsheetId: id}
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, sheets.DriveScope, sheets.DriveFileScope, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}
	ssc.sheetsService, err = sheets.New(client)
	if err != nil {
		return nil, err
	}
	return ssc, err
}

func (ssc *StatusSheetConnection) fetch() (*sheets.ValueRange, error) {
	ctx := context.Background()
	return ssc.sheetsService.Spreadsheets.Values.Get(
		ssc.spreadsheetId,
		"Summary!A:A").
		ValueRenderOption("FORMATTED_VALUE").
		DateTimeRenderOption("SERIAL_NUMBER").
		MajorDimension("ROWS").
		Context(ctx).
		Do()
}

func (ssc *StatusSheetConnection) updateRows(cellRange string, values [][]interface{}) {
	ctx := context.Background()
	valueRange := &sheets.ValueRange{
		Range:          cellRange,
		MajorDimension: "ROWS",
		Values:         values,
	}
	_, err := ssc.sheetsService.Spreadsheets.Values.Update(
		ssc.spreadsheetId, valueRange.Range, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		log.Fatal(err)
	}
}

func matching(checkerResults []CheckerResult, api string) *CheckerResult {
	for _, checkerResult := range checkerResults {
		if checkerResult.api == api {
			return &checkerResult
		}
	}
	return nil
}

func rowForResult(checkerResult *CheckerResult, withDetails bool) []interface{} {
	row := make([]interface{}, 0)
	row = append(row, (*checkerResult).api)
	if withDetails {
		row = append(row, (*checkerResult).command)
		row = append(row, strings.Join((*checkerResult).messages, "\n"))
	} else {
		row = append(row, len((*checkerResult).messages))
	}
	return row
}

func (ssc *StatusSheetConnection) updateWithCheckerResults(checkerResults []CheckerResult) {
	// reset headings
	ssc.updateRows("Summary!A1:B1",
		[]([]interface{}){[]interface{}{"API", "Errors"}},
	)
	ssc.updateRows("Details!A1:C1",
		[]([]interface{}){[]interface{}{"API", "Command", "Messages"}},
	)

	// get current sheet contents
	contents, _ := ssc.fetch()

	// prepare update
	summary := make([][]interface{}, 0)
	details := make([][]interface{}, 0)
	seen := make(map[string]bool, 0)

	// for each API in the sheet...
	for i, v := range contents.Values {
		if i == 0 {
			continue // skip header
		}
		// get the analysis results
		if len(v) > 0 {
			api := v[0].(string)
			checkerResult := matching(checkerResults, api)
			if checkerResult != nil {
				// if we have it, add it to the table
				summary = append(summary, rowForResult(checkerResult, false))
				details = append(details, rowForResult(checkerResult, true))
			} else {
				// if we don't have it, mark the api as unknown
				summary = append(summary, []interface{}{api, "unknown"})
				details = append(details, []interface{}{api, "unknown", "unknown"})
			}
			seen[api] = true
		} else {
			// mark any weird table rows
			summary = append(summary, []interface{}{"unknown", "unknown"})
			details = append(details, []interface{}{"unknown", "unknown", "unknown"})
		}
	}
	// now go back through the checker results and add any that weren't in the table
	for _, checkerResult := range checkerResults {
		if !seen[checkerResult.api] {
			summary = append(summary, rowForResult(&checkerResult, false))
			details = append(details, rowForResult(&checkerResult, true))
		}
	}
	ssc.updateRows(fmt.Sprintf("Summary!A2:B%d", 1+len(summary)), summary)
	ssc.updateRows(fmt.Sprintf("Details!A2:C%d", 1+len(details)), details)
}

func run() {
	dir, err := ioutil.TempDir(".", "temp")
	if err != nil {
		log.Printf("error making directory: %+v", err)
		return
	}
	defer os.RemoveAll(dir)
	log.Printf("Running checks on APIs in %s.", dir)

	err = os.Chdir(dir)
	{
		command := "curl -L -O https://github.com/googleapis/googleapis/archive/master.zip"
		out, err := exec.Command("sh", "-c", command).CombinedOutput()
		log.Printf("curl results %s (err=%v)", out, err)
	}
	{
		command := "unzip master.zip"
		_, err := exec.Command("sh", "-c", command).CombinedOutput()
		log.Printf("unzip results (err=%v)", err)
	}
	{
		command := "mv googleapis-master googleapis"
		out, err := exec.Command("sh", "-c", command).CombinedOutput()
		log.Printf("mv results %s (err=%v)", out, err)
	}

	err = os.Chdir("googleapis")

	checkerResults := checkConfigurations()
	log.Printf("checker results %+v", checkerResults)
	log.Printf("Updating sheet.")
	ssc, err := NewStatusSheetConnection(SHEETID)
	if err != nil {
		log.Printf("connection to spreadsheet API failed with error %+v", err)
	} else {
		ssc.updateWithCheckerResults(checkerResults)
	}
	log.Printf("Done.")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Config checker received a request.")
	run()
	fmt.Fprintf(w, "checker completed\n")
}

func main() {
	log.Print("config-checker started.")

	http.HandleFunc("/check", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
