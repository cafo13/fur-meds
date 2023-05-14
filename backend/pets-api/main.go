package main

import (
	"context"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"

	"github.com/cafo13/fur-meds/backend/pets-api/auth"
	"github.com/cafo13/fur-meds/backend/pets-api/cors"
	"github.com/cafo13/fur-meds/backend/pets-api/repository"
	"github.com/cafo13/fur-meds/backend/pets-api/router"

	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
)

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setupAuthMiddleware() *auth.AuthMiddleware {
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); mockAuth {
		authMiddleware := auth.NewMockAuthMiddleware()
		return &authMiddleware
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
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

func setupPetsRepository(ctx context.Context, gcpProject string) *repository.PetsRepository {
	client, err := firestore.NewClient(ctx, gcpProject)
	if err != nil {
		panic(err)
	}

	petsFirestoreRepository := repository.NewPetsFirestoreRepository(client)
	return &petsFirestoreRepository
}

func setupRouter(authHandler *auth.AuthMiddleware, corsHandler *cors.CORSMiddleware, petsRepository *repository.PetsRepository) router.GinRouter {
	router := router.NewRouter(*authHandler, *corsHandler, *petsRepository)
	return router
}

func main() {
	setupLogger()

	apiPort := os.Getenv("API_PORT")

	authMiddleware := setupAuthMiddleware()
	corsMiddleware := cors.NewAllowingCORSMiddleware()
	petsRepository := setupPetsRepository(context.Background(), os.Getenv("GCP_PROJECT"))
	router := setupRouter(authMiddleware, &corsMiddleware, petsRepository)

	router.StartRouter(apiPort)
}
