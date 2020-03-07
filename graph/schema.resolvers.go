// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import "C"
import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/jonathanho/gqlgen/graph/generated"
	"github.com/jonathanho/gqlgen/graph/model"
	"google.golang.org/api/iterator"
	"googlemaps.github.io/maps"
)

func (r *queryResolver) Uniqueness(ctx context.Context, newEmail string, newUserName string) (*model.UniquenessResult, error) {
	dbclient, _:= getDbClient()

	//check for email uniquenes
	docs := dbclient.Collection(USER_COLLECTION_NAME).Where(EMAIL_FIELD_NAME, "==", newEmail).Documents(ctx)
	var countEmail = 0
	var err error
	for {

		_, err = docs.Next()
		if err == iterator.Done {
			break
		}
		countEmail++
		if countEmail != 0 {
			break
		}
	}
	if countEmail != 0 {
		result := &model.UniquenessResult{
			EmailUnique:    false,
			UsernameUnique: true,
			Message:        "This Email already has an account",
		}
		defer dbclient.Close()
		return result, err
	}
	// check for username uniqueness
	docs = dbclient.Collection(USER_COLLECTION_NAME).Where(USERNAME_FIELD_NAME, "==", newUserName).Documents(ctx)
	var countUser = 0

	for {

		_, err = docs.Next()
		if err == iterator.Done {
			break
		}
		countUser++
		if countUser != 0 {
			break
		}
	}
	if countUser != 0 {
		result := &model.UniquenessResult{
			EmailUnique:    true,
			UsernameUnique: false,
			Message:        "This username is taken",
		}
		defer dbclient.Close()
		return result, err
	}
	//is valid
	result := &model.UniquenessResult{
		EmailUnique:    true,
		UsernameUnique: true,
		Message:        "This is a valid account",
	}
	defer dbclient.Close()
	return result , nil

}
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	//get the next userID from the server

	userName := input.Username
	password := GetMD5Hash(input.Password)
	email := input.Email
	//network call to put the user in the database
	client, _ := getDbClient()
	dbUser, _, err := client.Collection(USER_COLLECTION_NAME).Add(ctx, map[string]interface{}{
		USERNAME_FIELD_NAME:     userName,
		EMAIL_FIELD_NAME: email,
		PASSWORD_FIELD_NAME:     password,
	})

	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	user := &model.User{
		UserID:   dbUser.ID,
		Password: password,
		Username: userName,
		Email:    email,
	}
	defer client.Close()
	return user, nil
}

func (r *mutationResolver) CreateHome(ctx context.Context, input model.NewHome) (*model.Home, error) {
	client, _ := getMapsClient()

	//get the lat and long
	results, err := client.Geocode(ctx, &maps.GeocodingRequest{Address: input.Address})
	if err != nil {
		fmt.Println("api request error")
	}

	result := results[0]
	latitude := result.Geometry.Location.Lat
	longitude := result.Geometry.Location.Lng

	dbclient, _ := getDbClient()
	user := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID)
	doc, _, dbErr := user.Collection(HOME_COLLECTION_NAME).Add(ctx, map[string]interface{}{
		NAME_FIELD_NAME:      input.Name,
		ADDRESS_FIELD_NAME:   input.Address,
		RENT_FIELD_NAME:      input.Rent,
		BEDS_FIELD_NAME:   input.NumBed,
		CITY_FIELD_NAME:      input.City,
		LATITUDE_FIELD_NAME:  latitude,
		LONGITUDE_FIELD_NAME: longitude,
	})
	if dbErr != nil {
		panic(dbErr)
		return nil, dbErr
	}

	home := &model.Home{
		HomeID:    doc.ID,
		Name:      input.Name,
		Address:   input.Address,
		Rent:      input.Rent,
		Latitude:  latitude,
		Longitude: longitude,
		NumBed:    input.NumBed,
		City:      input.City,
	}
	defer dbclient.Close()
	return home, nil
}

func (r *mutationResolver) CreateWork(ctx context.Context, input model.NewWork) (*model.Work, error) {
	client, ctx := getMapsClient()

	//get the lat and long
	results, err := client.Geocode(ctx, &maps.GeocodingRequest{Address: input.Address})
	if err != nil {
		fmt.Println("api request error")
	}

	result := results[0]
	latitude := result.Geometry.Location.Lat
	longitude := result.Geometry.Location.Lng

	dbclient, dbctx := getDbClient()
	user := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID)
	doc, _, dbErr := user.Collection(WORK_COLLECTION_NAME).Add(dbctx, map[string]interface{}{
		NAME_FIELD_NAME:      input.Name,
		ADDRESS_FIELD_NAME:   input.Address,
		CITY_FIELD_NAME:      input.City,
		LATITUDE_FIELD_NAME:  latitude,
		LONGITUDE_FIELD_NAME: longitude,
	})
	if dbErr != nil {
		panic(dbErr)
		return nil, dbErr
	}

	home := &model.Work{
		WorkID:    doc.ID,
		Name:      input.Name,
		Address:   input.Address,
		Latitude:  latitude,
		Longitude: longitude,
		City:      input.City,
	}
	defer dbclient.Close()
	return home, nil
}

