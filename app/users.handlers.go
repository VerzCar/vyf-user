package app

import (
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/gin-gonic/gin"
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

		if err != nil {
			s.log.Error(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		if err := s.validate.Struct(userReq); err != nil {
			s.log.Warn(err)
			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		user, err := s.userService.User(ctx, &userReq.IdentityID)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Data:   user,
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

		user, err := s.userService.UpdateUser(ctx, userUpdateReq)

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.JSON(http.StatusInternalServerError, errResponse)
			return
		}

		response := model.Response{
			Status: model.ResponseSuccess,
			Data:   user,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
