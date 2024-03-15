package api

import (
	"context"
	"fmt"
	awsx "github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/app/config"
	routerContext "github.com/VerzCar/vyf-user/app/router/ctx"
	"github.com/VerzCar/vyf-user/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type UserUploadService interface {
	UploadImage(
		ctx context.Context,
		multiPartFile *multipart.FileHeader,
	) (string, error)
}

type userUploadService struct {
	userService       UserService
	extStorageService awsx.S3Service
	config            *config.Config
	log               logger.Logger
}

func NewUserUploadService(
	userService UserService,
	extStorageService awsx.S3Service,
	config *config.Config,
	log logger.Logger,
) UserUploadService {
	return &userUploadService{
		userService:       userService,
		extStorageService: extStorageService,
		config:            config,
		log:               log,
	}
}

func (u *userUploadService) UploadImage(
	ctx context.Context,
	multiPartFile *multipart.FileHeader,
) (string, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting auth claims: %s", err)
		return "", err
	}

	contentFile, err := multiPartFile.Open()

	if err != nil {
		u.log.Errorf("error opening multipart file: %s", err)
		return "", err
	}

	defer contentFile.Close()

	err = u.detectMimeType(contentFile)

	if err != nil {
		return "", err
	}

	decodedImage, _, err := image.Decode(contentFile)

	size, calculated := utils.CalculatedImageSize(decodedImage, image.Point{X: 600, Y: 400})

	userImage := decodedImage

	if calculated {
		userImage = utils.ResizeImage(decodedImage, size.Max)
	}

	tempImageFile, err := os.CreateTemp("", "user-image")

	if err != nil {
		u.log.Errorf("error creating temp file: %s", err)
		return "", err
	}

	defer os.Remove(tempImageFile.Name())

	err = png.Encode(tempImageFile, userImage)

	if err != nil {
		u.log.Errorf("error encoding image file to png: %s", err)
		return "", err
	}

	_, _ = tempImageFile.Seek(0, 0)

	filePath := fmt.Sprintf("profile/image/%s/%s", authClaims.Subject, "avatar.png")

	_, err = u.extStorageService.Upload(
		ctx,
		filePath,
		tempImageFile,
	)

	if err != nil {
		u.log.Errorf("error uploading file to external storage service: %s", err)
		return "", err
	}

	imageEndpoint := fmt.Sprintf("%s/%s", u.extStorageService.ObjectEndpoint(), filePath)

	updateUserReq := &model.UserUpdateRequest{
		Profile: &model.ProfileRequest{
			Bio:       nil,
			WhyVoteMe: nil,
			ImageSrc:  &imageEndpoint,
		},
	}

	_, err = u.userService.UpdateUser(ctx, updateUserReq)

	if err != nil {
		return "", err
	}

	return imageEndpoint, nil
}

func (u *userUploadService) detectMimeType(contentFile multipart.File) error {
	bytes, err := io.ReadAll(contentFile)

	if err != nil {
		u.log.Errorf("error reading content of file: %s", err)
		return err
	}

	mimeType := http.DetectContentType(bytes)

	if !utils.IsImageMimeType(mimeType) {
		u.log.Infof("content type is not of mime type image")
		return fmt.Errorf("image is of wrong type")
	}

	_, err = contentFile.Seek(0, 0)

	if err != nil {
		u.log.Errorf("error resetting file to start position: %s", err)
		return err
	}

	return nil
}
