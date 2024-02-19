package routes

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/realtime-tl-be/api/rtserver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

type Employee struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Designation string             `json:"designation,omitempty" bson:"designation,omitempty"`
	Salary      string             `json:"salary,omitempty" bson:"salary,omitempty"`
}

const (
	Dbname         = "simple-backend"
	CollectionName = "collectiontest"
)

var db *mongo.Database

func InitializeRouter(port string) {
	r := gin.Default()

	// to allow for CORS from all origins
	r.Use(cors.Default())

	// Routes
	r.GET("/api/employees", getEmployees)
	r.GET("/api/employees/:id", getEmployeeById)
	r.POST("/api/employees", createEmployee)
	r.PUT("/api/employees/:id", updateEmployee)
	r.DELETE("/api/employees/:id", deleteEmployee)
	r.GET("/api/socket/ws", rtserver.HandleCommunication)

	initClient()

	// run server
	r.Run(":" + port)

}

func initClient() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	db = client.Database(Dbname)

}

var mongoCollection *mongo.Collection

func getEmployees(c *gin.Context) {

	mongoCollection = db.Collection(CollectionName)

	if mongoCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mongo collection is nil"})
	}

	collection, err := mongoCollection.Find(context.Background(), bson.D{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding documents"})
		log.Fatal(err)
	}

	var employees []Employee

	for collection.Next(context.Background()) {
		var employee Employee
		if err := collection.Decode(&employee); err != nil {
			log.Fatal(err)
		}
		employees = append(employees, employee)
	}
	collection.Close(context.Background())

	c.JSON(http.StatusOK, employees)
}

func getEmployeeById(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var employee []Employee
	if err := db.Collection(CollectionName).FindOne(context.Background(), Employee{ID: id}).Decode(&employee); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, employee)
}

func createEmployee(c *gin.Context) {
	var employee Employee

	if err := c.BindJSON(&employee); err != nil {
		log.Fatal(err)
	}

	// insert employee
	result, err := db.Collection(CollectionName).InsertOne(context.Background(), employee)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)

}

func updateEmployee(c *gin.Context) {
	// get id
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var employee Employee

	if err := c.BindJSON(&employee); err != nil {
		log.Fatal(err)
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "phone", Value: employee.Phone},
			{Key: "destination", Value: employee.Designation},
			{Key: "salary", Value: employee.Salary},
		}},
	}

	result, err := db.Collection(CollectionName).UpdateOne(context.Background(), Employee{ID: id}, update)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)

}

func deleteEmployee(c *gin.Context) {
	// get id from param
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	result, err := db.Collection(CollectionName).DeleteOne(context.Background(), Employee{ID: id})

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, result)
}
