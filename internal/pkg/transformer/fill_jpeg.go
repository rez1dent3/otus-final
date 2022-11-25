package transformer

import (
	"bytes"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

func NewJpeg() TransformInterface {
	return &jpegImpl{}
}

type jpegImpl struct{}

func (j *jpegImpl) FillCenter(source []byte, width, height int) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	dst := imaging.Fill(src, width, height, imaging.Center, imaging.Box)

	var buff bytes.Buffer
	err = jpeg.Encode(&buff, dst, nil)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (j *jpegImpl) IsSupported(source []byte) bool {
	return strings.Contains(http.DetectContentType(source), "image/jpeg")
}
