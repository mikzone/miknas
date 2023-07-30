package drive

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	getFolderSize "github.com/markthree/go-get-folder-size/src"
	"github.com/mikzone/miknas/server/miknas"
	cp "github.com/otiai10/copy"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FileClientInfo struct {
	Name   string `json:"name"`
	IsFile bool   `json:"isFile"`
	Size   int64  `json:"size"`
	Modify int64  `json:"modify"`
}

func GetFileClientInfo(fullpath string, withFolderSize bool) (*FileClientInfo, error) {
	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}
	isDir := fileInfo.IsDir()
	fsize := fileInfo.Size()
	if isDir {
		fsize = 0
		if withFolderSize {
			size, err := getFolderSize.Parallel(fullpath)
			if err != nil {
				fmt.Printf("GetFileSizeErr: %v", err)
			}
			fsize = size
		}
	}
	clientInfo := FileClientInfo{
		Name:   fileInfo.Name(),
		IsFile: !isDir,
		Size:   fsize,
		Modify: fileInfo.ModTime().Unix(),
	}
	return &clientInfo, nil
}

type inDataFsLocate struct {
	Fsid   string `json:"fsid" form:"fsid" binding:"required"`
	Fspath string `json:"fspath" form:"fspath"`
}

type inDataFsRename struct {
	Fsid   string `json:"fsid" form:"fsid" binding:"required"`
	Fspath string `json:"fspath" form:"fspath" binding:"required"`
	Toname string `json:"toname" binding:"required"`
}

type inDataFsWithDest struct {
	Fsid   string `json:"fsid" form:"fsid" binding:"required"`
	Fspath string `json:"fspath" form:"fspath" binding:"required"`
	Topath string `json:"topath"`
}

type inDataFsThumb struct {
	Fsid    string `json:"fsid" form:"fsid" binding:"required"`
	Fspath  string `json:"fspath" form:"fspath" binding:"required"`
	MaxSize uint   `json:"maxSize" form:"maxSize" binding:"required"`
}

func listFiles(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "r")
	fullpath := fs.MustAbs(loc.Fspath)
	files, err := os.ReadDir(fullpath)
	if err != nil {
		ch.FailResp(err.Error())
		// ch.FailResp("读取目录失败")
		return
	}
	fileInfos := []FileClientInfo{}
	for _, file := range files {
		newfullpath := filepath.Join(fullpath, file.Name())
		fCInfo, _ := GetFileClientInfo(newfullpath, true)
		if fCInfo != nil {
			fileInfos = append(fileInfos, *fCInfo)
		}
	}
	ch.SucResp(gin.H{
		"fspath": loc.Fspath,
		"files":  fileInfos,
	})
}

func queryFileInfo(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "r")
	fullpath := fs.MustAbs(loc.Fspath)
	fCInfo, err := GetFileClientInfo(fullpath, true)
	if err != nil {
		ch.FailResp(err.Error())
		// ch.FailResp("读取目录失败")
		return
	}
	ch.SucResp(*fCInfo)
}

func download(ch *miknas.ContextHelper) {
	c := ch.Ctx
	fsid := c.Param("fsid")
	fspath := c.Param("fspath")
	fs := ch.OpenFs(fsid, "r")
	fullpath := fs.MustAbs(fspath)
	EnsureIsExistedFile(fullpath)
	filename := filepath.Base(fullpath)
	c.FileAttachment(fullpath, filename)
}

func view(ch *miknas.ContextHelper) {
	c := ch.Ctx
	fsid := c.Param("fsid")
	fspath := c.Param("fspath")
	fs := ch.OpenFs(fsid, "r")
	fullpath := fs.MustAbs(fspath)
	EnsureIsExistedFile(fullpath)
	c.File(fullpath)
}

func viewTxt(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "r")
	fullpath := fs.MustAbs(loc.Fspath)
	fInfo := EnsureIsExistedFile(fullpath)
	if fInfo.Size() > 10<<20 { // 10mb
		ch.FailResp("文件太大，无法查看")
		return
	}
	data, err := os.ReadFile(fullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(string(data))
}

