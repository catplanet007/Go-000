package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

func main() {
	client := struct{ Call func() (*response, error) }{api}

	ch := time.Tick(time.Second)

	for true {
		if _, err := client.Call(); err != nil {
			color.Green.Printf("client side get error: %+v\n", err)
		}
		color.Gray.Println("=============================================")
		<-ch
	}
}

type response struct{}

func api() (*response, error) {
	r, err := biz()
	if err != nil {
		if IsBizErr(err) {
			return nil, err
		}
		color.Red.Printf("server side log error: %+v\n", err)
		return nil, NewBizErr(777777, "something unexpected happend")
	}
	return r, nil
}

func biz() (*response, error) {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) < 50 {
		return nil, NewBizErr(123456, "do biz stuff failed")
	}
	return dao()
}

func dao() (*response, error) {
	return nil, errors.WithStack(sql.ErrNoRows)
}

// ============================================================================
// package biz/errors.go
type bizErr struct {
	code   int64
	detail interface{}
}

func NewBizErr(code int64, args ...interface{}) error {
	e := &bizErr{code, args}
	if len(args) == 1 {
		e.detail = args[0]
	}
	return e
}

func (b *bizErr) Error() string {
	return fmt.Sprintf("code: %d, detail: %v", b.code, b.detail)
}

func IsBizErr(err error) bool {
	_, ok := err.(*bizErr)
	return ok
}
