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
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
//
// When not voting for :doughnut:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

var (
	client = &http.Client{Transport: otelhttp.NewTransport(nil)}

	ocagentHost = os.Getenv("OC_AGENT_HOST")
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
	var deadline time.Time = time.Unix(0, 0)

	if timeToLive != 0 {
		deadline = time.Now().Add(time.Second * time.Duration(timeToLive))
	}

	// setting the the request rate is optional, thus invalid numbers are simply ignored
	requestRate, _ := strconv.Atoi(os.Getenv("REQUEST_RATE"))
	if requestRate < 1 {
		requestRate = 1
	}

	// Identify service
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("bot-vote"),
	)

	// Create exporter client
	client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(ocagentHost), otlptracegrpc.WithInsecure())
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("creating OTLP trace exporter: %v", err))
	}

	// Create Tracer Provider
	// TracerProvider will provide instrumentations with an impl of Tracer.
	// Tracer funnels data to export pipelines (span processors)
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(r),
		trace.WithSampler(trace.AlwaysSample()),
	)
	log.Print("Always sampling")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	webURL := "http://" + webHost
	if _, err := url.Parse(webURL); err != nil {
		log.Fatalf("WEB_HOST %s is invalid", webHost)
	}

	ctx := context.Background()
	for {
		//traceCtx, span := otel.Tracer("bot-vote").Start(ctx, "Main routine")
		// check if deadline has been reached, when TTL has been set.
		if (!deadline.IsZero()) && time.Now().After(deadline) {
			fmt.Printf("Time to live of %d seconds reached, completing\n", timeToLive)
			os.Exit(0)
		}

		time.Sleep(time.Second / time.Duration(requestRate))

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(ctx, webURL, hostOverride)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(ctx, webURL, hostOverride, ":doughnut:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(ctx, webURL, hostOverride, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		//defer span.End()
	}
}

func shortcodes(_ctx context.Context, webURL string, hostOverride string) ([]string, error) {
	//ctx, span := otel.Tracer("bot-vote").Start(ctx, "Shortcodes")
	//defer span.End()
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

func vote(ctx context.Context, webURL string, hostOverride string, shortcode string) error {
	//ctx, span := otel.Tracer("bot-vote").Start(ctx, "Vote")
	//defer span.End()
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
