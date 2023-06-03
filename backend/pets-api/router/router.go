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
	AddPetMedicine(ctx *gin.Context)
	AddPetFood(ctx *gin.Context)
	UpdatePet(ctx *gin.Context)
	UpdatePetMedicine(ctx *gin.Context)
	UpdatePetFood(ctx *gin.Context)
	SharePet(ctx *gin.Context)
	GetPetShareInvites(ctx *gin.Context)
	AcceptPetShare(ctx *gin.Context)
	DenyPetShare(ctx *gin.Context)
	DeletePet(ctx *gin.Context)
	DeletePetMedicine(ctx *gin.Context)
	DeletePetFood(ctx *gin.Context)
	GetPetMedicines(ctx *gin.Context)
	GetPetFoods(ctx *gin.Context)
	GetToDos(ctx *gin.Context)
	GenerateToDos(ctx *gin.Context)

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

func (r Router) AddPetMedicine(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	petMedicine := &repository.PetMedicine{}
	err := ctx.BindJSON(&petMedicine)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet medicine from json body")
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

	pets, err := r.PetsRepository.AddPetMedicine(ctx, user.UID, petMedicine)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusCreated, pets)
		return
	}
}

func (r Router) AddPetFood(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	petFood := &repository.PetFood{}
	err := ctx.BindJSON(&petFood)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet food from json body")
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

	pets, err := r.PetsRepository.AddPetFood(ctx, user.UID, petFood)
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

func (r Router) UpdatePetMedicine(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "PUT")

	petMedicine := &repository.PetMedicine{}
	err := ctx.BindJSON(&petMedicine)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet medicine from json body")
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

	_, err = r.PetsRepository.GetPetMedicine(ctx, user.UID, petMedicine.UUID.String())
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet medicine with UUID '%s'", petMedicine.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.PetsRepository.UpdatePetMedicine(
		ctx,
		user.UID,
		petMedicine.UUID.String(),
		func(context context.Context, firestorePetMedicine *repository.PetMedicine) (*repository.PetMedicine, error) {
			if petMedicine.Name != "" && petMedicine.Name != firestorePetMedicine.Name {
				firestorePetMedicine.Name = petMedicine.Name
			}
			if petMedicine.Dosage != 0 && petMedicine.Dosage != firestorePetMedicine.Dosage {
				firestorePetMedicine.Dosage = petMedicine.Dosage
			}
			if petMedicine.Unit != "" && petMedicine.Unit != firestorePetMedicine.Unit {
				firestorePetMedicine.Unit = petMedicine.Unit
			}
			if petMedicine.Stock != 0 && petMedicine.Stock != firestorePetMedicine.Stock {
				firestorePetMedicine.Stock = petMedicine.Stock
			}
			if len(petMedicine.Frequencies) != 0 && !reflect.DeepEqual(petMedicine.Frequencies, firestorePetMedicine.Frequencies) {
				firestorePetMedicine.Frequencies = petMedicine.Frequencies
			}

			return firestorePetMedicine, nil
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating pet medicine")
		log.Error(wrappedError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": wrappedError})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) UpdatePetFood(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "PUT")

	petFood := &repository.PetFood{}
	err := ctx.BindJSON(&petFood)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet food from json body")
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

	_, err = r.PetsRepository.GetPetFood(ctx, user.UID, petFood.UUID.String())
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet food with UUID '%s'", petFood.UUID)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.PetsRepository.UpdatePetFood(
		ctx,
		user.UID,
		petFood.UUID.String(),
		func(context context.Context, firestorePetFood *repository.PetFood) (*repository.PetFood, error) {
			if petFood.Name != "" && petFood.Name != firestorePetFood.Name {
				firestorePetFood.Name = petFood.Name
			}
			if petFood.Dosage != 0 && petFood.Dosage != firestorePetFood.Dosage {
				firestorePetFood.Dosage = petFood.Dosage
			}
			if petFood.Unit != "" && petFood.Unit != firestorePetFood.Unit {
				firestorePetFood.Unit = petFood.Unit
			}
			if petFood.Stock != 0 && petFood.Stock != firestorePetFood.Stock {
				firestorePetFood.Stock = petFood.Stock
			}
			if len(petFood.Frequencies) != 0 && !reflect.DeepEqual(petFood.Frequencies, firestorePetFood.Frequencies) {
				firestorePetFood.Frequencies = petFood.Frequencies
			}

			return firestorePetFood, nil
		},
	)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating pet food")
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

func (r Router) DeletePetMedicine(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "DELETE")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petMedicineUUID := ctx.Params.ByName("uuid")
	pets, err := r.PetsRepository.DeletePetMedicine(ctx, user.UID, petMedicineUUID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) DeletePetFood(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "DELETE")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petFoodUUID := ctx.Params.ByName("uuid")
	pets, err := r.PetsRepository.DeletePetFood(ctx, user.UID, petFoodUUID)
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
				firestorePet.SharedWithUsers = []repository.PetShares{}
			}
			for _, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == userUidToSharePetWith {
					return nil, fmt.Errorf("user '%s' is already invited to accept share for pet '%s'", sharePetInviteRequest.UserMailToInvite, pet.Name)
				}
			}

			firestorePet.SharedWithUsers = append(firestorePet.SharedWithUsers, repository.PetShares{UserUid: userUidToSharePetWith, ShareAccepted: false})

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

	acceptPetShareRequest := &repository.AnswerPetShareRequest{}
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

