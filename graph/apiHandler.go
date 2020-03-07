package graph
import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"googlemaps.github.io/maps"
	"log"

	"cloud.google.com/go/firestore"
)

func getDbClient() (*firestore.Client, context.Context) {
	projectID := PROJECT_ID

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(CREDENTIALS_FILE_PATH))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client, ctx
}
func getMapsClient() (*maps.Client, context.Context) {
	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(GOOGLE_PLACES_API_KEY))
	if err != nil {
		fmt.Println("error connecting to places api")
	}
	return client, ctx
}
//func createUser(username string, userIDNum int, emailAddress string, password string) firestore.DocumentRef{}
