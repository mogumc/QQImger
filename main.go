package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type EmojiResponse struct {
	Author        string `json:"author"`
	ChildEmojiId  string `json:"childEmojiId"`
	CommDiyText   []any  `json:"commDiyText"`
	DownloadCount int    `json:"downloadcount"`
	FeetType      string `json:"feetype"`
	FileSize      string `json:"filesize"`
	ID            string `json:"id"`
	Imgs          []struct {
		DiyText        []string `json:"diyText"`
		ID             string   `json:"id"`
		Keywords       []string `json:"keywords"`
		Name           string   `json:"name"`
		Param          string   `json:"param"`
		WHeightInPhone int      `json:"wHeightInPhone"`
		WWidthInPhone  int      `json:"wWidthInPhone"`
	} `json:"imgs"`
	IsApng        int    `json:"isApng"`
	IsOriginal    int    `json:"isOriginal"`
	Mark          string `json:"mark"`
	Name          string `json:"name"`
	OperationInfo []struct {
		MaxVersion string `json:"maxVersion"`
		MinVersion string `json:"minVersion"`
		Platform   int    `json:"platform"`
	} `json:"operationInfo"`
	Price           int    `json:"price"`
	Rights          string `json:"rights"`
	RingType        string `json:"ringtype"`
	Status          string `json:"status"`
	SupportApngSize []any  `json:"supportApngSize"`
	SupportSize     []struct {
		Height int `json:"Height"`
		Width  int `json:"Width"`
	} `json:"supportSize"`
	Type       int    `json:"type"`
	UpdateTime int64  `json:"updateTime"`
	ValidArea  string `json:"validArea"`
}

func downloadImage(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func getimgs(emojiID string) {
	lastChar := emojiID[len(emojiID)-1:]

	url := fmt.Sprintf("https://i.gtimg.cn/club/item/parcel/%s/%s_android.json", lastChar, emojiID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	var emojiResponse EmojiResponse
	if err := json.Unmarshal(body, &emojiResponse); err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return
	}

	downloadDir := filepath.Join("download", emojiID)
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}

	for _, img := range emojiResponse.Imgs {
		imageID := img.ID
		imageURL := fmt.Sprintf("https://i.gtimg.cn/club/item/parcel/item/%s/%s/%dx%d.png", imageID[:2], imageID, img.WHeightInPhone, img.WWidthInPhone)
		imagePath := filepath.Join(downloadDir, fmt.Sprintf("%s.png", imageID))

		if err := downloadImage(imageURL, imagePath); err != nil {
			fmt.Printf("下载图片 %s 失败: %v\n", imageID, err)
		} else {
			fmt.Printf("成功下载图片 %s 到 %s\n", imageID, imagePath)
		}
	}
}

func main() {
	for {
		var emojiID string
		fmt.Print("请输入表情包ID: ")
		fmt.Scanln(&emojiID)
		getimgs(emojiID)
	}

}
