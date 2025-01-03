package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/structs"
)

func ConvertArrayToCsv(array []interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ", "), "[]")
}

func ConvertArrayOfStringToCsv(array []string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ", "), "[]")
}

func AddElemntInCsvString(csvString string, element string) string {
	if csvString == "" {
		return element
	} else {
		csvStringArray := strings.Split(csvString, ",")
		csvStringArray = append(csvStringArray, element)
		updtatedCsvStringArray := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(csvStringArray)), ", "), "[]")
		return updtatedCsvStringArray
	}
}

func ConvertStructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

func ConvertCsvToArrayOfString(csv string) []string {
	return strings.Split(csv, ",")
}

func Includes(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveElement(slice []string, elem string) []string {
	for i, v := range slice {
		if v == elem {
			// Remove the element at index i
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // Return the original slice if element is not found
}

func CheckElementMatch(slice1 []string, slice2 []string) bool {
	// Create a map to hold elements of slice2 for faster lookup
	elements := make(map[string]bool)
	for _, el := range slice2 {
		elements[el] = true
	}

	// Iterate through each element in slice1 and check if it exists in slice2
	for _, el := range slice1 {
		if elements[el] {
			return true
		}
	}

	// If no match is found, return false
	return false
}

func ConvertDuplicatesArrtoUniqueArr(arr []string) []string {
	occurred := map[string]bool{}
	result := []string{}
	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if occurred[arr[e]] != true {
			occurred[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

func ConvertToLowerAndIgnoreSpecial(str string) string {
	var result []rune
	for _, char := range str {
		if unicode.IsLetter(char) {
			result = append(result, unicode.ToLower(char))
		} else if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func ConvertInterfaceToString(arr []interface{}) []string {
	var newSlice []string
	for _, url := range arr {
		strURL, ok := url.(string)
		if ok {
			newSlice = append(newSlice, strURL)
		}
	}
	return newSlice
}

func TrimString(str string) string {
	regexPattern := `^"(.*)"$`
	replacement := "$1"
	regex, err := regexp.Compile(regexPattern)

	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return str
	}

	// Replace using the regular expression
	result := regex.ReplaceAllString(str, replacement)
	return result
}

func ArrayDifference(arr1, arr2 []string) []string {
	var result []string
	exists := make(map[string]bool)

	for _, item := range arr2 {
		exists[item] = true
	}

	for _, item := range arr1 {
		if _, ok := exists[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

func RemoveSpacesAndSpecialCharacters(str string) string {
	result := str

	result = strings.ReplaceAll(str, " ", "")
	result = strings.TrimSpace(result)
	result = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(result, "")
	result = strings.ToLower(result)

	return result
}

func IsValidEmail(email string) bool {
	rfc5322 := "(?i)(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
	validRfc5322Regexp := regexp.MustCompile(fmt.Sprintf("^%s*$", rfc5322))

	if !validRfc5322Regexp.MatchString(email) {
		return false
	} else {
		return true
	}
}

func BuildUserName(email string) string {
	result := email

	i := strings.LastIndexByte(email, '@')
	result = email[:i]
	result = RemoveSpacesAndSpecialCharacters(result)

	return result
}

func ReplaceSpacesWithUnderscores(filename string) string {
	return strings.ReplaceAll(filename, " ", "_")
}

// fixMalformedJSON tries to correct common JSON formatting issues in the input string.
func FixMalformedJSON(input string) string {
	if strings.Contains(input, `":false`) {
		input = strings.Replace(input, `:false`, ``, 1)
	}
	if strings.Contains(input, `":true`) {
		input = strings.Replace(input, `:true`, ``, 1)
	}
	if strings.Contains(input, `":null`) {
		input = strings.Replace(input, `:null`, ``, 1)
	}
	return input
}
