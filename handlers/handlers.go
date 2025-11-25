package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Dooform/test-data-api/database"
	"github.com/Dooform/test-data-api/models"
	"github.com/gin-gonic/gin"
)

func ListBoundaries(c *gin.Context) {
	var boundaries []models.AdministrativeBoundary
	if err := database.DB.Find(&boundaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, boundaries)
}

// QueryBoundaries allows querying the administrative_boundaries table with query parameters.
// Example: /query?name1=Bangkok&name2=Phra Nakhon
func QueryBoundaries(c *gin.Context) {
	var boundaries []models.AdministrativeBoundary

	query := database.DB

	if name1 := c.Query("name1"); name1 != "" {
		query = query.Where("name1 = ?", name1)
	}

	if name2 := c.Query("name2"); name2 != "" {
		query = query.Where("name2 = ?", name2)
	}

	if name3 := c.Query("name3"); name3 != "" {
		query = query.Where("name3 = ?", name3)
	}

	if err := query.Find(&boundaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boundaries)
}

// SearchBoundaries performs a full-text and infix search on the administrative_boundaries table.
// Example: /search?q=Bangkok
func SearchBoundaries(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q query parameter is required"})
		return
	}

	// Sanitize the query to be used with to_tsquery
	// Allow Thai characters and spaces
	reg, err := regexp.Compile(`[^a-zA-Z0-9\s\p{Thai}]`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compile regex"})
		return
	}
	sanitizedQuery := reg.ReplaceAllString(q, "")

	if sanitizedQuery == "" {
		c.JSON(http.StatusOK, []models.AdministrativeBoundary{})
		return
	}

	// We will split the query by spaces and join them with '&' to perform an AND search
	// with prefix matching on the last word.
	words := strings.Fields(sanitizedQuery)
	for i, word := range words {
		if i == len(words)-1 {
			words[i] = word + ":*"
		}
	}
	searchQuery := strings.Join(words, " & ")

	likeQuery := "%" + sanitizedQuery + "%"

	var boundaries []models.AdministrativeBoundary

	if err := database.DB.Where("search_vector @@ to_tsquery('thai', ?) OR search_vector @@ to_tsquery('english', ?) OR (name1 || ' ' || name2 || ' ' || name3 || ' ' || name_eng1 || ' ' || name_eng2 || ' ' || name_eng3) ILIKE ?", searchQuery, searchQuery, likeQuery).Limit(10).Find(&boundaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boundaries)
}