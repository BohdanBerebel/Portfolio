package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/models"
	"github.com/BohdanBerebel/Portfolio/go_app/pkg/response"
	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1,max=10000"`
}

type NoteHandler struct {
	noteService NoteService
}

type NoteService interface {
	Create(
		ctx context.Context,
		userID int64,
		title string,
		content string,
	) (*models.Note, error)

	GetByUserID(
		ctx context.Context,
		userID int64,
	) ([]models.Note, error)

	GetByID(
		ctx context.Context,
		noteID int64,
		userID int64,
	) (*models.Note, error)

	Update(
		ctx context.Context,
		noteID int64,
		userID int64,
		title string,
		content string,
	) (*models.Note, error)

	Delete(
		ctx context.Context,
		noteID int64,
		userID int64,
	) error
}

type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1,max=10000"`
}

func NewNoteHandler(
	noteService NoteService,
) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

func (h *NoteHandler) Create(c *gin.Context) {
	var input CreateNoteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	note, err := h.noteService.Create(
		c.Request.Context(),
		userID,
		input.Title,
		input.Content,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to create note",
		)
		return
	}

	response.JSON(c, http.StatusCreated, note)
}

func (h *NoteHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	notes, err := h.noteService.GetByUserID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to get notes",
		)
		return
	}

	response.JSON(
		c,
		http.StatusOK,
		notes,
	)
}

func (h *NoteHandler) GetByID(c *gin.Context) {
	noteID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
		)
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	note, err := h.noteService.GetByID(
		c.Request.Context(),
		noteID,
		userID,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.JSON(
		c,
		http.StatusOK,
		note,
	)
}

func (h *NoteHandler) Update(c *gin.Context) {
	noteID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
		)
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var input UpdateNoteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	note, err := h.noteService.Update(
		c.Request.Context(),
		noteID,
		userID,
		input.Title,
		input.Content,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.JSON(
		c,
		http.StatusOK,
		note,
	)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	noteID, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid note id",
		)
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		return
	}

	err = h.noteService.Delete(
		c.Request.Context(),
		noteID,
		userID,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
