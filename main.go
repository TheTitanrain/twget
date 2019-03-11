// https://dev.twitch.tv/docs/api/reference/#get-videos
// curl -H "Client-ID: vmseauicwir0sq9vojejn18btclbw9" -X GET "https://api.twitch.tv/helix/users?login=stopgameru"
// curl -H "Client-ID: vmseauicwir0sq9vojejn18btclbw9" -X GET "https://api.twitch.tv/helix/videos?id=287782398"
// ffmpeg -y -i "file:Кинологи. Битва у Красной Скалы Ричарда II и Чёрного клановца.mp4" -vn -acodec libmp3lame -q:a 6 "file:Кинологи. Битва у Красной Скалы Ричарда II и Чёрного клановца.mp3"

package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	baseUrl      = viper.GetString("baseUrl")
	clientID     = viper.GetString("clientID")
	keywords     = viper.GetString("keywords")
	period       = viper.GetString("period")
	first        = viper.GetString("first")
	argsFilename = viper.GetString("argsFilename")
	argsEncode   = viper.GetString("argsEncode")
	argsMp3gain  = viper.GetString("argsMp3gain")
	userName     = viper.GetString("userName")
)

func main() {
	getConfig()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	id := getUserId(userName)
	json := getVideos(id)

	//id := gjson.Get(json, "data.#.id")
	titles := gjson.Get(json, "data.#.title")
	//urls := gjson.Get(json, "data.#.url")
	//durations := gjson.Get(json, "data.#.duration")
	//created :=gjson.get(json, "data.#.created_at")
	keywords = viper.GetString("keywords")
	for i, title := range titles.Array() {
		//println(i)
		for _, keyword := range strings.Split(keywords, ";") { // перебираем все ключевые слова
			keyword = strings.Join(strings.Fields(keyword), " ") // убираем лишние пробелы
			if strings.Contains(title.Str, keyword) {
				println(title.String())
				s := strconv.Itoa(i)
				url := gjson.Parse(json).
					Get("data").
					Get(s).
					Get("url")
				idvideo := gjson.Parse(json).
					Get("data").
					Get(s).
					Get("id")
				println(url.Str)

				argsFilename = dir + argsFilename + url.Str
				println("Command to get file name: " + argsFilename)
				filenameext := execute(argsFilename, "") // getting file name with ext
				println("Filename with ext: " + filenameext)
				filename := filenameext[0 : len(filenameext)-5]
				println("Filename: " + filename)

				argsEncode = dir + argsEncode + url.Str
				println("Command to download: " + argsEncode)
				println("Downloading...")
				execute(argsEncode, "")

				argsFFMPEG := dir + "/tools/ffmpeg.exe " +
					"-y " +
					"-i file:v" + idvideo.Str + ".mp4 " +
					"-vn " +
					"-acodec libmp3lame " +
					"-q:a 5 " +
					"-ac 1 " +
					"file:" + idvideo.Str + ".mp3"
				println("Command for FFMPEG: " + argsFFMPEG)
				println("Encoding...")
				execute(argsFFMPEG, "")

				argsMp3gain = dir + argsMp3gain
				println("Command to MP3Gain: " + argsMp3gain) // normalize mp3 audio
				println("Normalizing...")
				execute(argsMp3gain, idvideo.Str+".mp3")

				oldFileName := idvideo.Str + ".mp3"
				newFileName := filename + ".mp3"
				os.Rename(oldFileName, newFileName)

				fmt.Println("Done! Press Enter to exit!")
				fmt.Scanln()
			}
		}
	}
}

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	baseUrl = viper.GetString("baseUrl")
	clientID = viper.GetString("clientID")
	keywords = viper.GetString("keywords")
	period = viper.GetString("period")
	first = viper.GetString("first")
	argsFilename = viper.GetString("argsFilename")
	argsEncode = viper.GetString("argsEncode")
	argsMp3gain = viper.GetString("argsMp3gain")
	userName = viper.GetString("userName")
}

func getUserId(username string) string {
	url := baseUrl + "users?login=" + username
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", clientID)
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	if response.StatusCode == 404 {
		fmt.Printf("Error: %s", "404. Page not found.")
		os.Exit(1)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	println(string(contents))
	json := string(contents)
	id := gjson.Get(json, "data.0.id")
	println(id.String())
	return id.String()
}

func getVideos(id string) string {
	url := baseUrl + "videos?user_id=" + id + "&first=" + first + "&period=" + period
	println(url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", clientID)
	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	if response.StatusCode == 404 {
		fmt.Printf("Error: %s", "404. Page not found.")
		os.Exit(1)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	println(string(contents))
	json := string(contents)

	return json

}

func execute(args string, fname string) string {
	parts := strings.Fields(args)
	head := parts[0]
	parts = parts[1:]
	if fname != "" {
		parts = append(parts, fname)
	}

	cmd := exec.Command(head, parts...)
	var buf bytes.Buffer
	var buferr bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &buferr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		//os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%s finished with error: %v\n", head, err.Error())
		fmt.Println(win1251toUtf8(buferr.String()))
		os.Exit(1)
	}

	sutf8 := win1251toUtf8(buf.String())

	fmt.Printf("%s finished with output: %v\n", head, sutf8)
	return sutf8

}

func win1251toUtf8(buf string) string { ///  Convert Windows-1251 to UTF8
	s := buf
	sr := strings.NewReader(s)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	st, err := ioutil.ReadAll(tr)
	if err != nil {
		log.Fatal(err)
	}
	sutf8 := string(st)
	return sutf8

}
