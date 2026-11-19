package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"compass-quiz/api"
)

type handler struct {
	usecase usecase
}

func (h handler) GetStudents(ctx *gin.Context, params api.GetStudentsParams) {
	if params.FacilitatorId < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)

		return
	}

	students, totalCount := h.usecase.listStudents(params.FacilitatorId)

	response := api.StudentsResponse{
		Students:   make([]api.Student, len(students)),
		TotalCount: totalCount,
	}

	for i, record := range students {
		response.Students[i] = api.Student{
			Id:      record.id,
			Name:    record.name,
			LoginId: record.account.loginID,
			Classroom: api.Classroom{
				Id:   record.class.id,
				Name: record.class.name,
			},
		}
	}

	ctx.JSON(http.StatusOK, response)
}
