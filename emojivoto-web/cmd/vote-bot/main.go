package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
//
// When not voting for :doughnut:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

var (
	client *http.Client

	ocagentHost = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

type emoji struct {
	Shortcode string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	webHost := os.Getenv("WEB_HOST")
	if webHost == "" {
		log.Fatalf("WEB_HOST environment variable must me set")
	}

	hostOverride := os.Getenv("HOST_OVERRIDE")

	// setting the the TTL is optional, thus invalid numbers are simply ignored
	timeToLive, _ := strconv.Atoi(os.Getenv("TTL"))
	var deadline time.Time // zero value of time.Time

	if timeToLive != 0 {
		deadline = time.Now().Add(time.Second * time.Duration(timeToLive))
	}

	// setting the the request rate is optional, thus invalid numbers are simply ignored
	requestRate, _ := strconv.Atoi(os.Getenv("REQUEST_RATE"))
	if requestRate < 1 {
		requestRate = 1
	}

	ctx := context.Background()
	ote, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithReconnectionPeriod(5*time.Second),
		otlptracegrpc.WithEndpoint(ocagentHost),
	)
	if err != nil {
		log.Fatalf("Failed to create oteltracegrpc-exporter: %v", err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("vote-bot"),
		),
	)
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(ote),
		sdktrace.WithResource(r))
	otel.SetTracerProvider(traceProvider)

	// Initialize HTTP client after trace provider is set up
	client = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	webURL := "http://" + webHost
	if _, err := url.Parse(webURL); err != nil {
		log.Fatalf("WEB_HOST %s is invalid", webHost)
	}

	for {
		// check if deadline has been reached, when TTL has been set.
		if (!deadline.IsZero()) && time.Now().After(deadline) {
			fmt.Printf("Time to live of %d seconds reached, completing\n", timeToLive)
			os.Exit(0)
		}

		time.Sleep(time.Second / time.Duration(requestRate))

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(webURL, hostOverride)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(webURL, hostOverride, ":doughnut:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(webURL, hostOverride, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func shortcodes(webURL string, hostOverride string) ([]string, error) {
	url := fmt.Sprintf("%s/api/list", webURL)
	req, _ := http.NewRequest("GET", url, nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emojis []*emoji
	err = json.Unmarshal(bytes, &emojis)
	if err != nil {
		return nil, err
	}

	shortcodes := make([]string, len(emojis))
	for i, e := range emojis {
		shortcodes[i] = e.Shortcode
	}

	return shortcodes, nil
}

func vote(webURL string, hostOverride string, shortcode string) error {
	fmt.Printf("✔ Voting for %s\n", shortcode)

	url := fmt.Sprintf("%s/api/vote?choice=%s", webURL, shortcode)
	req, _ := http.NewRequest("GET", url, nil)
	if hostOverride != "" {
		req.Host = hostOverride
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
