package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name    string
	Path    string
	InType  string
	OutType string
	Comment string
}

// HTTPRule 表示 HTTP 规则
type HTTPRule struct {
	Method string
	Path   string
	Body   string
}

func parseService(g *protogen.GeneratedFile, service *protogen.Service) *Service {
	serv := &Service{
		Name:    service.GoName,
		Methods: make([]Method, 0),
	}

	for _, method := range service.Methods {
		serv.Methods = append(serv.Methods, Method{
			Name:    method.GoName,
			InType:  g.QualifiedGoIdent(method.Input.GoIdent),
			OutType: g.QualifiedGoIdent(method.Output.GoIdent),
			Path:    getMethodHTTPRule(method).Path,
			Comment: strings.TrimSpace(method.Comments.Leading.String()), // 获取方法前面的注释
		})
	}

	return serv
}

func render(serv *Service) []byte {
	parse, err := template.New("tpl").Parse(tmpl)
	if err != nil {
		return nil
	}

	buf := &bytes.Buffer{}
	_ = parse.Execute(buf, map[string]any{
		"ServiceName": serv.Name,
		"Methods": lo.Map(serv.Methods, func(item Method, index int) map[string]string {
			return map[string]string{
				"Name":     item.Name,
				"Path":     item.Path,
				"Request":  item.InType,
				"Response": item.OutType,
				"Comment":  item.Comment,
			}
		}),
	})

	return buf.Bytes()
}

// getMethodHTTPRule 从方法中提取 HTTP 规则
func getMethodHTTPRule(method *protogen.Method) *HTTPRule {
	// 解析真实的 google.api.http 注解
	options := method.Desc.Options().(*descriptorpb.MethodOptions)
	if options != nil {
		// 检查是否有 google.api.http 注解
		if proto.HasExtension(options, annotations.E_Http) {
			httpRule := proto.GetExtension(options, annotations.E_Http).(*annotations.HttpRule)
			if httpRule != nil {
				return parseHTTPRule(httpRule)
			}
		}
	}

	// 如果没有找到注解，使用默认规则
	defaultPath := generateRoutePath(method.Parent.GoName, method.GoName)
	return &HTTPRule{
		Method: "post",
		Path:   defaultPath,
		Body:   "*",
	}
}

// parseHTTPRule 解析 HTTP 规则
func parseHTTPRule(httpRule *annotations.HttpRule) *HTTPRule {
	rule := &HTTPRule{
		Body: httpRule.Body,
	}

	// 根据不同的 HTTP 方法设置
	switch pattern := httpRule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		rule.Method = "get"
		rule.Path = pattern.Get
	case *annotations.HttpRule_Post:
		rule.Method = "post"
		rule.Path = pattern.Post
	case *annotations.HttpRule_Put:
		rule.Method = "put"
		rule.Path = pattern.Put
	case *annotations.HttpRule_Delete:
		rule.Method = "delete"
		rule.Path = pattern.Delete
	case *annotations.HttpRule_Patch:
		rule.Method = "patch"
		rule.Path = pattern.Patch
	default:
		rule.Method = "post"
		rule.Path = "/"
	}

	return rule
}

func generateRoutePath(serviceName, methodName string) string {
	// 生成默认路由路径
	serviceNameLower := strings.ToLower(serviceName)
	methodNameLower := strings.ToLower(methodName)

	return fmt.Sprintf("/%s/%s", serviceNameLower, methodNameLower)
}
