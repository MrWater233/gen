package gen

import (
	"log"
	"net/http"
	"strings"
)

// Define router to save route information
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// Create router
// roots key eg, roots["GET"] roots["post"]
// handlers key eg, handlers["GET-/p/:lang/doc"] handlers["POST-/p/book"]
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Parse pattern to array
func parsePattern(pattern string) []string {
	// Split the pattern eg /p/test => "","p","test"
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		// Exclude the first item
		if item != "" {
			parts = append(parts, item)
			// Omit path after *
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// Add route information
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// Get route information
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	// Split the pattern to be matched
	searchParts := parsePattern(pattern)
	// When ":" or "*" in pattern, save parse result in map
	// eg: /p/go/doc(/p/:lang/doc)=>{lang: "go"}
	//     /static/css/geektutu.css(/static/*filepath)=>{filepath: "css/geektutu.css"}
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// Find pattern result
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			} else if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// Handle request
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// append router handler to the handler-list
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// execuate the handler-list
	c.Next()
}