func (r *queryResolver) Homes(ctx context.Context, userID string) ([]*model.Home, error) {
	client, ctx := getDbClient()
	user := client.Collection(USER_COLLECTION_NAME).Doc(userID)
	homeDocsIter := user.Collection(HOME_COLLECTION_NAME).Documents(ctx)
	var retArray []*model.Home
	for {
		doc, err := homeDocsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		docData := doc.Data()
		homeEntry := &model.Home{
			HomeID: doc.Ref.ID,
			Name: docData[NAME_FIELD_NAME].(string),
			Address: docData[ADDRESS_FIELD_NAME].(string),
			Rent: docData[RENT_FIELD_NAME].(float64),
			NumBed: int(docData[BEDS_FIELD_NAME].(int64)),
			City: docData[CITY_FIELD_NAME].(string),
			Latitude: docData[LATITUDE_FIELD_NAME].(float64),
			Longitude: docData[LONGITUDE_FIELD_NAME].(float64),

		}
		retArray = append(retArray, homeEntry)
	}
	defer client.Close()
	return retArray, nil
}

func (r *queryResolver) Works(ctx context.Context, userID string) ([]*model.Work, error) {
	client, ctx := getDbClient()
	user := client.Collection(USER_COLLECTION_NAME).Doc(userID)
	homeDocsIter := user.Collection(WORK_COLLECTION_NAME).Documents(ctx)
	var retArray []*model.Work
	for {
		doc, err := homeDocsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		docData := doc.Data()
		homeEntry := &model.Work{
			WorkID: doc.Ref.ID,
			Name: docData[NAME_FIELD_NAME].(string),
			Address: docData[ADDRESS_FIELD_NAME].(string),
			City: docData[CITY_FIELD_NAME].(string),
			Latitude: docData[LATITUDE_FIELD_NAME].(float64),
			Longitude: docData[LONGITUDE_FIELD_NAME].(float64),

		}
		retArray = append(retArray, homeEntry)
	}
	defer client.Close()
	return retArray, nil
}
func (r *queryResolver) Authentication(ctx context.Context, authenticaionDetails model.AuthQuery) (*model.AuthResult, error) {
	client, dbctx := getDbClient()
	var user *firestore.DocumentSnapshot

	if authenticaionDetails.Email != "" {
		docs := client.Collection(USER_COLLECTION_NAME).Where(EMAIL_FIELD_NAME, "==", authenticaionDetails.Email).Documents(dbctx)
		user,  _ = docs.Next()
		defer docs.Stop()
	} else if authenticaionDetails.Username != "" {
		docs := client.Collection(USER_COLLECTION_NAME).Where(USERNAME_FIELD_NAME, "==", authenticaionDetails.Username).Documents(dbctx)
		user,  _ = docs.Next()
		defer docs.Stop()

	} else {
		defer client.Close()
		return &model.AuthResult{
			Success:      false,
			VaildRequest: false,
			UserID:       "",
		}, nil
	}
	if user.Data()[PASSWORD_FIELD_NAME] == GetMD5Hash(authenticaionDetails.Password) {
		defer client.Close()
		return &model.AuthResult{
			Success:      true,
			VaildRequest: true,
			UserID:       user.Ref.ID,
		}, nil
	} else {
		defer client.Close()
		return &model.AuthResult{
			Success:      false,
			VaildRequest: true,
			UserID:       "",
		}, nil
	}


}

