package drive

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikzone/miknas/server/miknas"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type inDataFsLocate struct {
	Fsid   string `json:"fsid" form:"fsid" binding:"required"`
	Fspath string `json:"fspath" form:"fspath"`
}

type inDataFsListFile struct {
	Fsid           string `json:"fsid" form:"fsid" binding:"required"`
	Fspath         string `json:"fspath" form:"fspath"`
	NeedFolderSize bool   `json:"needFolderSize" form:"needFolderSize"`
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

type inDataFsUpload struct {
	Fsid   string `json:"fsid" form:"fsid" binding:"required"`
	Fspath string `json:"fspath" form:"fspath"`
	Name   string `json:"name" form:"name" binding:"required"`
}

func listFiles(ch *miknas.ContextHelper) {
	var loc inDataFsListFile
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "r")
	// needDirSize := ch.GetApp().ConfMgr.Get("MIKNAS_DRIVE_NEED_DIR_SIZE").(string) == "1"
	needDirSize := loc.NeedFolderSize
	fileInfos, err := fsd.ListDir(fsd, loc.Fspath, needDirSize)
	ch.EnsureNoErr(err)
	fCInfos := []FileClientInfo{}
	for _, fInfo := range fileInfos {
		fCInfo := PackFileClientInfo(fInfo)
		fCInfos = append(fCInfos, *fCInfo)
	}
	ch.SucResp(gin.H{
		"fspath": loc.Fspath,
		"files":  fCInfos,
	})
}

func queryFileInfo(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "r")
	// 查询文件时不会查询到文件夹大小
	fInfo, err := fsd.Stat(fsd, loc.Fspath, false)
	ch.EnsureNoErr(err)
	fCInfo := PackFileClientInfo(fInfo)
	ch.EnsureNoErr(err)
	ch.SucResp(*fCInfo)
}

func download(ch *miknas.ContextHelper) {
	c := ch.Ctx
	fsid := c.Param("fsid")
	fspath := c.Param("fspath")
	fsd := ch.OpenFs(fsid, "r")
	EnsureIsExistedFile(fsd, fspath)
	fsd.WebDownload(fsd, ch, fspath)
}

func view(ch *miknas.ContextHelper) {
	c := ch.Ctx
	fsid := c.Param("fsid")
	fspath := c.Param("fspath")
	fsd := ch.OpenFs(fsid, "r")
	EnsureIsExistedFile(fsd, fspath)
	fsd.WebDownload(fsd, ch, fspath)
}

