package service

import (
	"context"
	"log"

	"gin/v2/dal"

	"gin/v2/model"
)

type CreateVideoService struct {
	Title string `json:"title" binding:"required"`
	Alias string `json:"alias" binding:""`
	Tags  string `json:"tags" binding:""`
	Desc  string `json:"desc" binding:"required"`
}

func (svc *CreateVideoService) CreateVideo(ctx context.Context) (*model.Video, error) {
	var (
		vd  *model.Video
		err error
	)
	// todo rebuild to use copier
	vd = &model.Video{
		Title: svc.Title,
		Alias: svc.Alias,
		Tags:  svc.Tags,
		Desc:  svc.Desc,
	}

	err = dal.Video.WithContext(ctx).Create(vd)
	if err != nil {
		log.Println("CreateVideo create err", err)
		return nil, err
	}

	return vd, err
}
