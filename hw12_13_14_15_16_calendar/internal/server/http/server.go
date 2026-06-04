package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/app"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	host        string
	port        int
	logger      *logger.Logger
	application *app.App
}

func NewServer(logger *logger.Logger, app *app.App, host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		logger:      logger,
		application: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/headers", headers)
	mux.HandleFunc("/swagger.json", swagger)
	mux.HandleFunc("/swagger/", swaggerUI)

	RegisterHandlers(mux, NewCalendarHandlers(s.application))

	srv := &http.Server{
		Addr:         s.host + ":" + strconv.Itoa(s.port),
		Handler:      corsMiddleware(loggingMiddleware(s.logger, mux)),
		ReadTimeout:  5 * time.Second,  // ⚠️ Критически важно
		WriteTimeout: 10 * time.Second, // ⚠️
		IdleTimeout:  60 * time.Second, // ⚠️ для HTTP/1.1 keep-alive
		// BaseContext: func(net.Listener) context.Context { return ctx }, // опционально
	}

	// Канал для ошибки запуска
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		return nil

	case err := <-errCh:
		// ✅ Безопасная проверка "это ошибка закрытия сервера?" даже при wrapping'е
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed to start or crashed: %w", err)
		}
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Println(ctx)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	_, _ = fmt.Fprint(w, "Hello, World!")
}

func hello(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	// Функции, служащие обработчиками, принимают
	// `http.ResponseWriter` и `http.Request` в качестве
	// аргументов. Объект `http.ResponseWriter` используется для заполнения
	// HTTP-ответа. Здесь наш простой ответ - это просто
	// "hello\n".
	_, _ = fmt.Fprint(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	// Этот обработчик делает что-то более сложное,
	// прочитав все HTTP-заголовки запроса и выведя их в тело ответа.
	for name, headers := range req.Header {
		for _, h := range headers {
			_, _ = fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func swagger(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, "api/swagger.json")
}

func swaggerUI(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Calendar API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "BaseLayout"
            });
        };
    </script>
</body>
</html>`))
}
