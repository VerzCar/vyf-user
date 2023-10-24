package app

import (
	"fmt"
	"github.com/VerzCar/vyf-user/api/model"
	routerContext "github.com/VerzCar/vyf-user/app/router/ctx"
	"github.com/VerzCar/vyf-user/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (s *Server) UserMe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find user",
			Data:   nil,
		}

		user, err := s.userService.User(ctx.Request.Context(), nil)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		userResponse := &model.UserResponse{
			ID:         user.ID,
			IdentityID: user.IdentityID,
			Username:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Gender:     user.Gender,
			Locale:     user.Locale,
			Address:    user.Address,
			Contact:    user.Contact,
			Profile:    user.Profile,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   userResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) UserX() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find user",
			Data:   nil,
		}

		userUriReq := &model.UserUriRequest{}

		err := ctx.ShouldBindUri(userUriReq)

		// because the json can be optionally empty check against EOF
		if err != nil && err != io.EOF {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(userUriReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		user, err := s.userService.User(ctx.Request.Context(), &userUriReq.IdentityID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		userResponse := &model.UserResponse{
			ID:         user.ID,
			IdentityID: user.IdentityID,
			Username:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Gender:     user.Gender,
			Locale:     user.Locale,
			Address:    user.Address,
			Contact:    user.Contact,
			Profile:    user.Profile,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   userResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "user cannot be updated",
			Data:   nil,
		}

		userUpdateReq := &model.UserUpdateRequest{}

		err := ctx.ShouldBindJSON(userUpdateReq)

		if err == io.EOF {
			response := model.Response{
				Status: model.ResponseNop,
				Msg:    "no update",
				Data:   nil,
			}
			ctx.JSON(http.StatusOK, response)
			return
		}

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(userUpdateReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		user, err := s.userService.UpdateUser(ctx.Request.Context(), userUpdateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		userResponse := &model.UserResponse{
			ID:         user.ID,
			IdentityID: user.IdentityID,
			Username:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Gender:     user.Gender,
			Locale:     user.Locale,
			Address:    user.Address,
			Contact:    user.Contact,
			Profile:    user.Profile,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   userResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) Users() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot find users",
			Data:   nil,
		}

		users, err := s.userService.Users(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		var paginatedUsersResponse []*model.UserPaginatedResponse

		for _, user := range users {
			userPaginatedResponse := &model.UserPaginatedResponse{
				IdentityID: user.IdentityID,
				Username:   user.Username,
				Profile: &model.ProfilePaginatedResponse{
					ImageSrc: user.ProfileImageSrc,
				},
			}
			paginatedUsersResponse = append(paginatedUsersResponse, userPaginatedResponse)
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   paginatedUsersResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) UploadProfileImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Msg:    "cannot upload file",
			Data:   nil,
		}

		authClaims, err := routerContext.ContextToAuthClaims(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("error getting auth claims: %s", err)
			ctx.JSON(http.StatusUnauthorized, errResponse)
			return
		}

		multiPartFile, err := ctx.FormFile("profileImageFile")

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		contentFile, err := multiPartFile.Open()
		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		defer contentFile.Close()

		bytes, err := io.ReadAll(contentFile)
		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		mimeType := http.DetectContentType(bytes)

		if !utils.IsImageMimeType(mimeType) {
			s.log.Errorf("file type is wrong type: %s", mimeType)
			errResponse.Msg = "file type is not an image"
			ctx.JSON(http.StatusNotAcceptable, errResponse)
			return
		}

		filePath := fmt.Sprintf("profile/image/%s/%s", authClaims.Subject, multiPartFile.Filename)

		_, err = s.extStorageService.Upload(
			ctx.Request.Context(),
			filePath,
			contentFile,
		)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		imageEndpoint := fmt.Sprintf("%s/%s", s.extStorageService.ObjectEndpoint(), filePath)

		updateUserReq := &model.UserUpdateRequest{
			Profile: &model.ProfileRequest{
				Bio:       nil,
				WhyVoteMe: nil,
				ImageSrc:  &imageEndpoint,
			},
		}

		user, err := s.userService.UpdateUser(ctx.Request.Context(), updateUserReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Msg:    "",
			Data:   user.Profile.ImageSrc,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
