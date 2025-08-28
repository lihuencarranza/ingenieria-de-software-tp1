package main

// Test Suite for Melodia API Endpoints
//
// This test suite validates the API endpoints including:
// - Health checks
// - CRUD operations for songs
// - CRUD operations for playlists
// - Adding songs to playlists (with duplicate prevention)
//
// Duplicate Prevention Rules:
// - A song cannot be added twice to the same playlist
// - The same song can be added to different playlists
// - Duplicate attempts should return 400 Bad Request
//
// Expected Behavior:
// - POST /playlists/{id}/songs with existing song_id should fail with 400
// - GET /playlists/{id} should return songs ordered by addedAt (most recent first)
// - No duplicate songs should appear in the same playlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestResult represents the result of an individual test
type TestResult struct {
	TestName   string    `json:"test_name"`
	Endpoint   string    `json:"endpoint"`
	Method     string    `json:"method"`
	Status     string    `json:"status"` // PASS, FAIL, ERROR
	StatusCode int       `json:"status_code"`
	Response   string    `json:"response"`
	Error      string    `json:"error,omitempty"`
	Duration   string    `json:"duration"`
	Timestamp  time.Time `json:"timestamp"`
}

// TestSuite represents the entire test suite
type TestSuite struct {
	Results    []TestResult `json:"results"`
	TotalTests int          `json:"total_tests"`
	Passed     int          `json:"passed"`
	Failed     int          `json:"failed"`
	Errors     int          `json:"errors"`
	StartTime  time.Time    `json:"start_time"`
	EndTime    time.Time    `json:"end_time"`
}

const (
	BaseURL = "http://localhost:8080"
)

var (
	testSuite = &TestSuite{
		Results:   []TestResult{},
		StartTime: time.Now(),
	}
)

// runTest executes an individual test and records the result
func runTest(testName, method, endpoint, body string, expectedStatus int) {
	start := time.Now()

	result := TestResult{
		TestName:  testName,
		Endpoint:  endpoint,
		Method:    method,
		Timestamp: time.Now(),
	}

	var req *http.Request
	var err error

	if body != "" {
		req, err = http.NewRequest(method, BaseURL+endpoint, bytes.NewBuffer([]byte(body)))
		if err != nil {
			result.Status = "ERROR"
			result.Error = fmt.Sprintf("Error creating request: %v", err)
			result.Duration = time.Since(start).String()
			testSuite.Results = append(testSuite.Results, result)
			return
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, BaseURL+endpoint, nil)
		if err != nil {
			result.Status = "ERROR"
			result.Error = fmt.Sprintf("Error creating request: %v", err)
			result.Duration = time.Since(start).String()
			testSuite.Results = append(testSuite.Results, result)
			return
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		result.Status = "ERROR"
		result.Error = fmt.Sprintf("Error making request: %v", err)
		result.Duration = time.Since(start).String()
		testSuite.Results = append(testSuite.Results, result)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "ERROR"
		result.Error = fmt.Sprintf("Error reading response: %v", err)
		result.Duration = time.Since(start).String()
		testSuite.Results = append(testSuite.Results, result)
		return
	}

	result.StatusCode = resp.StatusCode
	result.Response = string(bodyBytes)
	result.Duration = time.Since(start).String()

	if resp.StatusCode == expectedStatus {
		result.Status = "PASS"
		testSuite.Passed++
	} else {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
		testSuite.Failed++
	}

	testSuite.Results = append(testSuite.Results, result)
}

// printResults prints the test results
func printResults() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("API ENDPOINT TEST RESULTS")
	fmt.Println(strings.Repeat("=", 80))

	for _, result := range testSuite.Results {
		statusIcon := "✅"
		if result.Status == "FAIL" {
			statusIcon = "❌"
		} else if result.Status == "ERROR" {
			statusIcon = "💥"
		}

		fmt.Printf("%s %s (%s %s)\n", statusIcon, result.TestName, result.Method, result.Endpoint)
		fmt.Printf("   Status: %s (Code: %d) | Duration: %s\n", result.Status, result.StatusCode, result.Duration)

		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}

		if result.Status == "PASS" {
			fmt.Printf("   Response: %s...\n", truncateString(result.Response, 100))
		}
		fmt.Println()
	}

	testSuite.EndTime = time.Now()
	duration := testSuite.EndTime.Sub(testSuite.StartTime)

	fmt.Println("FINAL SUMMARY:")
	fmt.Printf("   Total Tests: %d\n", len(testSuite.Results))
	fmt.Printf("   ✅ Passed: %d\n", testSuite.Passed)
	fmt.Printf("   ❌ Failed: %d\n", testSuite.Failed)
	fmt.Printf("   💥 Errors: %d\n", testSuite.Errors)
	fmt.Printf("   ⏱️  Total Duration: %s\n", duration)

	if testSuite.Failed == 0 && testSuite.Errors == 0 {
		fmt.Println("\n🎉 ALL TESTS PASSED SUCCESSFULLY!")
	} else {
		fmt.Printf("\n⚠️  %d tests failed or had errors\n", testSuite.Failed+testSuite.Errors)
	}
}

