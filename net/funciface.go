package net

//WSConnIface hook函数接口
type WSConnIface interface {
	SetProperty(string, interface{})
	GetPropertyByKey(string) (interface{}, error)
	RemoveProperty(string)
	GetAddr() string
}
