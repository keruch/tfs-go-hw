package services

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/keruch/tfs-go-hw/hw4/internal/domain"
	"github.com/keruch/tfs-go-hw/hw4/internal/handlers"
	"github.com/keruch/tfs-go-hw/hw4/internal/repository"
	"github.com/keruch/tfs-go-hw/hw4/pkg/log"
	"github.com/keruch/tfs-go-hw/hw4/pkg/token"
)

type key int

const ctxTokenKey key = iota

type ChatService struct {
	Repo       repository.Repo
	Controller handlers.Controller
	Logger     *log.Logger
}

func NewChatService(repo repository.Repo, controller handlers.Controller, logger *log.Logger) *ChatService {
	return &ChatService{
		Repo:       repo,
		Controller: controller,
		Logger:     logger,
	}
}

func (cs *ChatService) Start() {
	server := &http.Server{
		Addr:    domain.ServerAddr,
		Handler: cs.service(),
	}

	serverCtx, shutdown := context.WithCancel(context.Background())
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-shutdownSig

		// give 5 seconds to shutdown
		forceShutdownCtx, forceShutdown := context.WithTimeout(serverCtx, time.Second*5)
		go func() {
			<-forceShutdownCtx.Done()
			if forceShutdownCtx.Err() == context.DeadlineExceeded {
				cs.Logger.Fatal("graceful shutdown timed out, forcing exit")
			}
			forceShutdown()
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(forceShutdownCtx)
		if err != nil {
			cs.Logger.Fatal(err)
		}

		shutdown()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		cs.Logger.Fatal(err)
	}

	<-serverCtx.Done()
	cs.Logger.Infof("server shutted down!")
}

func (cs *ChatService) service() http.Handler {
	root := chi.NewRouter()
	root.Use(middleware.Logger)

	// /register and /login endpoints
	root.Post(domain.Register, cs.register)
	root.Post(domain.Login, cs.login)

	r := chi.NewRouter()
	r.Use(cs.auth)

	// /users endpoint
	r.Get(domain.Users, cs.getUsers)

	// /messages endpoint
	r.Get(domain.AllMessages, cs.getMessages)
	r.Get(domain.NumMessages, cs.getNumMessages)
	r.Post(domain.AllMessages, cs.postMessage)
	r.Get(domain.PrivateMsg, cs.getPrivateMessages)
	r.Post(domain.PrivateTo, cs.postPrivateMessage)

	root.Mount("/", r)

	return root
}

func (cs *ChatService) auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get("Authorization")
		username, err := token.ValidateUserToken(userToken)
		if err != nil {
			cs.Logger.Errorf("auth: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknCtx := context.WithValue(r.Context(), ctxTokenKey, username)
		handler.ServeHTTP(w, r.WithContext(tknCtx))
	}

	return http.HandlerFunc(fn)
}
