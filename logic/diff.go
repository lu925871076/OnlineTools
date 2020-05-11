package logic

import (
	"com.lu/OnlineTools/data"
	"com.lu/OnlineTools/logic/base/diff"
	"errors"
	"github.com/labstack/echo"
)

func Diff(ctx echo.Context, req *data.DiffReq) (*data.DiffResp, error) {
	if req.Type == data.TypeTextExactDiff {
		html := diff.DoTextExactDiff(req.Text1, req.Text2)
		return &data.DiffResp{Html: html}, nil
	} else if req.Type == data.TypeTextDiffByLine {
		html := diff.DoTextDiffByLine(req.Text1, req.Text2)
		return &data.DiffResp{Html: html}, nil
	} else {
		return nil, errors.New("不支持此类型")
	}
}