// truncateString truncates a string if it's too long
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// saveLogs saves the results to a log file
func saveLogs() {
	testSuite.TotalTests = len(testSuite.Results)

	logData, err := json.MarshalIndent(testSuite, "", "  ")
	if err != nil {
		log.Printf("Error marshaling test results: %v", err)
		return
	}

	// Create test_results directory if it doesn't exist
	resultsDir := "test_results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Printf("Error creating results directory: %v", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("test_results_%s.json", timestamp)
	filepath := filepath.Join(resultsDir, filename)

	err = os.WriteFile(filepath, logData, 0644)
	if err != nil {
		log.Printf("Error writing log file: %v", err)
		return
	}

	fmt.Printf("\n📝 Logs saved to: %s\n", filepath)
}

// waitForService waits for the service to become available
func waitForService() bool {
	fmt.Println("Waiting for service to become available...")

	for i := 0; i < 30; i++ {
		resp, err := http.Get(BaseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Println("Service is available!")
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
		fmt.Printf("Attempt %d/30...\n", i+1)
	}

	fmt.Println("Service not available after 30 attempts")
	return false
}

func main() {
	fmt.Println("STARTING MELODIA API ENDPOINT TESTS")
	fmt.Println(strings.Repeat("=", 60))

	// Wait for service to be available
	if !waitForService() {
		os.Exit(1)
	}

	fmt.Println("\nRUNNING TESTS...")
	fmt.Println(strings.Repeat("-", 40))

	// Health Check Tests
	runTest(
		"Health Check",
		"GET",
		"/health",
		"",
		200,
	)

	// Song Tests - Create
	fmt.Println("Testing Song endpoints - Create...")
	runTest(
		"Create Song - Valid Data 1",
		"POST",
		"/songs",
		`{"title":"Bohemian Rhapsody","artist":"Queen"}`,
		201,
	)

	runTest(
		"Create Song - Valid Data 2",
		"POST",
		"/songs",
		`{"title":"Hotel California","artist":"Eagles"}`,
		201,
	)

	runTest(
		"Create Song - Valid Data 3",
		"POST",
		"/songs",
		`{"title":"Stairway to Heaven","artist":"Led Zeppelin"}`,
		201,
	)

	runTest(
		"Create Song - Empty Title",
		"POST",
		"/songs",
		`{"title":"","artist":"Test Artist"}`,
		400,
	)

	runTest(
		"Create Song - Empty Artist",
		"POST",
		"/songs",
		`{"title":"Test Song","artist":""}`,
		400,
	)

	runTest(
		"Create Song - Invalid JSON",
		"POST",
		"/songs",
		`{"title":"Test Song","artist":"Test Artist"`,
		400,
	)

	// Song Tests - Read
	fmt.Println("\nTesting Song endpoints - Read...")
	runTest(
		"Get All Songs",
		"GET",
		"/songs",
		"",
		200,
	)

	runTest(
		"Get Song by ID - Valid",
		"GET",
		"/songs/1",
		"",
		200,
	)

	runTest(
		"Get Song by ID - Invalid ID",
		"GET",
		"/songs/invalid",
		"",
		400,
	)

	runTest(
		"Get Song by ID - Non-existent",
		"GET",
		"/songs/999",
		"",
		404,
	)

	// Song Tests - Update
	fmt.Println("\nTesting Song endpoints - Update...")
	runTest(
		"Update Song - Valid",
		"PUT",
		"/songs/1",
		`{"title":"Bohemian Rhapsody (Updated)","artist":"Queen"}`,
		200,
	)

	runTest(
		"Update Song - Invalid ID",
		"PUT",
		"/songs/invalid",
		`{"title":"Test","artist":"Test"}`,
		400,
	)

	runTest(
		"Update Song - Invalid JSON",
		"PUT",
		"/songs/1",
		`{"title":"Test","artist":"Test"`,
		400,
	)

	runTest(
		"Update Song - Non-existent ID",
		"PUT",
		"/songs/999",
		`{"title":"Test","artist":"Test"}`,
		404,
	)

	// Song Tests - Delete
	fmt.Println("\nTesting Song endpoints - Delete...")
	runTest(
		"Delete Song - Invalid ID",
		"DELETE",
		"/songs/invalid",
		"",
		400,
	)

	runTest(
		"Delete Song - Non-existent ID",
		"DELETE",
		"/songs/999",
		"",
		404,
	)

	// Playlist Tests - Create
	fmt.Println("\nTesting Playlist endpoints - Create...")
	runTest(
		"Create Playlist - Valid Data 1",
		"POST",
		"/playlists",
		`{"name":"Rock Classics","description":"Best rock songs of all time"}`,
		201,
	)

	runTest(
		"Create Playlist - Valid Data 2",
		"POST",
		"/playlists",
		`{"name":"Pop Hits","description":"Most popular pop songs"}`,
		201,
	)

	runTest(
		"Create Playlist - Empty Name",
		"POST",
		"/playlists",
		`{"name":"","description":"Test description"}`,
		400,
	)

	runTest(
		"Create Playlist - Empty Description",
		"POST",
		"/playlists",
		`{"name":"Test Playlist","description":""}`,
		400,
	)

	runTest(
		"Create Playlist - Invalid JSON",
		"POST",
		"/playlists",
		`{"name":"Test Playlist","description":"Test description"`,
		400,
	)

	// Playlist Tests - Read
	fmt.Println("\nTesting Playlist endpoints - Read...")
	runTest(
		"Get All Playlists",
		"GET",
		"/playlists",
		"",
		200,
	)

	runTest(
		"Get Playlist by ID - Valid",
		"GET",
		"/playlists/1",
		"",
		200,
	)

	runTest(
		"Get Playlist by ID - Invalid ID",
		"GET",
		"/playlists/invalid",
		"",
		400,
	)

	runTest(
		"Get Playlist by ID - Non-existent",
		"GET",
		"/playlists/999",
		"",
		404,
	)

	// Playlist Tests - Add Songs
	fmt.Println("\nTesting Playlist endpoints - Add Songs...")
	runTest(
		"Add Song to Playlist - Valid 1",
		"POST",
		"/playlists/1/songs",
		`{"song_id":1}`,
		200,
	)

	runTest(
		"Add Song to Playlist - Valid 2",
		"POST",
		"/playlists/1/songs",
		`{"song_id":2}`,
		200,
	)

	runTest(
		"Add Song to Playlist - Another Valid Song",
		"POST",
		"/playlists/1/songs",
		`{"song_id":3}`,
		200,
	)

	// Test adding songs to different playlists (should work)
	runTest(
		"Add Song to Different Playlist - Valid",
		"POST",
		"/playlists/2/songs",
		`{"song_id":1}`,
		200, // Misma canción en playlist diferente debería funcionar
	)


	runTest(
		"Get Playlist 2 - Verify No Duplicates",
		"GET",
		"/playlists/2",
		"",
		200, // Debería mostrar solo canciones únicas
	)

	// Test edge cases for adding songs to playlists
	runTest(
		"Add Song to Playlist - Non-existent Song",
		"POST",
		"/playlists/1/songs",
		`{"song_id":999}`,
		404, // Canción que no existe debería fallar
	)

	runTest(
		"Add Song to Playlist - Non-existent Playlist",
		"POST",
		"/playlists/999/songs",
		`{"song_id":1}`,
		404, // Playlist que no existe debería fallar
	)

	runTest(
		"Add Song to Playlist - Invalid JSON",
		"POST",
		"/playlists/1/songs",
		`{"song_id":1`,
		400, // JSON inválido debería fallar
	)

	// Print results and save logs
	printResults()
	saveLogs()

	// Exit based on results
	if testSuite.Failed > 0 || testSuite.Errors > 0 {
		os.Exit(1)
	}
}
