package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"compass-quiz/api"
)

type handler struct {
	useCase studentUseCase
}

func newHandler(useCase studentUseCase) handler {
	return handler{useCase: useCase}
}

func (h handler) GetStudents(ctx *gin.Context, params api.GetStudentsParams) {
	students, totalCount, err := h.useCase.listStudents(ctx.Request.Context(), params.FacilitatorId)
	if err != nil {
		_ = ctx.Error(fmt.Errorf("get students: %w", err))
		ctx.Abort()

		return
	}

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
