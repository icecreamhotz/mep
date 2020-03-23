package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	ut "github.com/go-playground/universal-translator"
	"github.com/gofrs/uuid"
	"github.com/icecreamhotz/mep-api/models"
	"github.com/icecreamhotz/mep-api/utils"
	"github.com/spf13/viper"
)

type TodoListHandler struct {
	Service   models.TodoListReporer
	Validator ut.Translator
}

func NewTodoListHandler(repository models.TodoListReporer, validator ut.Translator) TodoListHandler {
	return TodoListHandler{
		Service:   repository,
		Validator: validator,
	}
}

func (handler *TodoListHandler) TodoListGet(c *gin.Context) {
	data, total, err := handler.Service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseObject(gin.H{
		"message": "Todolist data",
		"total":   total,
		"data":    data,
	}))
}

func (handler *TodoListHandler) TodoListPost(c *gin.Context) {
	var todoList models.TodoLists
	if err := c.ShouldBind(&todoList); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(handler.Validator, err))
		return
	}

	file := todoList.ImageFile
	imageName := ""
	if file != nil {
		filename := filepath.Base(file.Filename)
		extension := filepath.Ext(filename)
		message, ok := utils.ValidateExtension(extension, utils.DefaultExtension)
		if !ok {
			c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
				"image": message,
			}}))
			return
		}
		message, ok = utils.ValidateFileSize(file.Size, 1)
		if !ok {
			c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
				"image": message,
			}}))
			return
		}
		imageName = utils.GetTimeNowFormatYYYYMMDDHHIIMM() + extension
		err := utils.ImageSaver(file, viper.GetString("base_dir.todo_list.path"), imageName, viper.GetStringMap("base_dir.todo_list.size"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
			return
		}
	}

	todoList.Image = imageName

	err := handler.Service.Create(todoList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}
	c.JSON(http.StatusCreated, utils.ResponseMessage("Created successful."))
}

func (handler *TodoListHandler) TodoListPut(c *gin.Context) {
	var todoList models.TodoLists
	if err := c.ShouldBind(&todoList); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(handler.Validator, err))
		return
	}

	todoListID := c.Param("id")
	todoListUUID, err := uuid.FromString(todoListID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}

	oldTodoList, err := handler.Service.GetById(todoListUUID)
	if err != nil {
		if err == pg.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ReponseNotFound("Todo list not found"))
			return
		}
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}

	file := todoList.ImageFile
	imageName := oldTodoList.Image
	if file != nil {
		filename := filepath.Base(file.Filename)
		extension := filepath.Ext(filename)
		message, ok := utils.ValidateExtension(extension, utils.DefaultExtension)
		if !ok {
			c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
				"image": message,
			}}))
			return
		}
		message, ok = utils.ValidateFileSize(file.Size, 1)
		if !ok {
			c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
				"image": message,
			}}))
			return
		}
		imageName = utils.GetTimeNowFormatYYYYMMDDHHIIMM() + extension
		err := utils.ImageSaver(file, viper.GetString("base_dir.todo_list.path"), imageName, viper.GetStringMap("base_dir.todo_list.size"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
			return
		}
		if oldTodoList.Image != "" {
			err := utils.RemoveImageAllResolution(viper.GetString("base_dir.todo_list.path"), oldTodoList.Image)
			if err != nil {
				c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
				return
			}
		}
	}
	todoList.ID = todoListUUID
	todoList.Image = imageName
	todoList.CreatedAt = oldTodoList.CreatedAt
	todoList.DeletedAt = oldTodoList.DeletedAt

	err = handler.Service.Update(&todoList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseMessage("Updated successful."))
}

func (handler *TodoListHandler) TodoListDelete(c *gin.Context) {
	todoListID := c.Param("id")

	todoListUUID, err := uuid.FromString(todoListID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}

	oldTodoList, err := handler.Service.GetById(todoListUUID)
	if err != nil {
		if err == pg.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ReponseNotFound("Todo list not found"))
			return
		}
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}

	if oldTodoList.Image != "" {
		err := utils.RemoveImageAllResolution(viper.GetString("base_dir.todo_list.path"), oldTodoList.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
			return
		}
	}

	err = handler.Service.DeleteById(todoListUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}
	c.JSON(http.StatusNoContent, utils.ResponseMessage("Deleted successful."))
}

func (handler *TodoListHandler) TodoListDonePatch(c *gin.Context) {
	todoListID := c.Param("id")

	todoListUUID, err := uuid.FromString(todoListID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}

	result, err := handler.Service.UpdateDoneStatus(todoListUUID)
	if err != nil {
		if err == pg.ErrNoRows {
			c.JSON(http.StatusNotFound, utils.ReponseNotFound("Todo list not found"))
			return
		}
		c.JSON(http.StatusBadRequest, utils.ResponseServerError("Please check your uuid again"))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseObject(gin.H{
		"message": "Todolist data",
		"data":    result,
	}))
}