func (r Router) DenyPetShare(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	denyPetShareRequest := &repository.AnswerPetShareRequest{}
	err := ctx.BindJSON(&denyPetShareRequest)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting deny pet share request from json body")
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
		denyPetShareRequest.PetUUID.String(),
		func(context context.Context, firestorePet *repository.Pet) (*repository.Pet, error) {
			noInviteFoundError := fmt.Errorf("no open invite exists for user '%s' at pet '%s'", user.UID, denyPetShareRequest.PetUUID.String())
			if firestorePet.SharedWithUsers == nil {
				return nil, noInviteFoundError
			}

			for index, sharedUser := range firestorePet.SharedWithUsers {
				if sharedUser.UserUid == user.UID {
					firestorePet.SharedWithUsers = append(firestorePet.SharedWithUsers[:index], firestorePet.SharedWithUsers[index+1:]...)
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

func (r Router) GetPetShareInvites(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	openSharedPets, err := r.PetsRepository.GetOpenSharedPets(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		var petShareInvites []*repository.PetShareInvites
		for _, pet := range openSharedPets {
			ownerUser, err := r.AuthMiddleware.GetUserByUid(ctx, pet.UserUID)
			if err != nil {
				errorMsg := fmt.Sprintf("error on getting UID of pet owner user '%s'", pet.UserUID)
				log.Error(errorMsg)
				ctx.JSON(http.StatusInternalServerError, gin.H{"Message": errorMsg})
				return
			}
			petShareInvites = append(petShareInvites, &repository.PetShareInvites{
				Pet:        *pet,
				OwnerEmail: ownerUser.Email,
			})
		}
		ctx.IndentedJSON(http.StatusOK, petShareInvites)
		return
	}
}

func (r Router) GetPetMedicines(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	petMedicines, err := r.PetsRepository.GetPetMedicines(ctx, user.UID, petUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, petMedicines)
		return
	}
}

func (r Router) GetPetFoods(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	petFoods, err := r.PetsRepository.GetPetFoods(ctx, user.UID, petUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, petFoods)
		return
	}
}

func (r Router) GetToDos(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	todos, err := r.PetsRepository.GetToDos(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, todos)
		return
	}
}

func (r Router) GenerateToDos(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	err = r.PetsRepository.GenerateToDos(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"Message": "successfully generated todos"})
		return
	}
}

func (r Router) StartRouter(port string) {
	r.Router.Use(r.CORSMiddleware.Middleware())
	r.Router.Use(r.AuthMiddleware.Middleware())

	v1 := r.Router.Group("/api/v1")
	{
		// pets
		v1.GET("/pets", r.GetPets)

		v1.POST("/pet", r.AddPet)

		v1.PUT("/pet", r.UpdatePet)

		v1.DELETE("/pet/:uuid", r.DeletePet)

		// pet medicines
		v1.GET("/pet/:petUuid/medicines", r.GetPetMedicines)

		v1.POST("/pet/medicine", r.AddPetMedicine)

		v1.PUT("/pet/medicine", r.UpdatePetMedicine)

		v1.DELETE("/pet/medicine/:uuid", r.DeletePetMedicine)

		// pet foods
		v1.GET("/pet/:petUuid/foods", r.GetPetFoods)

		v1.POST("/pet/food", r.AddPetFood)

		v1.PUT("/pet/food", r.UpdatePetFood)

		v1.DELETE("/pet/food/:uuid", r.DeletePetFood)

		// pet shares
		v1.POST("/pet/share/invite", r.SharePet)

		v1.GET("/pet/share/invites", r.GetPetShareInvites)

		v1.POST("/pet/share/accept", r.AcceptPetShare)

		v1.POST("/pet/share/deny", r.DenyPetShare)

		// todos
		v1.GET("/todos", r.GetToDos)

		v1.POST("/todos/generate", r.GenerateToDos)
	}

	r.Router.Run(":" + port)
}
