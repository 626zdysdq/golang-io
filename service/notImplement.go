package service

import (
	"net/http"
)

// 找不到文件
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented", 501)
}
