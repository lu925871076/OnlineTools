package controller

import (
	"com.lu/OnlineTools/data"
	"com.lu/OnlineTools/logic"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func Diff(ctx echo.Context) error  {
	req := &data.DiffReq{}
	err := ctx.Bind(req)
	if err != nil {
		ctx.Logger().Warnf("Diff bind err: %s", err.Error())
		return err
	}
	//todo
	fmt.Println(req.Text1)
	fmt.Println(req.Text2)
	resp, err := logic.Diff(ctx, req)
	if err != nil {
		ctx.Logger().Warnf("Diff logic err: %s", err.Error())
		return err
	}
	return ctx.String(http.StatusOK, resp.Html)
}
