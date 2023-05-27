package router

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/cafo13/fur-meds/backend/pets-api/auth"
	"github.com/cafo13/fur-meds/backend/pets-api/cors"
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
	CORSMiddleware cors.CORSMiddleware
	PetsRepository repository.PetsRepository
}

func NewRouter(
	authMiddleware auth.AuthMiddleware,
	corsMiddleware cors.CORSMiddleware,
	petsRepository repository.PetsRepository,
) GinRouter {
	return Router{
		Router:         gin.Default(),
		AuthMiddleware: authMiddleware,
		CORSMiddleware: corsMiddleware,
		PetsRepository: petsRepository,
	}
}

func (r Router) GetPets(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetsRepository.GetPets(ctx, user.UID)
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
	ctx.Header("Access-Control-Allow-Methods", "POST")

	pet := &repository.Pet{}
	err := ctx.BindJSON(&pet)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetsRepository.AddPet(ctx, user.UID, pet)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusCreated, pets)
		return
	}
}

func (r Router) UpdatePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "PUT")

	pet := &repository.Pet{}
	err := ctx.BindJSON(&pet)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	_, err = r.PetsRepository.GetPet(ctx, user.UID, pet.UUID.String())
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet with UUID '%s'", pet.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.PetsRepository.UpdatePet(
		ctx,
		user.UID,
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
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) DeletePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "DELETE")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUUID := ctx.Params.ByName("uuid")
	pets, err := r.PetsRepository.DeletePet(ctx, user.UID, petUUID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) SharePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	sharePetInviteRequest := &repository.SharePetInviteRequest{}
	err := ctx.BindJSON(&sharePetInviteRequest)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting share pet invite request from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pet, err := r.PetsRepository.GetPet(ctx, user.UID, sharePetInviteRequest.PetUUID.String())
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet with UUID '%s'", pet.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	userUidToSharePetWith, err := r.AuthMiddleware.GetUserUidByMail(ctx, sharePetInviteRequest.UserMailToInvite)
	if err != nil {
		errorMsg := fmt.Sprintf("error on getting UID of user '%s' to invite to pet share for pet with UUID '%s'", sharePetInviteRequest.UserMailToInvite, pet.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.PetsRepository.UpdatePet(
		ctx,
		user.UID,
		pet.UUID.String(),
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			if firestorePet.SharedWithUsers == nil {
				firestorePet.SharedWithUsers = []repository.SharedUsers{}
			}
			for _, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == userUidToSharePetWith {
					return nil, fmt.Errorf("user '%s' is already invited to accept share for pet '%s'", sharePetInviteRequest.UserMailToInvite, pet.Name)
				}
			}

			firestorePet.SharedWithUsers = append(firestorePet.SharedWithUsers, repository.SharedUsers{UserUid: userUidToSharePetWith, ShareAccepted: false})

			return firestorePet, nil
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error inviting user to accept share of pet")
		log.Error(wrappedError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) AcceptPetShare(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	acceptPetShareRequest := &repository.AcceptPetShareRequest{}
	err := ctx.BindJSON(&acceptPetShareRequest)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting accept pet share request from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetsRepository.UpdatePet(
		ctx,
		user.UID,
		acceptPetShareRequest.PetUUID.String(),
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			noInviteFoundError := fmt.Errorf("no open invite exists for user '%s' at pet '%s'", user.UID, acceptPetShareRequest.PetUUID.String())
			if firestorePet.SharedWithUsers == nil {
				return nil, noInviteFoundError
			}

			for index, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == user.UID {
					firestorePet.SharedWithUsers[index].ShareAccepted = true
					return firestorePet, nil
				}
			}

			return nil, noInviteFoundError
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error accepting invite to share of pet")
		log.Error(wrappedError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) StartRouter(port string) {
	r.Router.Use(r.CORSMiddleware.Middleware())
	r.Router.Use(r.AuthMiddleware.Middleware())

	v1 := r.Router.Group("/api/v1")
	{
		v1.GET("/pets", r.GetPets)

		v1.POST("/pet", r.AddPet)

		v1.PUT("/pet", r.UpdatePet)

		v1.POST("/pet/share/invite", r.SharePet)

		v1.POST("/pet/share/accept", r.AcceptPetShare)

		v1.DELETE("/pet/:uuid", r.DeletePet)
	}

	r.Router.Run(":" + port)
}
