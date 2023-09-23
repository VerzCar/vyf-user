package app

import (
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (s *Server) User() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Data:   "cannot find user",
		}

		userReq := &model.UserRequest{}

		err := ctx.ShouldBindJSON(userReq)

		// because the json can be optionally empty check against EOF
		if err != nil && err != io.EOF {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(userReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		user, err := s.userService.User(ctx.Request.Context(), userReq.IdentityID)

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
			Data:   userResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func (s *Server) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errResponse := model.Response{
			Status: model.ResponseError,
			Data:   "user cannot be updated",
		}

		userUpdateReq := &model.UserUpdateRequest{}

		err := ctx.ShouldBindJSON(userUpdateReq)

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
			Data:   userResponse,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
