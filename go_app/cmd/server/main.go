package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/config"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/database"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/handlers"
	middlewares "github.com/BohdanBerebel/Portfolio/go_app/internal/middleware"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/repositories"
	"github.com/BohdanBerebel/Portfolio/go_app/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.Connect(
		context.Background(),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
	)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	userRepository := repositories.NewUserRepository(db)
	noteRepository := repositories.NewNoteRepository(db)

	userService := services.NewUserService(
		userRepository,
		cfg.JWTSecret,
	)

	noteService := services.NewNoteService(
		noteRepository,
	)

	authHandler := handlers.NewAuthHandler(
		userService,
	)

	noteHandler := handlers.NewNoteHandler(
		noteService,
	)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	router.POST(
		"/api/auth/register",
		authHandler.Register,
	)

	router.POST(
		"/api/auth/login",
		authHandler.Login,
	)

	protected := router.Group("/api")
	protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	protected.GET("/notes", noteHandler.GetAll)
	protected.GET("/notes/:id", noteHandler.GetByID)
	protected.POST("/notes", noteHandler.Create)
	protected.PUT("/notes/:id", noteHandler.Update)
	protected.DELETE("/notes/:id", noteHandler.Delete)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		slog.Info("server started", "addr", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}

	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("closing database connection pool")

	db.Close()

	slog.Info("server stopped")
}
