package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"runtime"
	"strings"
)

var (
	ctxType = reflect.TypeOf(&gin.Context{})
)

type ProxyRouter struct {
	apiHandlers []interface{}
	Routes      []RouteConfig
}

func (b *ProxyRouter) add(method, path string, handler gin.HandlerFunc) {
	b.Routes = append(b.Routes, RouteConfig{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (b *ProxyRouter) RegisterHandlersWithTags(handlers ...interface{}) {
	b.apiHandlers = append(b.apiHandlers, handlers...)
}

// LoadRouter quét tất cả các thư mục và tìm kiếm các route từ các tệp Go
func (b *ProxyRouter) LoadRouter() {
	if len(b.apiHandlers) == 0 {
		return
	}

	for _, apiHandler := range b.apiHandlers {
		val := reflect.ValueOf(apiHandler)

		filePath, err := getFilePathOfStruct(apiHandler)
		if err != nil {
			log.Printf("failed to get file path of struct %T: %v", apiHandler, err)
			continue
		}

		methodApiMap := parseApiTags(filePath)
		for methodName, route := range methodApiMap {
			method := val.MethodByName(methodName)
			if !method.IsValid() {
				log.Printf("method %s not found in handler %T", methodName, apiHandler)
				continue
			}

			// Ex: (h *MyApiHandler) SayHello(c *gin.Context) {}
			// Kiểm tra kiểu signature method
			methodType := method.Type()
			// method là bound method, nên đầu vào phải có 1 tham số (ctx)
			if methodType.NumIn() != 1 {
				log.Printf("method %s in %T must have exactly one input parameter (*gin.Context)", methodName, apiHandler)
				continue
			}

			// Kiểm tra kiểu *gin.Context
			if methodType.In(0) != ctxType {
				log.Printf("method %s in %T input parameter must be *gin.Context", methodName, apiHandler)
				continue
			}

			// Đảm bảo function không có giá trị trả về
			if methodType.NumOut() != 0 {
				log.Printf("method %s in %T must have no return value", methodName, apiHandler)
				continue
			}

			h := gin.HandlerFunc(func(ctx *gin.Context) {
				method.Call([]reflect.Value{reflect.ValueOf(ctx)})
			})

			b.add(route.Method, route.Path, h)
		}
	}
}

type ParseRoute struct {
	Path   string
	Method string
}

func parseApiTags(filename string) map[string]ParseRoute {
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	result := make(map[string]ParseRoute)
	baseUrl := ""

	// Tìm @BaseUrl
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				if strings.HasPrefix(comment.Text, "// @BaseUrl") {
					parts := strings.Fields(comment.Text)
					if len(parts) >= 2 {
						baseUrl = parts[2]
					}
				}
			}
		}
	}

	// Tìm @Api
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}
		for _, comment := range fn.Doc.List {
			if strings.HasPrefix(comment.Text, "// @Api") {
				parts := strings.Fields(comment.Text)
				if len(parts) != 4 {
					log.Printf("Invalid @Api comment format: %s", comment.Text)
					continue
				}
				method := parts[2]
				path := parts[3]
				fullPath := path
				if baseUrl != "" {
					fullPath = strings.TrimRight(baseUrl, "/") + "/" + strings.TrimLeft(path, "/")
				}
				result[fn.Name.Name] = ParseRoute{
					Method: method,
					Path:   fullPath,
				}
			}
		}
	}

	return result
}

func getFilePathOfStruct(i interface{}) (string, error) {
	typ := reflect.TypeOf(i)
	if typ.NumMethod() == 0 {
		return "", fmt.Errorf("struct %T has no methods", i)
	}
	pc := typ.Method(0).Func.Pointer()
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "", fmt.Errorf("cannot find function for struct")
	}
	file, _ := fn.FileLine(pc)
	return file, nil
}
