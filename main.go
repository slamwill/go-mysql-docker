package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"database/sql"
    "fmt"
    _ "time"
	"strconv"
	"strings"
    _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv/autoload"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {

		var envs map[string]string
		envs, err := godotenv.Read(".env")
		if err != nil {
			// log.Fatal("Error loading .env file")
		}
		mysql_root_password := envs["MYSQL_ROOT_PASSWORD"]
		db_user := envs["DB_USER"]
		db_password := envs["DB_PASSWORD"]
		// fmt.Printf("%s uses %s\n", db_user, db_password)
		fmt.Println("11111")
		fmt.Println(mysql_root_password)
		fmt.Println(db_user)
		fmt.Println(db_password)
		fmt.Println("11111")


		// godotenv.Load()
		// fmt.Println("11111")
		// os.Setenv("DB_USER", "admin1234")
		// fmt.Printf("os.Getenv(): %s=%s\n", "DB_USER", os.Getenv("DB_USER"))
		// fmt.Printf("os.Getenv(): %s=%s\n", "DB_PASSWORD", os.Getenv("DB_PASSWORD"))

		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/someroute", func(c echo.Context) error {
		// Get the value of the "id" parameter from the query string
		id := c.QueryParam("id")
		return c.JSON(http.StatusOK, struct {
			Status string
			ID     string // Add a field for the "id" parameter value
		}{
			Status: "OK",
			ID:     id, // Set the "id" field to the parameter value
		})
	})

	e.GET("/test", func(c echo.Context) error {
		retStr := hello()
		fmt.Println(retStr)
		return c.JSON(http.StatusOK, struct{ Status string }{Status: retStr})
	})

	e.GET("/getSingleDataFromDB", func(c echo.Context) error {
		userID := c.QueryParam("userID")
		intUserID, err := strconv.Atoi(userID)
		name, err := getSingleDataFromDB(intUserID)
		_ = err
		retStr := "Hi routeGetSingleDataFromDB " + name;
		return c.JSON(http.StatusOK, struct{ Status string }{Status: retStr})
	})

	e.GET("/getMultipleDataFromDB", func(c echo.Context) error {
		// Get the array of ID strings from the "ids" parameter
		idStrings, ok := c.QueryParams()["ids[]"]
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'ids' parameter"})
		}
	
		// Convert the array of strings to an array of integers
		ids := make([]int, len(idStrings))
		for i, idString := range idStrings {
			id, err := strconv.Atoi(idString)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID: " + idString})
			}

			ids[i] = id
		}
	
		// Call the function to retrieve data from the database
		data, err := getMultipleDataFromDB(ids)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve data from database"})
		}
	
		// Return the data as a JSON response
		return c.JSON(http.StatusOK, data)
	})
	

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func getMultipleDataFromDB(ids []int) ([]string, error) {
    // Open a connection to the database
	db, err := connectToDatabase()
	if err != nil {
		// Handle error
	}
    defer db.Close()

    // Construct a query to retrieve the names of the users with the specified IDs
    query := fmt.Sprintf("SELECT name FROM users WHERE id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]"))

    // Execute the query
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Retrieve the names of the users from the query results
    names := make([]string, 0)
    for rows.Next() {
        var name string
        err := rows.Scan(&name)
        if err != nil {
            return nil, err
        }
        names = append(names, name)
    }

    // Check for any errors that occurred during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return names, nil
}

func getSingleDataFromDB(userID int) (string, error) {
	db, err := connectToDatabase()
	if err != nil {
		// Handle error
	}
	name, err := getUserInfo(db, userID)
	if err != nil {
		// Handle error
	}
	defer db.Close()
	return name, nil
}

func connectToDatabase() (*sql.DB, error) {

	const (
		UserName     string = "vt_test_0906"
		Password     string = "JSnPK7eduQ8v8QCA7YdynDxVrewgb3yP"
		Addr         string = "34.80.26.167"
		Port         int    = 3306
		Database     string = "voicetube"
		MaxLifetime  int    = 10
		MaxOpenConns int    = 10
		MaxIdleConns int    = 10
	)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
    db, err := sql.Open("mysql", connectionString)
    if err != nil {
        return nil, err
    }
    return db, nil
}

func getUserInfo(db *sql.DB, userID int) (string, error) {
    var name string
    err := db.QueryRow("SELECT name FROM users WHERE id=?", userID).Scan(&name)
    if err != nil {
        return "", err
    }
    return name, nil
}

func hello() string {
	
	var str = "Hello, World!"
	return str
}
