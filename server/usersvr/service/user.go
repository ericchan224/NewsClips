package service

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ericchan224/NewsClips/server/pkg/pb"
	"usersvr/log"
	"usersvr/repository"
)

var (
	Secret = []byte("TikTok")
	// TokenExpireDuration = time.Hour * 2 过期时间
)

type JWTClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

// UpdateUserFavoritedCount 更新用户 获赞数，ActionType 1：表示+1 2：-1
func (u UserService) UpdateUserFavoritedCount(ctx context.Context, req *pb.UpdateUserFavoritedCountReq) (*pb.UpdateUserFavoritedCountRsp, error) {
	err := repository.UpdateUserFavoritedNum(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateUserFavoritedCount err", req.UserId)
		return nil, err
	}
	return &pb.UpdateUserFavoritedCountRsp{}, nil
}