package transformer

import "errors"

type TransformInterface interface {
	FillCenter(source []byte, width, height int) ([]byte, error)
	IsSupported(source []byte) bool
}

type stack struct {
	transforms []TransformInterface
}

var ErrFileNotSupported = errors.New("file not supported")

func NewStack() TransformInterface {
	return NewStackBy(NewJpeg(), NewPng())
}

func NewStackBy(transforms ...TransformInterface) TransformInterface {
	return &stack{transforms: transforms}
}

func (s *stack) FillCenter(source []byte, width, height int) ([]byte, error) {
	for _, transform := range s.transforms {
		if transform.IsSupported(source) {
			return transform.FillCenter(source, width, height)
		}
	}

	return nil, ErrFileNotSupported
}

func (s *stack) IsSupported(source []byte) bool {
	for _, transform := range s.transforms {
		if transform.IsSupported(source) {
			return true
		}
	}
	return false
}
