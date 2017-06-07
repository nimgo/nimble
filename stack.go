package nimble

import (
	"net/http"
)

// Stack is a stack of Middleware Handlers that can be invoked as an http.Handler.
// The middleware stack is run in the sequence that they are added to the stack.
type Stack struct {
	handlers   []HandlerFunc
	middleware middleware
}

// HandlerFunc is a linked-list handler interface that provides
// every middleware a forward reference to the next middleware in the stack.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Handler exposes an adapter to support specific middleware that uses this signature
// ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// Each Middleware should yield to the next middleware in the chain by invoking the next http.HandlerFunc
type middleware struct {
	fn   HandlerFunc
	next *middleware
}

// Make sure Stack conforms with the http.Handler interface
var _ http.Handler = New()

// Stack itself is a http.Handler. This allows it to used as a substack manager
func (n *Stack) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(Writer); ok { // handle substacks
		n.middleware.serve(w, r)
	} else {
		n.middleware.serve(newWriter(w), r)
	}
}

// With adds a http.Handler onto the middleware stack.
func (n *Stack) With(handler http.Handler) *Stack {
	return n.WithHandlerFunc(wrap(handler))
}

// WithFunc adds a http.HandlerFunc onto the middleware stack.
func (n *Stack) WithFunc(handlerFunc http.HandlerFunc) *Stack {
	return n.WithHandlerFunc(wrapHandlerFunc(handlerFunc))
}

// WithHandler adds a Stack.Handler onto the middleware stack.
func (n *Stack) WithHandler(handler Handler) *Stack {
	return n.WithHandlerFunc(handler.ServeHTTP)
}

// WithHandlerFunc adds a Stack.HandlerFunc function onto the middleware stack.
func (n *Stack) WithHandlerFunc(handlerFunc HandlerFunc) *Stack {
	if handlerFunc == nil {
		panic("handlerFunc cannot be nil")
	}

	n.handlers = append(n.handlers, handlerFunc)
	n.middleware = build(n.handlers)
	return n
}

// The next http.HandlerFunc is automatically called after the Handler is executed.
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
func (m *middleware) serve(w http.ResponseWriter, r *http.Request) {
	m.fn(w, r, m.next.serve)
}

// Wrap converts a http.Handler into a Stack.HandlerFunc
func wrap(handler http.Handler) HandlerFunc {
	if handler == nil {
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(w, r)
		next(w, r)
	}
}

// wrapFunc converts a http.HandlerFunc into a Stack.HandlerFunc.
func wrapHandlerFunc(fn http.HandlerFunc) HandlerFunc {
	if fn == nil {
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fn(w, r)
		next(w, r)
	}
}

func build(handles []HandlerFunc) middleware {
	var next middleware

	if len(handles) == 0 {
		return empty()
	} else if len(handles) > 1 {
		next = build(handles[1:])
	} else {
		next = empty()
	}

	return middleware{handles[0], &next}
}

func empty() middleware {
	return middleware{
		func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { /* do nothing */ },
		&middleware{},
	}
}
