package main

import (
	"context"
	"errors"
	"github.com/ncruces/zenity"
)

type Backend struct {
}

var backend = Backend{}

type Empty struct {
}

type BackendRequest struct {
	Paths []string
}

type PlayingRequest struct {
	IsPlay bool `json:"is_play"`
}
type HelloResponse struct {
	Message string `json:"message"`
}

func (b *Backend) Hello(ctx context.Context, request *Empty) (*HelloResponse, error) {
	return &HelloResponse{Message: "hello"}, nil
}

func (b *Backend) Playing(ctx context.Context, request *PlayingRequest) (*Empty, error) {
	PlayerList.Playing(request.IsPlay)
	return &Empty{}, nil
}

func (b *Backend) Next(ctx context.Context, request *Empty) (*Empty, error) {
	PlayerList.Next()
	return &Empty{}, nil
}

func (b *Backend) Prev(ctx context.Context, request *Empty) (*Empty, error) {
	PlayerList.Prev()
	return &Empty{}, nil
}

func (b *Backend) Mode(ctx context.Context, request *Empty) (*Empty, error) {
	PlayerList.Player.Mode = (PlayerList.Player.Mode + 1) % 4
	return &Empty{}, nil
}

func (b *Backend) FileMulti(ctx context.Context, request *Empty) (*Empty, error) {
	files, err := zenity.SelectFileMultiple(
		zenity.FileFilters{
			{"选择mp3文件", []string{"*.mp3"}},
		})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return &Empty{}, err
	} else if len(files) < 1 {
		return &Empty{}, errors.New("selected files is empty")
	}
	err = PlayerList.PlayList(files)
	if err != nil {
		return nil, err
	}
	PlayerList.Play(0)
	return &Empty{}, nil
}
func RangeConvert(value float64, s1 float64, s2 float64, d1 float64, d2 float64) float64 {
	w1 := s2 - s1
	w2 := d2 - d1
	return (value+s1)/w1*w2 + d1
}
