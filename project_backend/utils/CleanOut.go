package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)








func Cleancode(output []byte) string {
	// Clean the output as its from the chaincode (terminal)
	// find the payload if present
	// remove the first 6 characters and last 2 characters
	// and return the string
	outputstring := string(output)
	log.Println("Output string", outputstring)
	// find the first occurrence of "payload" and remove everything before it
	start := 0
	for i := 0; i < len(outputstring); i++ {
		if outputstring[i] == '{' {
			start = i
			break
		}
	}
	end := len(outputstring) - 2
	for i := len(outputstring) - 1; i >= 0; i-- {
		if outputstring[i] == '}' {
			end = i
			break
		}
	}

	return outputstring[start : end+1]
}

func Cleancode2(c *gin.Context, output []byte) map[string]interface{} {

	log.Println("First output", string(output))

	// Clean the output as its from the chaincode (terminal)
	// extract the payload if present after cleaning
	outputstring := Cleancode(output)
	log.Println("Second output", outputstring)

	// now this string contains (/")'s as read from terminal we need to make them normal
	newcleanstring := ""
	for i := 0; i < len(outputstring); i++ {
		if outputstring[i] == '"' {
			newcleanstring += `"`
		} else if outputstring[i] == '\\' {
			// newcleanstring += `\`
			continue
		} else {
			newcleanstring += string(outputstring[i])
		}
	}
	log.Println("Third output", newcleanstring)
	// now unmarshal the string to a map
	var result map[string]interface{}
	err := json.Unmarshal([]byte(newcleanstring), &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal JSON"})
		return nil
	}
	log.Println("Fourth output", result)
	return result
}
