package client

// ShowClient makes HTTP calls to Rajesh's Java show service.
//
// Java Feign equivalent:
//   @FeignClient(name = "show-service", url = "${show.service.url}")
//   public interface ShowClient {
//       @GetMapping("/internal/shows/{showId}/booked-seats")
//       List<String> getBookedSeats(@PathVariable String showId);
//   }

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ShowClient is the HTTP client for Rajesh's Java show service.
type ShowClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewShowClient constructs the client with a 5-second timeout.
// baseURL example: "http://192.168.1.7:8082"
//
// Java equivalent:
//   @Bean
//   public RestTemplate showServiceRestTemplate() {
//       RestTemplate rt = new RestTemplate();
//       rt.setRequestFactory(new HttpComponentsClientHttpRequestFactory());
//       return rt;
//   }
func NewShowClient(baseURL string) *ShowClient {
	return &ShowClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			// If the Java service doesn't respond in 5s, the call is cancelled.
			// Java equivalent: restTemplate.setConnectTimeout(5000)
			Timeout: 2 * time.Second,
		},
	}
}

// GetBookedSeats calls GET /internal/shows/{showId}/booked-seats on the Java service.
// Returns a slice of seat ID strings e.g. ["A1", "B3", "C10"].
//
// Java RestTemplate equivalent:
//   List<String> seats = restTemplate.exchange(
//       baseUrl + "/internal/shows/" + showId + "/booked-seats",
//       HttpMethod.GET, null,
//       new ParameterizedTypeReference<List<String>>() {}
//   ).getBody();
func (c *ShowClient) GetBookedSeats(ctx context.Context, showID string) ([]string, error) {
	// Step 1: Build the full URL.
	// fmt.Sprintf = String.format() in Java.
	url := fmt.Sprintf("%s/internal/shows/%s/booked-seats", c.baseURL, showID)

	log.Printf("[ShowClient] Calling: GET %s", url)

	// Step 2: Create request with the caller's context.
	// If the incoming Gin request times out or is cancelled, this call cancels too.
	// Java WebClient equivalent: webClient.get().uri(url).retrieve()
	//                                     .bodyToMono(...).timeout(Duration.ofSeconds(5))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Step 3: Execute the HTTP call.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[ShowClient] ERROR - network failure calling %s: %v", url, err)
		return nil, fmt.Errorf("show-service unreachable at %s: %w", url, err)
	}
	// Step 4: Always close body to avoid connection leaks.
	// Like try-with-resources on HttpURLConnection in Java.
	defer resp.Body.Close()

	// Step 5: Check HTTP status.
	// Go never throws on non-2xx — must check manually.
	log.Printf("[ShowClient] Response status: %d from GET %s", resp.StatusCode, url)

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[ShowClient] ERROR - show %s not found (404)", showID)
		return nil, fmt.Errorf("show %s not found in show-service", showID)
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("[ShowClient] ERROR - unexpected status %d from show-service", resp.StatusCode)
		return nil, fmt.Errorf("show-service returned unexpected status %d", resp.StatusCode)
	}

	// Step 6: Decode JSON array into []string.
	// Java: objectMapper.readValue(body, new TypeReference<List<String>>() {})
	var bookedSeats []string
	if err := json.NewDecoder(resp.Body).Decode(&bookedSeats); err != nil {
		log.Printf("[ShowClient] ERROR - failed to decode response for show %s: %v", showID, err)
		return nil, fmt.Errorf("failed to decode booked-seats response: %w", err)
	}

	log.Printf("[ShowClient] SUCCESS - show %s returned %d booked seats: %v", showID, len(bookedSeats), bookedSeats)
	return bookedSeats, nil
}