func (r *mutationResolver) DeleteItem(ctx context.Context, input *model.DeleteInfo) (*model.DeleteResult, error)  {
	client, dbctx := getDbClient()
	if input.ItemType == WORK_TYPE_NAME {
		_, err := client.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(WORK_COLLECTION_NAME).Doc(input.ItemID).Delete(dbctx)
		if err != nil {
			fmt.Println(err.Error())
			defer client.Close()
			return &model.DeleteResult{
				Success:     false,
				Description: "failed to delete item",
			}, nil
		}
	} else if input.ItemType == HOME_TYPE_NAME {
		_, err := client.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(HOME_COLLECTION_NAME).Doc(input.ItemID).Delete(dbctx)
		if err != nil {
			fmt.Println(err.Error())
			defer client.Close()
			return &model.DeleteResult{
				Success:     false,
				Description: "failed to delete item",
			}, nil
		}
	}
	return &model.DeleteResult{
		Success:     true,
		Description: "Deleted item successfully",
	} ,nil
}
func (r *mutationResolver) ModifyUser(ctx context.Context, input model.ModUser) (*model.User, error) {
	client , _:= getDbClient()

	_, err := client.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Set(ctx, map[string]interface{}{
		USERNAME_FIELD_NAME:     input.Username,
		EMAIL_FIELD_NAME: input.Email,
		PASSWORD_FIELD_NAME:     input.Password,
	}, firestore.MergeAll)
	if err != nil {
		defer client.Close()
		return nil, err
	}
	newDoc, err := client.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Get(ctx)
	docData := newDoc.Data()
	if err != nil {
		defer client.Close()
		return nil, err
	}
	result := &model.User{
		UserID:   input.UserID,
		Password: docData[PASSWORD_FIELD_NAME].(string),
		Username: docData[USERNAME_FIELD_NAME].(string),
		Email:    docData[EMAIL_FIELD_NAME].(string),
	}
	defer client.Close()
	return result, nil

}

func (r *mutationResolver) ModifyHome(ctx context.Context, input model.ModHome) (*model.Home, error) {
	dbclient , _:= getDbClient()
	client, _ := getMapsClient()
	//get the lat and long
	results, err := client.Geocode(ctx, &maps.GeocodingRequest{Address: input.Address})
	if err != nil {
		fmt.Println("api request error")
	}

	result := results[0]
	latitude := result.Geometry.Location.Lat
	longitude := result.Geometry.Location.Lng

	_, dberr := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(HOME_COLLECTION_NAME).Doc(input.HomeID).Set(ctx, map[string]interface{}{
		NAME_FIELD_NAME:      input.Name,
		ADDRESS_FIELD_NAME:   input.Address,
		RENT_FIELD_NAME:      input.Rent,
		BEDS_FIELD_NAME:   input.NumBed,
		CITY_FIELD_NAME:      input.City,
		LATITUDE_FIELD_NAME:  latitude,
		LONGITUDE_FIELD_NAME: longitude,
	}, firestore.MergeAll)
	if dberr != nil {
		defer dbclient.Close()
		return nil, err
	}
	newDoc, dberr := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(HOME_COLLECTION_NAME).Doc(input.HomeID).Get(ctx)
	docData := newDoc.Data()
	if err != nil {
		defer dbclient.Close()
		return nil, err
	}
	homeEntry := &model.Home{
		HomeID: newDoc.Ref.ID,
		Name: docData[NAME_FIELD_NAME].(string),
		Address: docData[ADDRESS_FIELD_NAME].(string),
		Rent: docData[RENT_FIELD_NAME].(float64),
		NumBed: docData[BEDS_FIELD_NAME].(int),
		City: docData[CITY_FIELD_NAME].(string),
		Latitude: docData[LATITUDE_FIELD_NAME].(float64),
		Longitude: docData[LONGITUDE_FIELD_NAME].(float64),
	}
	defer dbclient.Close()
	return homeEntry, nil

}

func (r *mutationResolver) ModifyWork(ctx context.Context, input model.ModWork) (*model.Work, error) {
	dbclient , _:= getDbClient()
	client, _ := getMapsClient()
	//get the lat and long
	results, err := client.Geocode(ctx, &maps.GeocodingRequest{Address: input.Address})
	if err != nil {
		fmt.Println("api request error")
	}

	result := results[0]
	latitude := result.Geometry.Location.Lat
	longitude := result.Geometry.Location.Lng

	_, dberr := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(WORK_COLLECTION_NAME).Doc(input.WorkID).Set(ctx, map[string]interface{}{
		NAME_FIELD_NAME:      input.Name,
		ADDRESS_FIELD_NAME:   input.Address,
		CITY_FIELD_NAME:      input.City,
		LATITUDE_FIELD_NAME:  latitude,
		LONGITUDE_FIELD_NAME: longitude,
	}, firestore.MergeAll)
	if dberr != nil {
		defer dbclient.Close()
		return nil, err
	}
	newDoc, dberr := dbclient.Collection(USER_COLLECTION_NAME).Doc(input.UserID).Collection(WORK_COLLECTION_NAME).Doc(input.WorkID).Get(ctx)
	docData := newDoc.Data()
	if err != nil {
		defer dbclient.Close()
		return nil, err
	}
	workEntry := &model.Work{
		WorkID: newDoc.Ref.ID,
		Name: docData[NAME_FIELD_NAME].(string),
		Address: docData[ADDRESS_FIELD_NAME].(string),
		City: docData[CITY_FIELD_NAME].(string),
		Latitude: docData[LATITUDE_FIELD_NAME].(float64),
		Longitude: docData[LONGITUDE_FIELD_NAME].(float64),
	}
	defer dbclient.Close()
	return workEntry, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }






