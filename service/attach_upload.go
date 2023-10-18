package service

import (
	"HiChat/common"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func UploadImage(ctx *gin.Context) {
	w := ctx.Writer
	req := ctx.Request

	// 從請求中獲取文件
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		common.RespFail(w, err.Error())
		return
	}

	//檢查文件後綴
	suffixText := ".png"
	imageAllowType := []string{"png", "jpg", "jpeg"}
	fileName := head.Filename
	tempSuffixText := strings.Split(fileName, ".")

	if !slices.Contains(imageAllowType, tempSuffixText[len(tempSuffixText)-1]){
		common.RespFail(w, "文件格式錯誤")
		return	
	}
	if len(tempSuffixText) > 1 {
		suffixText = "." + tempSuffixText[len(tempSuffixText)-1]
	}

	// 保存文件
	newFileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffixText)
	dstFile, err := os.Create("./asset/upload/" + newFileName)
	if err != nil {
		common.RespFail(w, err.Error())
		return
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		common.RespFail(w, err.Error())
		return
	}

	url := "./asset/upload/" + newFileName
	common.RespOK(w, url, "上傳成功")
}
