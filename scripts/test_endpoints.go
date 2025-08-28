package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TestResult representa el resultado de un test individual
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

// TestSuite representa toda la suite de tests
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

// runTest ejecuta un test individual y registra el resultado
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

// printResults imprime los resultados de los tests
func printResults() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("RESULTADOS DE LOS TESTS DE ENDPOINTS")
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

	fmt.Println("RESUMEN FINAL:")
	fmt.Printf("   Total Tests: %d\n", len(testSuite.Results))
	fmt.Printf("   ✅ Passed: %d\n", testSuite.Passed)
	fmt.Printf("   ❌ Failed: %d\n", testSuite.Failed)
	fmt.Printf("   💥 Errors: %d\n", testSuite.Errors)
	fmt.Printf("   ⏱️  Total Duration: %s\n", duration)

	if testSuite.Failed == 0 && testSuite.Errors == 0 {
		fmt.Println("\n🎉 ¡TODOS LOS TESTS PASARON EXITOSAMENTE!")
	} else {
		fmt.Printf("\n⚠️  %d tests fallaron o tuvieron errores\n", testSuite.Failed+testSuite.Errors)
	}
}

// truncateString trunca una cadena si es muy larga
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// saveLogs guarda los resultados en un archivo de log
func saveLogs() {
	testSuite.TotalTests = len(testSuite.Results)

	logData, err := json.MarshalIndent(testSuite, "", "  ")
	if err != nil {
		log.Printf("Error marshaling test results: %v", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("test_results_%s.json", timestamp)

	err = os.WriteFile(filename, logData, 0644)
	if err != nil {
		log.Printf("Error writing log file: %v", err)
		return
	}

	fmt.Printf("\n📝 Logs guardados en: %s\n", filename)
}

// waitForService espera a que el servicio esté disponible
func waitForService() bool {
	fmt.Println("Esperando a que el servicio esté disponible...")

	for i := 0; i < 30; i++ {
		resp, err := http.Get(BaseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Println("Servicio disponible!")
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
		fmt.Printf("Intento %d/30...\n", i+1)
	}

	fmt.Println("Servicio no disponible después de 30 intentos")
	return false
}

func main() {
	fmt.Println("INICIANDO TESTS DE ENDPOINTS DE LA API MELODIA")
	fmt.Println(strings.Repeat("=", 60))

	// Esperar a que el servicio esté disponible
	if !waitForService() {
		os.Exit(1)
	}

	fmt.Println("\nEJECUTANDO TESTS...")
	fmt.Println(strings.Repeat("-", 40))

	// Tests de Songs
	fmt.Println("Testing Songs endpoints...")

	// Crear canciones
	runTest(
		"Create Song - Bohemian Rhapsody",
		"POST",
		"/songs",
		`{"title":"Bohemian Rhapsody","artist":"Queen"}`,
		201,
	)

	runTest(
		"Create Song - Hotel California",
		"POST",
		"/songs",
		`{"title":"Hotel California","artist":"Eagles"}`,
		201,
	)

	runTest(
		"Create Song - Stairway to Heaven",
		"POST",
		"/songs",
		`{"title":"Stairway to Heaven","artist":"Led Zeppelin"}`,
		201,
	)

	// Obtener canciones
	runTest(
		"Get All Songs",
		"GET",
		"/songs",
		"",
		200,
	)

	runTest(
		"Get Song by ID",
		"GET",
		"/songs/1",
		"",
		200,
	)

	// Actualizar canción
	runTest(
		"Update Song",
		"PUT",
		"/songs/1",
		`{"title":"Bohemian Rhapsody (Updated)","artist":"Queen"}`,
		200,
	)

	// Tests de validación de Songs
	runTest(
		"Create Song - Empty Title (Should Fail)",
		"POST",
		"/songs",
		`{"title":"","artist":"Test Artist"}`,
		400,
	)

	runTest(
		"Create Song - Empty Artist (Should Fail)",
		"POST",
		"/songs",
		`{"title":"Test Song","artist":""}`,
		400,
	)

	runTest(
		"Create Song - Invalid JSON (Should Fail)",
		"POST",
		"/songs",
		`{"title":"Test Song","artist":"Test Artist"`,
		400,
	)

	// Tests de Playlists
	fmt.Println("\nTesting Playlists endpoints...")

	// Crear playlists
	runTest(
		"Create Playlist - Rock Classics",
		"POST",
		"/playlists",
		`{"name":"Rock Classics","description":"Las mejores canciones de rock"}`,
		201,
	)

	runTest(
		"Create Playlist - Pop Hits",
		"POST",
		"/playlists",
		`{"name":"Pop Hits","description":"Canciones pop más populares"}`,
		201,
	)

	// Obtener playlists
	runTest(
		"Get All Playlists",
		"GET",
		"/playlists",
		"",
		200,
	)

	runTest(
		"Get Playlist by ID",
		"GET",
		"/playlists/1",
		"",
		200,
	)

	// Agregar canciones a playlists
	runTest(
		"Add Song to Playlist",
		"POST",
		"/playlists/1/songs",
		`{"song_id":1}`,
		200,
	)

	runTest(
		"Add Another Song to Playlist",
		"POST",
		"/playlists/1/songs",
		`{"song_id":2}`,
		200,
	)

	// Tests de validación de Playlists
	runTest(
		"Create Playlist - Empty Name (Should Fail)",
		"POST",
		"/playlists",
		`{"name":"","description":"Test description"}`,
		400,
	)

	runTest(
		"Create Playlist - Empty Description (Should Fail)",
		"POST",
		"/playlists",
		`{"name":"Test Playlist","description":""}`,
		400,
	)

	runTest(
		"Create Playlist - Invalid JSON (Should Fail)",
		"POST",
		"/playlists",
		`{"name":"Test Playlist","description":"Test description"`,
		400,
	)

	// Tests de casos edge
	runTest(
		"Get Song - Invalid ID (Should Fail)",
		"GET",
		"/songs/invalid",
		"",
		400,
	)

	runTest(
		"Get Song - Non-existent ID (Should Fail)",
		"GET",
		"/songs/999",
		"",
		404,
	)

	runTest(
		"Update Song - Non-existent ID (Should Fail)",
		"PUT",
		"/songs/999",
		`{"title":"Test","artist":"Test"}`,
		404,
	)

	runTest(
		"Delete Song - Non-existent ID (Should Fail)",
		"DELETE",
		"/songs/999",
		"",
		404,
	)

	// Tests de Health Check
	runTest(
		"Health Check",
		"GET",
		"/health",
		"",
		200,
	)

	// Imprimir resultados
	printResults()

	// Guardar logs
	saveLogs()

	// Exit code basado en resultados
	if testSuite.Failed > 0 || testSuite.Errors > 0 {
		os.Exit(1)
	}
}
