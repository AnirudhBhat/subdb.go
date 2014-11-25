package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
)

/*

	function get_hash takes video path as input and returns
	hash of the video file calculated by taking first and last
	64kb of video file.


*/

func GetHash(name string) string {
	readsize := 64 * 1024
	// open file
	f, err := os.Open(name)
	if err != nil {
		fmt.Println("error")
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Println("error")
	}
	size := fi.Size()
	buf := make([]byte, readsize)
	buf1 := make([]byte, readsize)
	for {
		// read a chunk
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("error")
		}
		if n == 0 {
			break
		}
		f.Seek(size-65536, 0)
		_, err = f.Read(buf1)
		buffer := append(buf, buf1...)
		hasher := md5.New()
		hasher.Write([]byte(buffer))
		//fmt.Println(hex.EncodeToString(hasher.Sum(nil)))
		return hex.EncodeToString(hasher.Sum(nil))
	}
	return " "
}

func SubDownloader(video_path string) {
	hash := GetHash(video_path)
	url := "http://api.thesubdb.com/?action=download&hash=" + hash + "&language=en"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		req.Header.Set("User-Agent", "SubDB/1.0 (SubDownloader/0.1; http://github.com/AnirudhBhat)")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		f, err := os.Create(path.Dir(video_path) + "/subtitle.srt")
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error")
		}
		f.Write(body)
	}
}

/*
	function notify notifies once the subtitle
	is downloaded.

*/
func notify(path string) {
	command := "notify-send"
	message := "subtitle for " + path + " downloaded!"
	cmd := exec.Command(command, message)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	movie_path := os.Args[1]
	SubDownloader(movie_path)
	notify(movie_path)
}
