package common

import (
	"Hertz-Scaffold/biz/constant"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type Api struct {
	Model   string
	Type    string
	Path    string
	Handler app.HandlerFunc
}

var GlobalApiList []*Api
var GlobalApiMap = make(map[string]*Api)

func Register(module, httpType, path string, handler app.HandlerFunc) {
	api := &Api{
		Model:   module,
		Type:    httpType,
		Path:    "",
		Handler: handler,
	}

	switch module {
	case constant.DefaultAPIModule:
		api.Path = fmt.Sprintf("%s/%s", constant.DefaultURLPrefix, path)
	case constant.DevOpsAPIModule:
		api.Path = fmt.Sprintf("%s/%s", constant.DevOpsURLPrefix, path)
	case constant.AdminApIModel:
		api.Path = fmt.Sprintf("%s/%s", constant.AdminURLPrefix, path)
	default:
		panic(fmt.Sprintf("Path %v no match model", path))
	}
	if !RegisterMap(api.Path, api) {
		panic(fmt.Sprintf("Same path %s", api.Path))
	}

	GlobalApiList = append(GlobalApiList, api)
}

func RegisterMap(path string, param *Api) bool {
	if _, ok := GlobalApiMap[path]; !ok {
		GlobalApiMap[path] = param
		return true
	} else {
		return false
	}
}

func EngineRegister(engine *route.Engine, param *Api) {
	switch param.Type {
	case constant.MethodGet:
		engine.GET(param.Path, param.Handler)
	case constant.MethodPost:
		engine.POST(param.Path, param.Handler)
	case constant.MethodDelete:
		engine.DELETE(param.Path, param.Handler)
	case constant.MethodHead:
		engine.HEAD(param.Path, param.Handler)
	case constant.MethodOptions:
		engine.OPTIONS(param.Path, param.Handler)
	case constant.MethodPut:
		engine.PUT(param.Path, param.Handler)
	case constant.MethodPatch:
		engine.PATCH(param.Path, param.Handler)
	default:
		panic("no match http type")
	}
}
