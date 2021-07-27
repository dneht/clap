package util

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"reflect"
)

func ResultEmpty(c *fiber.Ctx) error {
	return c.SendString("")
}

func ResultParamEmpty(c *fiber.Ctx, e error) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return paramEmpty(c)
	}
}

func ResultParamMapOne(c *fiber.Ctx, e error, k1 string, v1 interface{}) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return paramResult(c, true, map[string]interface{}{
			k1: v1,
		})
	}
}

func ResultParamMapTwo(c *fiber.Ctx, e error, k1 string, v1 interface{}, k2 string, v2 interface{}) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return paramResult(c, true, map[string]interface{}{
			k1: v1,
			k2: v2,
		})
	}
}

func ResultParam(c *fiber.Ctx, e error, p bool, b interface{}) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return paramResult(c, p, b)
	}
}

func ResultParamWithMessage(c *fiber.Ctx, e error, p bool, msg string, b interface{}) error {
	if nil != e {
		return ErrorWithMessage(c, e, msg)
	} else {
		return paramResult(c, p, b)
	}
}

func ResultParamWithCodeMessage(c *fiber.Ctx, e error, p bool, code int, msg string, b interface{}) error {
	if nil != e {
		return ErrorWithCodeMessage(c, e, code, msg)
	} else {
		return paramResult(c, p, b)
	}
}

func paramResult(c *fiber.Ctx, p bool, b interface{}) error {
	if p {
		return paramNormal(c, b)
	} else {
		return paramEmpty(c)
	}
}

func paramNormal(c *fiber.Ctx, b interface{}) error {
	return c.JSON(b)
}

func paramEmpty(c *fiber.Ctx) error {
	return c.JSON(EmptyMap)
}

func ResultPageOrList(c *fiber.Ctx, p *MainInput, t func(*MainInput) (int64, error), f func(*MainInput) (int, interface{}, error)) error {
	size := len(p.Ids)
	if size > 100 {
		return ErrorInput(c, "input size error")
	}
	if size > 0 {
		_, result, err := f(p)
		if nil != err {
			return ErrorInputOrDirect(c, err, "can not get list")
		}
		return listNormal(c, result)
	} else {
		total, err := t(p)
		if nil != err {
			return ErrorInputOrDirect(c, err, "can not get page")
		}
		if total == 0 {
			return ResultPageEmpty(c, nil)
		} else {
			return ResultPage(c, total, p, func(input *MainInput) (int, interface{}, error) {
				return f(input)
			})
		}
	}
}

func ResultListEmpty(c *fiber.Ctx, e error) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return listEmpty(c)
	}
}

func ResultList(c *fiber.Ctx, e error, l int, b interface{}) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		if l > 0 {
			return listNormal(c, b)
		} else {
			return listEmpty(c)
		}
	}
}

func listNormal(c *fiber.Ctx, b interface{}) error {
	return c.JSON(b)
}

func listEmpty(c *fiber.Ctx) error {
	return c.JSON(EmptyList)
}

func ResultPageEmpty(c *fiber.Ctx, e error) error {
	if nil != e {
		return ErrorInternal(c, e)
	} else {
		return paramEmpty(c)
	}
}

func ResultPage(c *fiber.Ctx, t int64, p *MainInput, f func(*MainInput) (int, interface{}, error)) error {
	if t > 0 {
		s, b, e := f(p)
		if nil != e {
			return ErrorInternal(c, e)
		} else {
			return pageBuild(c, t, p, s, b)
		}
	} else {
		return pageEmpty(c)
	}
}

// t=totalNum, s=pageSize
func pageBuild(c *fiber.Ctx, t int64, p *MainInput, s int, b interface{}) error {
	if t <= 0 {
		return pageEmpty(c)
	} else {
		var tp int64
		sl := int64(s)
		if t%sl == 0 {
			tp = t / sl
		} else {
			tp = t/sl + 1
		}
		hn := false
		if p.PageNum < tp {
			hn = true
		}
		return c.JSON(&PageResult{
			TotalNum:  t,
			TotalPage: tp,
			PageSize:  int(p.PageSize),
			HasNext:   hn,
			Results:   b,
		})
	}
}

func pageEmpty(c *fiber.Ctx) error {
	return c.JSON(EmptyPage)
}

func ErrorInternal(c *fiber.Ctx, e error) error {
	log.Printf("[error] internal error: %v\n", e)
	if reflect.TypeOf(e) == reflect.TypeOf(fiber.ErrBadRequest) {
		return e
	} else {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"message": "Error",
		})
	}
}

func ErrorInput(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(map[string]string{
		"message": message,
	})
}

func ErrorInputOrDirect(c *fiber.Ctx, e error, message string) error {
	if reflect.TypeOf(e) == reflect.TypeOf(fiber.ErrBadRequest) {
		return e
	} else {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{
			"message": message,
		})
	}
}

func ErrorInputShowMessage(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(map[string]string{
		"message": message,
	})
}

func ErrorInputErrorMessage(c *fiber.Ctx, e error, message string) error {
	log.Printf("[error] input error: %v, %v\n", message, e)
	return c.Status(http.StatusBadRequest).JSON(map[string]string{
		"message": message,
	})
}

func ErrorWithMessage(c *fiber.Ctx, e error, message string) error {
	log.Printf("[error] result error: %v, %v\n", message, e)
	return c.Status(http.StatusInternalServerError).JSON(map[string]string{
		"message": message,
	})
}

func ErrorWithCodeMessage(c *fiber.Ctx, e error, code int, message string) error {
	log.Printf("[error] result error with code: %d, %v, %v\n", code, message, e)
	return c.Status(code).JSON(map[string]string{
		"message": message,
	})
}

var EmptyMap = make(map[string]string)
var EmptyList = make([]string, 0)
var EmptyPage = &PageResult{
	TotalNum:  0,
	TotalPage: 0,
	PageSize:  10,
	HasNext:   false,
	Results:   EmptyList,
}

type PageResult struct {
	TotalNum  int64       `json:"total" form:"total" query:"total"`
	TotalPage int64       `json:"-" form:"totalPage" query:"totalPage"`
	PageSize  int         `json:"size" form:"size" query:"size"`
	HasNext   bool        `json:"next" form:"next" query:"next"`
	Results   interface{} `json:"results" form:"results" query:"results"`
}
