package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/88250/gulu"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var logger = gulu.Log.NewLogger(os.Stdout)

func fetchPaperMCVersionGroups() []interface{} {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://papermc.io/api/v2/projects/paper/").Timeout(7*time.Second).
		Set("User-Agent", "Profile Bot; +https://github.com/OhMyMC/PaperMC-Docker").EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Fatalf("fetch events failed: %+v, %s", errors, data)
		panic("fetch events failed")
	}
	return result["versions"].([]interface{})
}

func fetchPaperMCBuild(version string) interface{} {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://papermc.io/api/v2/projects/paper/versions/"+version).Timeout(7*time.Second).
		Set("User-Agent", "Profile Bot; +https://github.com/OhMyMC/PaperMC-Docker").EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Fatalf("fetch events failed: %+v, %s", errors, data)
		panic("fetch events failed")
	}
	builds := result["builds"].([]interface{})
	return builds[len(builds)-1].(interface{})
}

func readPaperMCBuildJson(path string) map[string]int64 {
	file, err := ioutil.ReadFile(path)
	if nil != err {
		logger.Fatalf("open file failed: %s", err)
		panic("open file failed")
	}
	var buildJson map[string]int64
	err = json.Unmarshal(file, &buildJson)
	if nil != err {
		logger.Fatalf("parse json failed: %s", err)
		panic("parse json failed")
	}
	return buildJson
}

func main() {
	paperMCJsonPath := "papermc-data/papermc-build.json"
	paperMCCIEnv := ".env.tmp"

	paperMCJson := readPaperMCBuildJson(paperMCJsonPath)
	versions := fetchPaperMCVersionGroups()
	PaperBuildVersion := ""
	PaperBuildNumber := int64(0)
	for _, event := range versions {
		version := event.(string)
		build := fetchPaperMCBuild(version).(float64)
		logger.Infof("version: %s, build: %f", version, build)
		if paperMCJson[version] != int64(build) {
			// update .env file
			paperMCJson[version] = int64(build)
			PaperBuildVersion = version
			PaperBuildNumber = int64(build)
			break
		}
		time.Sleep(1 * time.Second)
	}

	bytes, err := json.Marshal(paperMCJson)
	if nil != err {
		logger.Fatalf("marshal failed: %+v", err)
		panic(err)
	}
	if PaperBuildNumber != 0 && PaperBuildVersion != "" {
		bytes := []byte("PAPER_BUILD_VERSION=" + PaperBuildVersion + "\n" + "PAPER_BUILD_NUMBER=" + strconv.FormatInt(PaperBuildNumber, 10))
		if err := ioutil.WriteFile(paperMCCIEnv, bytes, 0644); nil != err {
			logger.Fatalf("write papermc-build.json failed: %s", bytes)
		}
	}
	if err := ioutil.WriteFile(paperMCJsonPath, bytes, 0644); nil != err {
		logger.Fatalf("write papermc-build.json failed: %s", bytes)
	}
}
