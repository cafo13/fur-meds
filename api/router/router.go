package router

import (
	"fmt"
	"net/http"

	"github.com/cafo13/fur-meds/api/auth"
	"github.com/cafo13/fur-meds/api/cors"
	"github.com/cafo13/fur-meds/api/handler"
	"github.com/cafo13/fur-meds/api/repository"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var petAccessError = errors.New("user has no access to pet")

type HandlerSet struct {
	PetHandler      handler.PetHandler
	MedicineHandler handler.MedicineHandler
	FoodHandler     handler.FoodHandler
	TodoHandler     handler.TodoHandler
}
type Router struct {
	Router         *gin.Engine
	AuthMiddleware auth.AuthMiddleware
	CORSMiddleware cors.CORSMiddleware
	HandlerSet
}

func NewRouter(
	authMiddleware auth.AuthMiddleware,
	corsMiddleware cors.CORSMiddleware,
	handlerSet HandlerSet,
) Router {
	return Router{
		Router:         gin.Default(),
		AuthMiddleware: authMiddleware,
		CORSMiddleware: corsMiddleware,
		HandlerSet:     handlerSet,
	}
}

func (r Router) GetPets(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetHandler.GetAllForUser(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) GetPet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetHandler.Get(ctx, user.UID, petUuid)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	pet := &repository.Pet{}
	err = ctx.BindJSON(&pet)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	pets, err := r.PetHandler.Create(ctx, user.UID, pet)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	medicine := &repository.Medicine{}
	err = ctx.BindJSON(&medicine)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet medicine from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.MedicineHandler.Create(ctx, user.UID, petUuid, medicine)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	food := &repository.Food{}
	err = ctx.BindJSON(&food)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet food from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.FoodHandler.Create(ctx, user.UID, petUuid, food)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("uuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	_, err = r.PetHandler.Get(ctx, user.UID, petUuid)
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading pet with UUID '%s'", petUuid)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.PetHandler.Update(ctx, user.UID, petUuid, pet)
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

	medicine := &repository.Medicine{}
	err := ctx.BindJSON(&medicine)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting medicine from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	medicineUuid := ctx.Params.ByName("uuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting medicine UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	_, err = r.MedicineHandler.Get(ctx, user.UID, medicineUuid)
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading medicine with UUID '%s'", medicineUuid)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.MedicineHandler.Update(ctx, user.UID, medicineUuid, medicine)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating medicine")
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

	food := &repository.Food{}
	err := ctx.BindJSON(&food)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting food from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	foodUuid := ctx.Params.ByName("uuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting food UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	_, err = r.FoodHandler.Get(ctx, user.UID, foodUuid)
	if err != nil {
		errorMsg := fmt.Sprintf("error on loading food with UUID '%s'", foodUuid)
		log.Error(errorMsg)
		ctx.JSON(http.StatusNotFound, gin.H{"Message": errorMsg})
		return
	}

	pets, err := r.FoodHandler.Update(ctx, user.UID, foodUuid, food)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on updating food")
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("uuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetHandler.Delete(ctx, user.UID, petUuid)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	petMedicineUUID := ctx.Params.ByName("uuid")
	pets, err := r.MedicineHandler.Delete(ctx, user.UID, petMedicineUUID)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	petFoodUUID := ctx.Params.ByName("uuid")
	pets, err := r.FoodHandler.Delete(ctx, user.UID, petFoodUUID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, pets)
		return
	}
}

func (r Router) InviteToSharePet(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	sharePetInviteRequest := &repository.SharePetInviteRequest{}
	err := ctx.BindJSON(&sharePetInviteRequest)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting share pet invite request from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request body")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	pet, err := r.PetHandler.Get(ctx, user.UID, petUuid)
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

	pets, err := r.PetHandler.CreatePetShareInvite(ctx, user.UID, petUuid, userUidToSharePetWith)

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

func (r Router) AnswerPetShareInvite(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	petShareRequestAnswer := &repository.AnswerPetShareRequest{}
	err := ctx.BindJSON(&petShareRequestAnswer)
	if err != nil {
		wrappedError := errors.Wrap(err, "error on getting pet share request answer from json body")
		log.Error(wrappedError)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": wrappedError})
		return
	}

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request body")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	pets, err := r.PetHandler.AnswerPetShareInvite(ctx, user.UID, petUuid, petShareRequestAnswer.Answer)

	if err != nil {
		wrappedError := errors.Wrap(err, "error answering pet share invite")
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request body")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	openSharedPets, err := r.PetHandler.GetOpenSharedPets(ctx, user.UID)
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

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request body")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	petMedicines, err := r.MedicineHandler.GetAllForPet(ctx, user.UID, petUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, petMedicines)
		return
	}
}

func (r Router) GetPetMedicine(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	medicineUuid := ctx.Params.ByName("uuid")
	medicine, err := r.MedicineHandler.Get(ctx, user.UID, medicineUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, medicine)
		return
	}
}

func (r Router) GetPetFoods(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request body")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	petFoods, err := r.FoodHandler.GetAllForPet(ctx, user.UID, petUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, petFoods)
		return
	}
}

func (r Router) GetPetFood(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	petUuid := ctx.Params.ByName("petUuid")
	if len(petUuid) == 0 {
		err := errors.New("error on getting pet UUID from request URL")
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hasAccess, err := r.PetHandler.UserHasAccess(ctx, user.UID, petUuid)
	if err != nil {
		err := errors.New("error on checking if user has access to pet")
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !hasAccess {
		err := petAccessError
		log.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	foodUuid := ctx.Params.ByName("uuid")
	food, err := r.FoodHandler.Get(ctx, user.UID, foodUuid)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, food)
		return
	}
}

func (r Router) GetToDos(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	todos, err := r.TodoHandler.GetAllForUser(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, todos)
		return
	}
}

func (r Router) SetToDoStatus(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")

	user, err := r.AuthMiddleware.UserFromCtx(ctx)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	todos, err := r.TodoHandler.SetToDoStatus(ctx, user.UID)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	} else {
		ctx.IndentedJSON(http.StatusOK, todos)
		return
	}
}

func (r Router) StartRouter(port string) {
	r.Router.Use(r.CORSMiddleware.Middleware())
	r.Router.Use(r.AuthMiddleware.Middleware())

	v1 := r.Router.Group("/api/v1")
	{
		pets := v1.Group("/pets")
		{
			pets.POST("/", r.AddPet)

			pets.GET("/", r.GetPets)

			pets.GET("/:uuid", r.GetPet)

			pets.PUT("/:uuid", r.UpdatePet)

			pets.DELETE("/:uuid", r.DeletePet)

			medicines := pets.Group("/:petUuid/medicines")
			{
				medicines.POST("/", r.AddPetMedicine)

				medicines.GET("/", r.GetPetMedicines)

				medicines.GET("/:uuid", r.GetPetMedicine)

				medicines.PUT("/:uuid", r.UpdatePetMedicine)

				medicines.DELETE("/:uuid", r.DeletePetMedicine)
			}

			foods := pets.Group("/:petUuid/foods")
			{
				foods.POST("/", r.AddPetFood)

				foods.GET("/", r.GetPetFoods)

				foods.GET("/:uuid", r.GetPetFood)

				foods.PUT("/:uuid", r.UpdatePetFood)

				foods.DELETE("/:uuid", r.DeletePetFood)
			}

			shares := pets.Group("/:petUuid/shares")
			{
				invites := shares.Group("/invites")
				{
					invites.POST("/", r.InviteToSharePet)

					invites.GET("/", r.GetPetShareInvites)

					invites.POST("/answer", r.AnswerPetShareInvite)
				}
			}
		}

		todos := v1.Group("/todos")
		{
			todos.GET("/", r.GetToDos)

			todos.POST("/:uuid/status", r.SetToDoStatus)
		}
	}

	r.Router.Run(":" + port)
}