func viewTxt(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "r")
	fspath := loc.Fspath
	fInfo := EnsureIsExistedFile(fsd, fspath)
	if fInfo.Size() > 10<<20 { // 10mb
		ch.FailResp("文件太大，无法查看")
		return
	}
	data, err := fsd.ReadFile(fsd, fspath)
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
	fsd := ch.OpenFs(fsid, "r")
	EnsureIsExistedFile(fsd, fspath)
	ext := CalcExt(fspath)
	if ext == "svg" {
		fsd.WebView(fsd, ch, fspath)
		return
	}
	fileType := GetFileType(ext)
	var srcBuf []byte
	if fileType == "video" {
		dfsd, ok := fsd.(miknas.IDiskFsDriver)
		if !ok {
			c.AbortWithError(400, fmt.Errorf("缩略图功能暂不支持该文件系统下的视频文件"))
			return
		}
		diskFullpath := dfsd.MustAbs(fspath)
		videoBuf, err := GetSnapshot(diskFullpath, 10)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		srcBuf = videoBuf.Bytes()
	} else if fileType == "img" {
		imgData, err := fsd.ReadFile(fsd, fspath)
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
	fsd := ch.OpenFs(fsid, "r")
	EnsureIsExistedFile(fsd, fspath)
	ext := CalcExt(fspath)
	if ext == "svg" {
		ch.FailResp("不支持压缩svg")
		return
	}
	fileType := GetFileType(ext)
	var srcBuf []byte
	if fileType == "video" {
		dfsd, ok := fsd.(miknas.IDiskFsDriver)
		if !ok {
			ch.FailResp("缩略图功能暂不支持该文件系统下的视频文件")
			return
		}
		diskFullpath := dfsd.MustAbs(fspath)
		videoBuf, err := GetSnapshot(diskFullpath, 10)
		if err != nil {
			ch.FailResp("获取视频截图发生错误: %s", err.Error())
			return
		}
		srcBuf = videoBuf.Bytes()
	} else if fileType == "img" {
		imgData, err := fsd.ReadFile(fsd, fspath)
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

	oldFileName := filepath.Base(fspath)
	newFileName := strings.TrimSuffix(oldFileName, filepath.Ext(oldFileName))
	toFolder := filepath.Join(filepath.Dir(fspath), ".mnthumbs")
	err = fsd.MkdirAll(fsd, toFolder)
	if err != nil && !os.IsExist(err) {
		ch.FailResp("创建.mnthumbs目录失败: %s", err.Error())
		return
	}
	toPath := filepath.Join(toFolder, newFileName+newFileExt)

	ferr := fsd.WriteFile(fsd, toPath, newImage)
	if ferr != nil {
		ch.FailResp("创建生成文件失败: %s", ferr.Error())
		return
	}
	ch.SucResp("生成成功")
}

func precheckUpload(ch *miknas.ContextHelper) {
	var loc inDataFsUpload
	ch.BindJSON(&loc)
	fsid := loc.Fsid
	fspath := loc.Fspath
	name := loc.Name
	fsd := ch.OpenFs(fsid, "w")
	EnsureIsExistedDir(fsd, fspath)
	filename := filepath.Base(name)
	dst := filepath.Join(fspath, filename)
	EnsureNotExisted(fsd, dst)
	ch.SucResp(true)
}

func uploadFiles(ch *miknas.ContextHelper) {
	ch.SetFailStatus(400)
	var loc inDataFsUpload
	err := ch.Ctx.ShouldBindQuery(&loc)
	ch.EnsureNoErr(err)
	fsid := loc.Fsid
	fspath := loc.Fspath
	name := loc.Name
	fsd := ch.OpenFs(fsid, "w")
	EnsureIsExistedDir(fsd, fspath)
	filename := filepath.Base(name)
	dst := filepath.Join(fspath, filename)
	EnsureNotExisted(fsd, dst)

	c := ch.Ctx
	form, formErr := c.MultipartForm()
	if formErr != nil {
		ch.FailResp(formErr.Error())
		return
	}
	files := form.File["files"]
	fileNum := len(files)
	if fileNum < 1 {
		ch.FailResp("没有携带要上传的文件")
		return
	} else if fileNum > 1 {
		ch.FailResp("不支持同时上传多个文件")
		return
	}

	for _, file := range files {
		if filename != filepath.Base(file.Filename) {
			ch.FailResp("请求错误，文件名和请求参数对不上")
			return
		}
		fsd.WebSaveUploadedFile(fsd, ch, file, dst)
	}
	ch.SucResp("文件上传成功!")
}

func createFolder(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "w")
	fspath := loc.Fspath
	EnsureNotExisted(fsd, fspath)
	fsDir := filepath.Dir(fspath)
	EnsureIsExistedDir(fsd, fsDir)
	err := fsd.MkdirAll(fsd, fspath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(loc.Fspath)
}

func removeFile(ch *miknas.ContextHelper) {
	var loc inDataFsLocate
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "w")
	fspath := loc.Fspath
	err := fsd.RemoveAll(fsd, fspath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp(loc.Fspath)
}

func renameFile(ch *miknas.ContextHelper) {
	var loc inDataFsRename
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "w")
	fspath := loc.Fspath
	EnsureExisted(fsd, fspath)
	toname := loc.Toname
	tofullpath := filepath.Join(filepath.Dir(fspath), toname)
	EnsureNotExisted(fsd, tofullpath)
	err := fsd.MoveAll(fsd, fspath, tofullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp("重命名成功")
}

func copyFile(ch *miknas.ContextHelper) {
	var loc inDataFsWithDest
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "w")
	fspath := loc.Fspath
	EnsureExisted(fsd, fspath)
	topath := loc.Topath
	EnsureIsExistedDir(fsd, topath)
	tofullpath := filepath.Join(topath, filepath.Base(fspath))
	EnsureNotExisted(fsd, tofullpath)
	if strings.HasPrefix(tofullpath, fspath) {
		ch.FailResp("非法操作:不能将目录复制到子目录")
		return
	}

	err := fsd.CopyAll(fsd, fspath, tofullpath)
	if err != nil {
		ch.FailResp(err.Error())
		return
	}
	ch.SucResp("复制成功")
}

func mvFile(ch *miknas.ContextHelper) {
	var loc inDataFsWithDest
	ch.BindJSON(&loc)
	fsd := ch.OpenFs(loc.Fsid, "w")
	fspath := loc.Fspath
	EnsureExisted(fsd, fspath)
	topath := loc.Topath
	EnsureIsExistedDir(fsd, topath)
	tofullpath := filepath.Join(topath, filepath.Base(fspath))
	EnsureNotExisted(fsd, tofullpath)
	if strings.HasPrefix(tofullpath, fspath) {
		ch.FailResp("非法操作:不能将目录移动到子目录")
		return
	}

	err := fsd.MoveAll(fsd, fspath, tofullpath)
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
	ext.POST("/precheckUpload", precheckUpload)
	ext.POST("/uploadFiles", uploadFiles)
	ext.POST("/createFolder", createFolder)
	ext.POST("/removeFile", removeFile)
	ext.POST("/renameFile", renameFile)
	ext.POST("/copyFile", copyFile)
	ext.POST("/mvFile", mvFile)
}
