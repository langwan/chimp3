package main

import (
	"context"
	"fmt"
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

var chiStreamer = ChiStreamer{}

func (b *Backend) UpdateList(ctx context.Context, request *BackendRequest) (*Empty, error) {

	err := Player.PlayList(request.Paths)
	if err != nil {
		return &Empty{}, err
	}

	Player.Play(0)

	return &Empty{}, nil
}

func (b *Backend) Playing(ctx context.Context, request *PlayingRequest) (*Empty, error) {
	fmt.Println("playing request", request.IsPlay)
	Player.Playing(request.IsPlay)
	return &Empty{}, nil
}

func (b *Backend) Next(ctx context.Context, request *Empty) (*Empty, error) {
	Player.Next()
	return &Empty{}, nil
}

func (b *Backend) Prev(ctx context.Context, request *Empty) (*Empty, error) {
	Player.Prev()
	return &Empty{}, nil
}
