package graphql

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	graphql_context "github.com/jerbob92/hoppscotch-backend/api/controllers/graphql/context"

	"github.com/graph-gophers/graphql-go"
)

// Request is the request content
type Request struct {
	OperationName string
	Query         string
	Variables     map[string]interface{}
	Context       context.Context
}

func set(v interface{}, m interface{}, path string) error {
	var parts []interface{}
	for _, p := range strings.Split(path, ".") {
		if isNumber, err := regexp.MatchString(`\d+`, p); err != nil {
			return err
		} else if isNumber {
			index, _ := strconv.Atoi(p)
			parts = append(parts, index)
		} else {
			parts = append(parts, p)
		}
	}
	for i, p := range parts {
		last := i == len(parts)-1
		switch idx := p.(type) {
		case string:
			if last {
				m.(map[string]interface{})[idx] = v
			} else {
				m = m.(map[string]interface{})[idx]
			}
		case int:
			if last {
				m.([]interface{})[idx] = v
			} else {
				m = m.([]interface{})[idx]
			}
		}
	}
	return nil
}

type File struct {
	File     multipart.File
	Filename string
	Size     int64
}

func Handle(c *graphql_context.Context, exec func(req *Request) *graphql.Response) {
	c.GinContext.Writer.Header().Set("Content-Type", "application/json")

	var operations interface{}

	logRequestErrors := func(r *Request, result *graphql.Response) {
		if result.Errors != nil && len(result.Errors) > 0 {
			for i, _ := range result.Errors {
				meta := map[string]interface{}{
					"graphql_query":     r.Query,
					"http_method":       c.GinContext.Request.Method,
					"http_content_type": c.GinContext.Request.Header.Get("Content-Type"),
				}
				if r.Variables != nil && len(r.Variables) > 0 {
					meta["graphql_variables"] = r.Variables
				}
				c.LogErr(result.Errors[i], meta)
			}
		}
	}

	makeRequest := func(r *Request) {
		result := exec(r)
		logRequestErrors(r, result)
		err := json.NewEncoder(c.GinContext.Writer).Encode(result)
		if err != nil {
			http.Error(c.GinContext.Writer, "Could not encode the graphql response", http.StatusInternalServerError)
			return
		}
	}

	switch c.GinContext.Request.Method {
	case "GET":
		request := Request{
			Context:   c.GinContext.Request.Context(),
			Variables: map[string]interface{}{},
		}

		for key, args := range c.GinContext.Request.URL.Query() {
			if args == nil || len(args) == 0 {
				continue
			}
			arg := args[0]
			if arg == "" {
				continue
			}
			switch strings.ToLower(key) {
			case "query":
				request.Query = arg
			case "variables":
				err := json.Unmarshal([]byte(arg), &request.Variables)
				if err != nil {
					http.Error(c.GinContext.Writer, "Graphql url variables are not valid json", http.StatusBadRequest)
					return
				}
			case "operationname":
				request.OperationName = arg
			}
		}

		if request.Query == "" {
			http.Error(c.GinContext.Writer, "Missing graphql query in url", http.StatusBadRequest)
			return
		}

		makeRequest(&request)
	default: // "POST", "PATCH", "DELETE", "PUT"
		contentType := strings.SplitN(c.GinContext.Request.Header.Get("Content-Type"), ";", 2)[0]

		switch contentType {
		case "text/plain", "application/json", "application/graphql":
			err := json.NewDecoder(c.GinContext.Request.Body).Decode(&operations)
			if err != nil {
				http.Error(c.GinContext.Writer, "Could not read the json request body", http.StatusBadRequest)
				return
			}
		case "multipart/form-data":
			// Parse multipart form
			err := c.GinContext.Request.ParseMultipartForm(8192)
			if err != nil {
				http.Error(c.GinContext.Writer, "Could not access uploaded file", http.StatusBadRequest)
				return
			}

			// Unmarshal uploads
			var uploads = map[File][]string{}
			var uploadsMap = map[string][]string{}
			if err := json.Unmarshal([]byte(c.GinContext.Request.Form.Get("map")), &uploadsMap); err != nil {
				panic(err)
			} else {
				for key, path := range uploadsMap {
					file, header, err := c.GinContext.Request.FormFile(key)
					if err != nil {
						http.Error(c.GinContext.Writer, "Could not access uploaded file", http.StatusInternalServerError)
						return
					}
					uploads[File{
						File:     file,
						Size:     header.Size,
						Filename: header.Filename,
					}] = path
				}
			}

			// Unmarshal operations
			if err := json.Unmarshal([]byte(c.GinContext.Request.Form.Get("operations")), &operations); err != nil {
				http.Error(c.GinContext.Writer, "the request form operations field doesn't exist or has invalid json data", http.StatusInternalServerError)
				return
			}

			// set uploads to operations
			for file, paths := range uploads {
				for _, path := range paths {
					if err := set(file, operations, path); err != nil {
						http.Error(c.GinContext.Writer, "Could not access uploaded file", http.StatusInternalServerError)
						return
					}
				}
			}
		}

		switch data := operations.(type) {
		case map[string]interface{}:
			request := Request{}

			for key, raw := range data {
				switch strings.ToLower(key) {
				case "operationname":
					if val, ok := raw.(string); ok {
						request.OperationName = val
					}
				case "query":
					if val, ok := raw.(string); ok {
						request.Query = val
					}
				case "variables":
					if val, ok := raw.(map[string]interface{}); ok {
						request.Variables = val
					}
				}
			}

			request.Context = c.GinContext.Request.Context()
			makeRequest(&request)
		case []interface{}:
			result := make([]interface{}, len(data))
			for index, operation := range data {
				data, ok := operation.(map[string]interface{})
				if !ok {
					http.Error(c.GinContext.Writer, "Invalid "+c.GinContext.Request.Method+" data", http.StatusInternalServerError)
					return
				}

				request := Request{}

				for key, raw := range data {
					switch strings.ToLower(key) {
					case "operationname":
						if val, ok := raw.(string); ok {
							request.OperationName = val
						}
					case "query":
						if val, ok := raw.(string); ok {
							request.Query = val
						}
					case "variables":
						if val, ok := raw.(map[string]interface{}); ok {
							request.Variables = val
						}
					}
				}

				request.Context = c.GinContext.Request.Context()
				reqResults := exec(&request)
				logRequestErrors(&request, reqResults)
				result[index] = reqResults
			}

			err := json.NewEncoder(c.GinContext.Writer).Encode(result)
			if err != nil {
				http.Error(c.GinContext.Writer, "Internal error", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(c.GinContext.Writer, "Could not encode the graphql response", http.StatusBadRequest)
			return
		}
	}
}
