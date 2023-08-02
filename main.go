package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"plsc/golang/handlers"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	//Load enviroment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// Connect to PlanetScale database using DSN environment variable.
	db, err := getDbConn()
	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
	}

	// Create an API handler which serves data from PlanetScale.
	handler := NewHandler(db)

	// Start an HTTP API server.
	const addr = ":8080"
	log.Printf("successfully connected to PlanetScale, starting HTTP server on %q", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

func getDbConn() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	// Connect to PlanetScale database using DSN environment variable.
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})

	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
		return nil, err
	}
	tx := db.Session(&gorm.Session{Logger: newLogger})

	return tx, nil
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// NewHandler creates an http.Handler which wraps a PlanetScale database
// connection.
func NewHandler(db *gorm.DB) http.Handler {

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("frontEndHost")},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	router.HandleFunc("/health", HealthCheckHandler)
	productHandler := handlers.NewProductHandler(*db)
	router.HandleFunc("/products", productHandler.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/add/product", productHandler.AddProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", productHandler.GetProduct).Methods(http.MethodGet)
	router.HandleFunc("/categories", productHandler.GetCategories).Methods(http.MethodGet)
	router.HandleFunc("/categories/{id}", productHandler.GetCategory).Methods(http.MethodGet)
	router.HandleFunc("/query", productHandler.GetRawQuery).Methods(http.MethodGet)

	return handler
}

// // // seedDatabase is the HTTP handler for GET /seed.
// func (h *Handler) seedDatabase(w http.ResponseWriter, r *http.Request) {
// 	// Perform initial schema migrations.
// 	if err := h.db.AutoMigrate(&Product{}); err != nil {
// 		http.Error(w, "failed to migrate products table", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := h.db.AutoMigrate(&Category{}); err != nil {
// 		http.Error(w, "failed to migrate categories table", http.StatusInternalServerError)
// 		return
// 	}

// 	// Seed categories and products for those categories.
// 	h.db.Create(&Category{
// 		Name:        "Phone",
// 		Description: "Description 1",
// 	})
// 	h.db.Create(&Category{
// 		Name:        "Video Game Console",
// 		Description: "Description 2",
// 	})

// 	h.db.Create(&Product{
// 		Name:        "iPhone",
// 		Description: "Description 1",
// 		Image:       "Image 1",
// 		Category:    Category{ID: 1},
// 	})
// 	h.db.Create(&Product{
// 		Name:        "Pixel Pro",
// 		Description: "Description 2",
// 		Image:       "Image 2",
// 		Category:    Category{ID: 1},
// 	})
// 	h.db.Create(&Product{
// 		Name:        "Playstation",
// 		Description: "Description 3",
// 		Image:       "Image 3",
// 		Category:    Category{ID: 2},
// 	})
// 	h.db.Create(&Product{
// 		Name:        "Xbox",
// 		Description: "Description 4",
// 		Image:       "Image 4",
// 		Category:    Category{ID: 2},
// 	})
// 	h.db.Create(&Product{
// 		Name:        "Galaxy S",
// 		Description: "Description 5",
// 		Image:       "Image 5",
// 		Category:    Category{ID: 1},
// 	})

// 	io.WriteString(w, "Migrations and Seeding of database complete\n")
// }
