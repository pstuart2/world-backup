package api

import (
	"github.com/labstack/echo"
)

func (api *API) getFolders(ctx echo.Context) error {
	//book, err := api.payload.CreateBook(ctx)
	//if err != nil {
	//	return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	//}
	//
	//log := getLogger(ctx)
	//db := getAppDb(ctx)
	//
	//if err := db.Insert(book); err != nil {
	//	log.Error("failed to insert book into database", err)
	//	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to create book"})
	//}
	//
	//log.WithFields(logrus.Fields{
	//	"_id":        book.Id.Hex(),
	//	"clubTypeId": book.ClubTypeId,
	//	"name":       book.Name,
	//}).Info("Book created")

	return nil //ctx.JSON(http.StatusCreated)
}