func GetSnapshot(videoPath string, frameNum int) (imgData *bytes.Buffer, err error) {
	srcBuf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(srcBuf, os.Stdout).
		Run()

	if err != nil {
		return nil, err
	}
	return srcBuf, nil
}

func fitSize(origWidth, origHeight, toWidth, toHeight int) (int, int) {
	fToWidth := float64(toWidth)
	fOrigWidth := float64(origWidth)
	ratioW := fToWidth / fOrigWidth
	fToHeight := float64(toHeight)
	fOrigHeight := float64(origHeight)
	ratioH := fToHeight / fOrigHeight
	if ratioH >= 1.0 && ratioW >= 1.0 {
		return origWidth, origHeight
	} else if ratioH == ratioW {
		return toWidth, toHeight
	} else if ratioH < ratioW {
		return int(fOrigWidth * ratioH), toHeight
	} else {
		return toWidth, int(fOrigHeight * ratioW)
	}
}

func thumb(ch *miknas.ContextHelper) {
	c := ch.Ctx
	fsid := c.Param("fsid")
	fspath := c.Param("fspath")
	fs := ch.OpenFs(fsid, "r")
	fullpath := fs.MustAbs(fspath)
	EnsureIsExistedFile(fullpath)
	// c.File(fullpath)
	ext := CalcExt(fullpath)
	if ext == "svg" {
		c.File(fullpath)
		return
	}
	fileType := GetFileType(ext)
	var srcBuf []byte
	if fileType == "video" {
		videoBuf, err := GetSnapshot(fullpath, 10)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		srcBuf = videoBuf.Bytes()
	} else if fileType == "img" {
		imgData, err := os.ReadFile(fullpath)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		srcBuf = imgData
	} else {
		c.AbortWithError(400, fmt.Errorf("不支持的文件"))
		return
	}

	newImage, err := doImgThumb(srcBuf, 256)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.Writer.Write(newImage)
}

func genThumb(ch *miknas.ContextHelper) {
	var loc inDataFsThumb
	ch.BindJSON(&loc)
	fsid := loc.Fsid
	fspath := loc.Fspath
	fs := ch.OpenFs(fsid, "r")
	fullpath := fs.MustAbs(fspath)
	EnsureIsExistedFile(fullpath)
	// c.File(fullpath)
	ext := CalcExt(fullpath)
	if ext == "svg" {
		ch.FailResp("不支持压缩svg")
		return
	}
	fileType := GetFileType(ext)
	var srcBuf []byte
	if fileType == "video" {
		videoBuf, err := GetSnapshot(fullpath, 10)
		if err != nil {
			ch.FailResp("获取视频截图发生错误: %s", err.Error())
			return
		}
		srcBuf = videoBuf.Bytes()
	} else if fileType == "img" {
		imgData, err := os.ReadFile(fullpath)
		if err != nil {
			ch.FailResp("生成图像截图发生错误: %s", err.Error())
			return
		}
		srcBuf = imgData
	} else {
		ch.FailResp("不支持的文件格式")
		return
	}

	needSize := int(loc.MaxSize)
	newImage, newFileExt, err := doGenImgThumb(srcBuf, needSize)
	if err != nil {
		ch.FailResp("Resize Error: %s", err.Error())
		return
	}

	oldFileName := filepath.Base(fullpath)
	newFileName := strings.TrimSuffix(oldFileName, filepath.Ext(oldFileName))
	toFolder := filepath.Join(filepath.Dir(fullpath), ".mnthumbs")
	err = os.MkdirAll(toFolder, 0750)
	if err != nil && !os.IsExist(err) {
		ch.FailResp("创建.mnthumbs目录失败: %s", err.Error())
		return
	}
	toPath := filepath.Join(toFolder, newFileName+newFileExt)

	ferr := os.WriteFile(toPath, newImage, 0644)
	if ferr != nil {
		ch.FailResp("创建生成文件失败: %s", ferr.Error())
		return
	}
	ch.SucResp("生成成功")
}

