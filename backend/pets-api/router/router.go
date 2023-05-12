package router

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/cafo13/fur-meds/backend/pets-api/auth"
	"github.com/cafo13/fur-meds/backend/pets-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GinRouter interface {
	GetPets(ctx *gin.Context)
	AddPet(ctx *gin.Context)
	UpdatePet(ctx *gin.Context)
	DeletePet(ctx *gin.Context)

	StartRouter(port string)
}

type Router struct {
	Router         *gin.Engine
	AuthMiddleware auth.AuthMiddleware
	PetsRepository repository.PetsRepository
}

func NewRouter(authMiddleware auth.AuthMiddleware, petsRepository repository.PetsRepository) GinRouter {
	return Router{Router: gin.Default(), AuthMiddleware: authMiddleware, PetsRepository: petsRepository}
}

func (r Router) GetPets(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET")

	pets, err := r.PetsRepository.GetPets(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) AddPet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST")

	pet := &repository.Pet{}
	err := ctx.BindJSON(&pet)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	err = r.PetsRepository.AddPet(ctx, pet)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusCreated, pet)
		return
	}
}

func (r Router) UpdatePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "PUT")

	pet := &repository.Pet{}
	err := ctx.BindJSON(&pet)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	_, err = r.PetsRepository.GetPet(ctx, pet.UUID.String())
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet with UUID '%s'", pet.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	err = r.PetsRepository.UpdatePet(
		ctx,
		pet.UUID.String(),
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			if pet.Name != "" && pet.Name != firestorePet.Name {
				firestorePet.Name = pet.Name
			}
			if pet.Species != "" && pet.Species != firestorePet.Species {
				firestorePet.Species = pet.Species
			}
			if pet.Image != "" && pet.Image != firestorePet.Image {
				firestorePet.Image = pet.Image
			}
			if len(pet.Medicines) != 0 && !reflect.DeepEqual(pet.Medicines, firestorePet.Medicines) {
				firestorePet.Medicines = pet.Medicines
			}
			if len(pet.Foods) != 0 && !reflect.DeepEqual(pet.Foods, firestorePet.Foods) {
				firestorePet.Foods = pet.Foods
			}

			return firestorePet, nil
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating pet")
		log.Error(wrappedError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Message": "updated pet successfully"})
		return
	}
}

func (r Router) DeletePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "DELETE")

	id := ctx.Params.ByName("uuid")
	err := r.PetsRepository.DeletePet(ctx, id)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Message": "deleted pet successfully"})
		return
	}
}

func (r Router) StartRouter(port string) {
	v1 := r.Router.Group("/api/v1")
	{
		v1.GET("/pets", r.AuthMiddleware.Middleware(), r.GetPets)

		v1.POST("/pet", r.AuthMiddleware.Middleware(), r.AddPet)

		v1.PUT("/pet", r.AuthMiddleware.Middleware(), r.UpdatePet)

		v1.DELETE("/pet/:uuid", r.AuthMiddleware.Middleware(), r.DeletePet)
	}

	r.Router.Run(":" + port)
}
