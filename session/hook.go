package session

import (
	"github.com/felixorbit/fexorm/log"
	"reflect"
)

// Hook 是留在 被调用对象身上的钩子，调用者可以尝试“勾住”想要的钩子，同时需要将自己作为参数传入

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// TODO: 可以改成接口实现
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	params := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		v := fm.Call(params)
		if len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
