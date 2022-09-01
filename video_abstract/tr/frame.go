package main

// import (
// 	"crypto/md5"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"strconv"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/pinguo-icc/go-lib/v2/uploader/qiniu"
// 	"github.com/pinguo-icc/video-beats/internal/infrastructure/conf"
// )

// type FrameExtract struct {
// 	rawUploadSvc *qiniu.QiNiu
// 	C            *resty.Client
// 	FI           *FrameInfo
// }

// func NewFrameExtrac(cfg *conf.AppQiNiuConfig) *FrameExtract {
// 	return &FrameExtract{
// 		rawUploadSvc: qiniu.NewQiNiu((*qiniu.Config)(cfg)),
// 		C:            &resty.Client{},
// 	}
// }

// type FrameInfo struct {
// 	Width  int64  `json:"w"`
// 	Height int64  `json:"h"`
// 	Format string `json:"vframe"`
// 	Offset int64  `json:"offset"`
// }

// func (fi *FrameInfo) toMap() map[string]string {
// 	resMap := make(map[string]string)
// 	resMap["vframe"] = fi.Format
// 	resMap["w"] = strconv.Itoa(int(fi.Width))
// 	resMap["h"] = strconv.Itoa(int(fi.Height))
// 	resMap["offset"] = strconv.Itoa(int(fi.Offset))
// 	return resMap
// }

// // func (f *FrameExtract) GetFrameUrl(ctx Context) (*qiniu.UploadResult, error) {
// // 	img, err := f.VideoFrameExtraction(ctx, 1, "jpg", )
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	res, err := f.UploadFrame(VideoBeats, img)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return res, nil
// // }

// func (f *FrameExtract) UploadFrame(prefix string, img io.Reader) (*qiniu.UploadResult, error) {
// 	buf, err := ioutil.ReadAll(img)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hashData := md5.Sum(buf)
// 	key := prefix + "/" + hex.EncodeToString(hashData[:])

// 	uploadRet, err := f.rawUploadSvc.Upload(key, buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return uploadRet, nil
// }

// func (f *FrameExtract) VideoFrameExtraction(ctx Context, offset int64, format string, uploadRet *qiniu.UploadResult) (io.Reader, error) {
// 	url := fmt.Sprintf("%ss", uploadRet.Url)
// 	_ = url
// 	f.SetFrameParam(offset, format, uploadRet)
// 	req, err := http.NewRequest("GET", fmt.Sprintln(""), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Host", "")
// 	// req := f.SetReqParam()
// 	// response, err := req.Get(url)
// 	// if err != nil {,
// 	// 	return nil, err
// 	// }
// 	// reader := response.RawResponse.Body
// 	return nil, nil
// }

// func (fe *FrameExtract) SetReqParam() *resty.Request {
// 	req := fe.C.R()
// 	req.SetPathParams(fe.FI.toMap())
// 	req.SetHeader("Host", "https://cdn-qa-all.c360dn.com")
// 	return req
// }

// func (f *FrameExtract) SetFrameParam(offset int64, format string, uploadRet *qiniu.UploadResult) {
// 	frameParam := &FrameInfo{
// 		Width:  uploadRet.Width,
// 		Height: uploadRet.Height,
// 		Format: format,
// 		Offset: offset,
// 	}
// 	f.FI = frameParam
// }
