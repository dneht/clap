package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"strconv"
)

func CheckIdQuery(c *fiber.Ctx, s string) (uint64, error) {
	in := c.Query(s)
	return checkIdInput(c, in)
}

func CheckIdInput(c *fiber.Ctx, s string) (uint64, error) {
	in := c.Params(s)
	return checkIdInput(c, in)
}

func checkIdInput(c *fiber.Ctx, in string) (uint64, error) {
	if "" == in {
		return 0, ErrorInput(c, "id is empty")
	}

	id, err := strconv.ParseUint(in, 10, 64)
	if nil == err {
		if id > 0 {
			return id, nil
		} else {
			return 0, ErrorInput(c, "id input error")
		}
	}
	return 0, err
}

func CheckWsIdQuery(c *websocket.Conn, s string) (uint64, error) {
	in := c.Query(s)
	return checkWsIdInput(c, in)
}

func CheckWsIdInput(c *websocket.Conn, s string) (uint64, error) {
	in := c.Params(s)
	return checkWsIdInput(c, in)
}

func checkWsIdInput(c *websocket.Conn, in string) (uint64, error) {
	if "" == in {
		return 0, fiber.ErrBadRequest
	}

	id, err := strconv.ParseUint(in, 10, 64)
	if nil == err {
		if id > 0 {
			return id, nil
		} else {
			return 0, fiber.ErrBadRequest
		}
	}
	return 0, err
}

func CheckMainInput(c *fiber.Ctx) (*MainInput, error) {
	param := new(MainInput)
	err := c.BodyParser(param)
	if nil != err {
		return nil, ErrorInputShowMessage(c, "main input error")
	}
	return param, nil
}

func CheckReferInput(c *fiber.Ctx) (*MainInput, error) {
	param, err := CheckMainInput(c)
	if nil != err {
		return nil, err
	}
	if (nil != param.Refer && len(param.Refer) > 0) || (nil != param.Ids && len(param.Ids) > 0) {
		return param, nil
	}
	return nil, ErrorInput(c, "refer must be set")
}

func CheckIntQueryInput(c *fiber.Ctx, s string, d int64) (int64, error) {
	in := c.Query(s)
	if "" == in {
		return d, nil
	}

	n, err := strconv.ParseInt(in, 10, 64)
	if nil != err {
		return 0, err
	}
	return n, nil
}

func CheckWsIntQueryInput(c *websocket.Conn, s string, d int64) (int64, error) {
	in := c.Query(s)
	if "" == in {
		return d, nil
	}

	n, err := strconv.ParseInt(in, 10, 64)
	if nil != err {
		return 0, err
	}
	return n, nil
}
