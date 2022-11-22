package handlers

import (
	"api/internal/commons"
	"api/internal/events"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type GeneratorEnvironment struct {
	IcsTemplate *template.Template
}

func (ge GeneratorEnvironment) Generator(ctx *gin.Context) {
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	eventsPage, err := events.GetEvents(commons.User(ctx.Param("user")), commons.Offset(offset), commons.Limit(limit))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	if len(eventsPage) < 1 {
		ctx.AbortWithStatus(http.StatusNoContent)
	}

	for _, event := range eventsPage {
		err = ge.IcsTemplate.Execute(ctx.Writer, event)
		if err != nil {
			log.Print(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
