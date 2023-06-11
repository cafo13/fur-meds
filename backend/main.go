package main

import (
	"context"
	"errors"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"

	"github.com/cafo13/fur-meds/backend/auth"
	"github.com/cafo13/fur-meds/backend/cors"
	"github.com/cafo13/fur-meds/backend/handler"
	"github.com/cafo13/fur-meds/backend/repository"
	"github.com/cafo13/fur-meds/backend/router"

	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setupAuthMiddleware(gcpProject string) *auth.AuthMiddleware {
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
		authMiddleware := auth.NewMockAuthMiddleware()
		return &authMiddleware
	}

	config := &firebase.Config{ProjectID: gcpProject}
	firebaseApp, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		panic(err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	authMiddleware := auth.NewFirebaseAuthMiddleware(authClient)
	return &authMiddleware
}

func setupFirestoreClient(ctx context.Context, gcpProject string) *firestore.Client {
	client, err := firestore.NewClient(ctx, gcpProject)
	if err != nil {
		panic(err)
	}

	return client
}

func setupRouter(authHandler *auth.AuthMiddleware, corsHandler *cors.CORSMiddleware, handlerSet *router.HandlerSet) router.Router {
	router := router.NewRouter(*authHandler, *corsHandler, *handlerSet)
	return router
}

func main() {
	setupLogger()

	apiPort := os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		apiPort = "80"
	}
	gcpProject := os.Getenv("GCP_PROJECT")
	if len(gcpProject) == 0 {
		panic(errors.New("GCP_PROJECT environment variable needs to be set"))
	}

	authMiddleware := setupAuthMiddleware(gcpProject)
	corsMiddleware := cors.NewAllowingCORSMiddleware()
	firestoreClient := setupFirestoreClient(context.Background(), gcpProject)
	router := setupRouter(authMiddleware, &corsMiddleware, &router.HandlerSet{
		PetHandler:      handler.NewPetHandler(repository.NewPetFirestoreRepository(firestoreClient)),
		MedicineHandler: handler.NewMedicineHandler(repository.NewMedicineFirestoreRepository(firestoreClient)),
		FoodHandler:     handler.NewFoodHandler(repository.NewFoodFirestoreRepository(firestoreClient)),
		TodoHandler:     handler.NewTodoHandler(repository.NewTodoFirestoreRepository(firestoreClient)),
	})

	router.StartRouter(apiPort)
}
