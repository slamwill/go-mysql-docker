package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"database/sql"
    "fmt"
    "time"
	"strconv"
	"strings"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
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

	e.GET("/getdb", func(c echo.Context) error {
		retStr := getSingleDataFromDBold()
		return c.JSON(http.StatusOK, struct{ Username string } {Username: "Hi " + retStr})
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

			// fmt.Printf("id is %d", id)

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

func getSingleDataFromDBold() string {
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

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		var retStr = "failed"
		return retStr
	}
    DB.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	DB.SetMaxOpenConns(MaxOpenConns)
	DB.SetMaxIdleConns(MaxIdleConns)
	row := DB.QueryRow("select id,name from users where id=?", 4814546)

	var id int
	var name string
	err = row.Scan(&id, &name)
	// checkErr(err)
	var retStr = name
	defer DB.Close()
	return retStr
}
