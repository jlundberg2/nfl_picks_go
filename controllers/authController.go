package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jlundberg2/nfl_picks_go/database"
	"github.com/jlundberg2/nfl_picks_go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var userCollection = database.DB().Database("nflPicks").Collection("users")
var secretKey = os.Getenv("SECRET_KEY")

func Register(c *fiber.Ctx) error {
	var data models.User
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 16)
	if err != nil {
		return err
	}
	data.Password = string(password)

	_, err = userCollection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(data)
}

func Login(c *fiber.Ctx) error {
	var data models.User

	err := c.BodyParser(&data)
	if err != nil {
		fmt.Println("Error parsing input data")
		return err
	}

	filter := bson.D{{"email", data.Email}}
	var result models.User
	err = userCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error. Unable to access database",
		})

	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(data.Password))
	if err != nil {
		fmt.Println("Incorrect password")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect Password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(result.Id),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
        fmt.Println("Unable to sign token")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error. Unable sign token",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Successfully Logged In",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims

	return c.JSON(claims)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func Check(c *fiber.Ctx) error {
	var data models.User
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	filter := bson.D{{"email", data.Email}}
	var result models.User
	err = userCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error. Unable to access database",
		})

	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "userExists",
	})
}
