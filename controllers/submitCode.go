package controllers

import (
	"CfBE/initializers"
	"CfBE/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var apiKey = (os.Getenv("API_KEY"))

const apiHost = "judge0-ce.p.rapidapi.com"

// RequestPayload defines the structure for input payload from the user
type RequestPayload struct {
	LanguageID int    `json:"language_id"` // User provides language ID
	SourceCode string `json:"source_code"` // User provides source code
}

// SubmitCode handles code submission requests
func SubmitCode(c *gin.Context) {
	// cehck if the problem is exist

	// get data from url:
	userIDStr := c.Param("id1")
	var user models.User

	userID, err := strconv.ParseUint(userIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	// --------------------------------------------------------------------------
	contestIDStr := c.Param("id2")
	var contest models.Contest

	contestID, err := strconv.ParseUint(contestIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	if err := initializers.DB.First(&contest, contestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contest not found"})
		return
	}
	// --------------------------------------------------------------------------
	problemIDStr := c.Param("id3")
	var problem models.Problem

	problemID, err := strconv.ParseUint(problemIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}
	if err := initializers.DB.Where("problem_id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problem not found"})
		return
	}

	var userPayload RequestPayload

	// Bind JSON from the request to our struct
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	// Ensure language_id and source_code are provided
	if userPayload.LanguageID == 0 || userPayload.SourceCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "language_id and source_code must be provided"})
		return
	}

	TestCaseInput := problem.TestCaseInput

	// Base64 encode the user-provided source code and input
	encodedSourceCode := base64.StdEncoding.EncodeToString([]byte(userPayload.SourceCode))
	encodedInput := base64.StdEncoding.EncodeToString([]byte(TestCaseInput))

	// Prepare the payload with user-provided source code and input
	payload := strings.NewReader(fmt.Sprintf(`{
		"language_id": %d,
		"source_code": "%s",
		"stdin": "%s"
	}`, userPayload.LanguageID, encodedSourceCode, encodedInput))

	// Prepare the request to send to Judge0 API
	url := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=false&fields=*"
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)
	req.Header.Add("Content-Type", "application/json")

	// Send request to Judge0 API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer res.Body.Close()

	// Read the response from Judge0
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Log the full Judge0 response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Judge0 response", "response": string(body)})
		return
	}

	// Check for the token
	token, ok := result["token"].(string)
	if !ok {
		if message, ok := result["message"]; ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not found in response", "response": result})
		}
		return
	}

	// Poll for the result using the token
	var finalResult map[string]interface{}
	for {
		finalResult, err = getResult(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get result"})
			return
		}

		// Check if the status is ready
		statusID := finalResult["status"].(map[string]interface{})["id"].(float64)
		if statusID != 1 && statusID != 2 { // 1 = "In Queue", 2 = "Processing"
			break
		}

		// Wait before the next request
		time.Sleep(2 * time.Second)
	}

	// Save the result to the database using initializers.DB
	newResult := models.Result{
		UserID:            uint(userID),
		ContestID:         uint(contestID),
		ProblemID:         uint(problemID),
		LanguageID:        userPayload.LanguageID,
		SourceCode:        userPayload.SourceCode,
		Input:             TestCaseInput,
		Output:            safeString(finalResult["stdout"]),
		ErrorOutput:       safeString(finalResult["stderr"]),
		StatusID:          uint(finalResult["status"].(map[string]interface{})["id"].(float64)),
		StatusDescription: finalResult["status"].(map[string]interface{})["description"].(string),
	}

	// Handle memory used
	if memoryValue, ok := finalResult["memory"]; ok && memoryValue != nil {
		if memoryFloat, ok := memoryValue.(float64); ok {
			newResult.MemoryUsed = int(memoryFloat)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected type for memory used"})
			return
		}
	} else {
		newResult.MemoryUsed = 0 // Default value
	}

	// Handle execution time type
	/*
		var executionTime float64
		if timeValue, ok := finalResult["time"]; ok {
			switch v := timeValue.(type) {
			case float64:
				executionTime = v
			case string:
				if parsedTime, err := strconv.ParseFloat(v, 64); err == nil {
					executionTime = parsedTime
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid execution time format"})
					return
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected type for execution time"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Execution time not found in result"})
			return
		}
		newResult.ExecutionTime = executionTime

		// Handle wall_time type
		var wallTime float64
		if wallTimeValue, ok := finalResult["wall_time"]; ok {
			switch v := wallTimeValue.(type) {
			case float64:
				wallTime = v
			case string:
				if parsedTime, err := strconv.ParseFloat(v, 64); err == nil {
					wallTime = parsedTime
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid wall time format"})
					return
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected type for wall time"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Wall time not found in result"})
			return
		}
		newResult.WallTime = wallTime

		// Handle exit code
		if exitCodeValue, ok := finalResult["exit_code"]; ok && exitCodeValue != nil {
			if exitCodeFloat, ok := exitCodeValue.(float64); ok {
				newResult.ExitCode = int(exitCodeFloat)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected type for exit code"})
				return
			}
		} else {
			newResult.ExitCode = 0 // Default value
		}
	*/
	// Save the result to the database

	// Show the status description
	// Decode the output
	decodedOutput, err := base64.StdEncoding.DecodeString(safeString(newResult.Output))
	if err != nil {
		// Handle error if decoding fails
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decode the output",
		})
		return
	}

	var resultMessage string
	if problem.TestCaseOutput == string(decodedOutput) { // Convert byte slice to string for comparison
		resultMessage = "Test case passed"
		newResult.ResultMessage = "Test case passed"
	} else {
		resultMessage = "Test case failed"
		newResult.ResultMessage = "Test case failed"

	}
	if err := initializers.DB.Create(&newResult).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Code":           newResult.StatusDescription,
		"TestCaseResult": resultMessage,
	})

}

// getResult fetches the result of a submission using the token
func getResult(token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s?base64_encoded=true&fields=*", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// safeString safely converts an interface{} to a string
func safeString(value interface{}) string {
	if value == nil {
		return ""
	}
	strValue, ok := value.(string)
	if !ok {
		return ""
	}
	return strValue
}
