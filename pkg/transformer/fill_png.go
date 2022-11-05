package transformer

import (
	"bytes"
	"image/png"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

func NewPng() TransformInterface {
	return &pngImpl{}
}

type pngImpl struct{}

func (j *pngImpl) FillCenter(source []byte, width, height int) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	dst := imaging.Fill(src, width, height, imaging.Center, imaging.Box)

	var buff bytes.Buffer
	err = png.Encode(&buff, dst)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (j *pngImpl) IsSupported(source []byte) bool {
	return strings.Contains(http.DetectContentType(source), "image/png")
}