func uploadFiles(ch *miknas.ContextHelper) {
	ch.SetFailRetStatus(400)
	var loc inDataFsLocate
	ch.MustBind(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	EnsureIsExistedDir(fullpath)
	c := ch.Ctx
	form, formErr := c.MultipartForm()
	if formErr != nil {
		ch.FailRespWithStatus(400, formErr.Error())
		return
	}
	files := form.File["files"]
	fileNum := len(files)
	if fileNum < 1 {
		ch.FailRespWithStatus(400, "没有携带要上传的文件")
		return
	} else if fileNum > 1 {
		ch.FailRespWithStatus(400, "不支持同时上传多个文件")
		return
	}

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if filename == "" {
			ch.FailRespWithStatus(400, "文件名不能为空")
			return
		}
		dst := filepath.Join(fullpath, filename)
		EnsureNotExisted(dst)
		// Upload the file to specific dst.
		c.SaveUploadedFile(file, dst)
	}
	ch.SucResp("文件上传成功!")
}

func createFolder(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	EnsureNotExisted(fullpath)
	fullDir := filepath.Dir(fullpath)
	EnsureIsExistedDir(fullDir)
	err := os.Mkdir(fullpath, 0750)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(loc.Fspath)
}

func removeFile(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	fInfo := EnsureExisted(fullpath)
	var err error
	if fInfo.IsDir() {
		err = os.RemoveAll(fullpath)
	} else {
		err = os.Remove(fullpath)
	}
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(loc.Fspath)
}

func renameFile(ch *miknas.ContextHelper) {
	var loc inDataFsRename
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	EnsureExisted(fullpath)
	toname := loc.Toname
	tofullpath := filepath.Join(filepath.Dir(fullpath), toname)
	EnsureNotExisted(tofullpath)
	err := os.Rename(fullpath, tofullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp("重命名成功")
}

func copyFile(ch *miknas.ContextHelper) {
	var loc inDataFsWithDest
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	EnsureExisted(fullpath)
	topath := loc.Topath
	tofulldir := fs.MustAbs(topath)
	EnsureIsExistedDir(tofulldir)
	tofullpath := filepath.Join(tofulldir, filepath.Base(fullpath))
	EnsureNotExisted(tofullpath)
	if strings.HasPrefix(tofullpath, fullpath) {
		ch.FailResp("非法操作:不能将目录复制到子目录")
		return
	}

	err := cp.Copy(fullpath, tofullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp("复制成功")
}

func mvFile(ch *miknas.ContextHelper) {
	var loc inDataFsWithDest
	ch.BindJSON(&loc)
	fs := ch.OpenFs(loc.Fsid, "w")
	fullpath := fs.MustAbs(loc.Fspath)
	EnsureExisted(fullpath)
	topath := loc.Topath
	tofulldir := fs.MustAbs(topath)
	EnsureIsExistedDir(tofulldir)
	tofullpath := filepath.Join(tofulldir, filepath.Base(fullpath))
	EnsureNotExisted(tofullpath)
	if strings.HasPrefix(tofullpath, fullpath) {
		ch.FailResp("非法操作:不能将目录移动到子目录")
		return
	}

	err := os.Rename(fullpath, tofullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp("移动成功")
}

func regRoutes(ext *MikNasExt) {
	ext.POST("/listFiles", listFiles)
	ext.POST("/queryFileInfo", queryFileInfo)
	ext.GET("/download/:fsid/*fspath", download)
	ext.GET("/view/:fsid/*fspath", view)
	ext.GET("/thumb/:fsid/*fspath", miknas.UseReqPool("drive/thumb", 1), thumb)
	ext.POST("/genThumb", miknas.UseReqPool("drive/genThumb", 1), genThumb)
	ext.POST("/viewTxt", viewTxt)
	ext.POST("/uploadFiles", uploadFiles)
	ext.POST("/createFolder", createFolder)
	ext.POST("/removeFile", removeFile)
	ext.POST("/renameFile", renameFile)
	ext.POST("/copyFile", copyFile)
	ext.POST("/mvFile", mvFile)
}
