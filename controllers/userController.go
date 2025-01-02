package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/database"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/helpers"
	"github.com/razyneko/jwt-auth-with-go-gin-gonic-mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Initialize the user collection from the database
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// Initialize the validator for struct validation
var validate = validator.New()

// HashPassword hashes user entered password using bcrypt package
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Hashing with cost of 14
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// VerifyPassword compares a hashed password with the user provided password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = "Email or Password is incorrect"
		check = false
	}
	return check, msg
}

// SignUp handles user registration
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Parse the JSON body into the user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the parsed user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if the email already exists in the database
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for email"})
			return
		}

		if count > 0{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
			return
		}

		// Hash the user's password and replace it in the struct
		password := HashPassword(*user.Password)
		user.Password = &password

		// Check if the phone number already exists in the database
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for phone"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone already exists"})
			return
		}

		// Set metadata for the user
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()           // Generate a new MongoDB ObjectID
		user.UserId = user.ID.Hex()                 // Convert ObjectID to string
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		// Insert the user into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Return the insertion result
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// Login handles user authentication
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout of 100 seconds
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User       // Incoming login data
		var foundUser models.User  // Retrieved user from the database

		// Parse the JSON body into the user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find the user by email
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		// Verify the password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Generate new tokens for the user
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		// Retrieve the updated user data
		err = userCollection.FindOne(ctx, bson.M{"userId": foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the user data
		c.JSON(http.StatusOK, foundUser)
	}
}

// GetUsers retrieves a paginated list of users from the database
// Supports pagination using `recordPerPage` and `page` query parameters
// Aggregates data to include the total count and paginated user records
func GetUsers() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Parse `recordPerPage` query parameter, defaulting to 10 if invalid or not provided
        recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
        if err != nil || recordPerPage < 1 {
            recordPerPage = 10
        }

        // Parse `page` query parameter, defaulting to 1 if invalid or not provided
        page, err := strconv.Atoi(c.Query("page"))
        if err != nil || page < 1 {
            page = 1
        }

        // Calculate the starting index for the page
        startIndex := (page - 1) * recordPerPage

        // MongoDB aggregation pipeline:
        // 1. `matchStage` filters data (currently matches all documents)
        // 2. `groupStage` groups documents and calculates total count and all user data
        // 3. `projectStage` slices the user data to return only paginated records
        matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
    	{Key: "_id", Value: "null"}, // Group all documents into one group
    	{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: 1}}}, // Total number of documents
    	{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}, // Push all documents into an array
		}}}
		projectStage := bson.D{
    	{Key: "$project", Value: bson.D{
        {Key: "_id", Value: 0}, // Exclude `_id` from the result
        {Key: "totalCount", Value: 1}, // Include total count of documents
        {Key: "userItems", Value: bson.D{
            {Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}, // Paginated user data
        }},
    }},
}


        // Execute the aggregation pipeline and retrieve results
        result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
            matchStage, groupStage, projectStage,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while listing user items"})
            return
        }

        // Parse the aggregation result into a slice of maps
        var allUsers []bson.M
        if err = result.All(ctx, &allUsers); err != nil {
            log.Fatal(err)
        }

        // Return the first document from the result as it contains the aggregated data
        c.JSON(http.StatusOK, allUsers[0])
    }
}

// GetUser retrieves a single user by their user ID
// The user ID is extracted from the URL parameter `userId`
// Returns user details if found, otherwise returns an error
func GetUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Extract `userId` parameter from the URL
        userId := c.Param("userId")

        // Find the user in the database using the `userId`
        var user models.User
        err := userCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the user"})
            return
        }

        // Return the user details if found
        c.JSON(http.StatusOK, user)
    }
}
