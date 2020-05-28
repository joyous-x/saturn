package bizs

import (
	"fmt"
	"io"
	"io/ioutil"
	"krotas/common/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
	xerr "github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
)

// Tran2Cartoon translate a picture to cartoon style
func Tran2Cartoon(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		reqresp.ResponseMarshal(c, errcode.ErrNoFileFound, nil)
		return
	}

	if header.Size > 10*1024*1024 {
		reqresp.ResponseMarshal(c, errcode.ErrTooBigFile, nil)
		return
	}
	filename := header.Filename
	contentType := c.Request.Header.Get("Content-Type")
	uniqueID := c.Request.Header.Get("unique_id")
	if len(uniqueID) != 32 {
		reqresp.ResponseMarshal(c, errcode.ErrInvalidUniqueId, nil)
		return
	}

	newContent, err := tran2Cartoon(file)
	if err != nil {
		reqresp.ResponseMarshal(c, xerr.NewError(xerr.ErrServerError.Code, err.Error()), nil)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	c.Header("Content-Type", func() string {
		if len(contentType) > 0 {
			return contentType
		} else {
			return "application/octet-stream"
		}
	}())
	c.Header("Accept-Length", fmt.Sprintf("%d", len(newContent)))
	c.Writer.Write([]byte(newContent))
}

func tran2Cartoon(raw io.Reader) ([]byte, error) {
	rawContent, err := ioutil.ReadAll(raw)
	if err != nil {
		return nil, err
	}

	//> TODO

	return rawContent, nil
}
