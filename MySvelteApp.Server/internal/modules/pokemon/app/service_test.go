package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode:    status,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        http.Header{"Content-Type": []string{"application/json"}},
	}
}

func TestGetRandomPokemonSuccess(t *testing.T) {
	t.Parallel()

	expectedID := rand.New(rand.NewSource(1)).Intn(3) + 1
	var countCalls int32

	transport := roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/species":
			atomic.AddInt32(&countCalls, 1)
			return jsonResponse(http.StatusOK, `{"count":3}`), nil
		case fmt.Sprintf("/pokemon/%d", expectedID):
			return jsonResponse(http.StatusOK, `{"name":"pikachu","types":[{"type":{"name":"electric"}}],"sprites":{"front_default":"https://img/pikachu.png"}}`), nil
		default:
			return nil, fmt.Errorf("unexpected request path %s", req.URL.Path)
		}
	})

	service := NewService(&http.Client{Transport: transport})
	service.random = rand.New(rand.NewSource(1))
	service.pokemonBaseURL = "https://test/pokemon/"
	service.countURL = "https://test/species"
	service.now = func() time.Time { return time.Unix(0, 0) }

	resp, err := service.GetRandomPokemon(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response")
	}
	if resp.Name == nil || *resp.Name != "pikachu" {
		t.Fatalf("expected pikachu, got %+v", resp)
	}
	if resp.Type == nil || *resp.Type != "electric" {
		t.Fatalf("expected electric type, got %+v", resp)
	}
	if resp.Image == nil || *resp.Image != "https://img/pikachu.png" {
		t.Fatalf("expected sprite url, got %+v", resp)
	}
	if got := atomic.LoadInt32(&countCalls); got != 1 {
		t.Fatalf("expected count endpoint to be called once, got %d", got)
	}
}

func TestGetRandomPokemonCachesSpeciesCount(t *testing.T) {
	t.Parallel()

	rnd := rand.New(rand.NewSource(7))
	firstID := rnd.Intn(5) + 1
	secondID := rnd.Intn(5) + 1

	var countCalls int32
	transport := roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/species":
			atomic.AddInt32(&countCalls, 1)
			return jsonResponse(http.StatusOK, `{"count":5}`), nil
		case fmt.Sprintf("/pokemon/%d", firstID):
			return jsonResponse(http.StatusOK, fmt.Sprintf(`{"name":"poke-%d","types":[],"sprites":{"front_default":null}}`, firstID)), nil
		case fmt.Sprintf("/pokemon/%d", secondID):
			return jsonResponse(http.StatusOK, fmt.Sprintf(`{"name":"poke-%d","types":[],"sprites":{"front_default":null}}`, secondID)), nil
		default:
			return nil, fmt.Errorf("unexpected request path %s", req.URL.Path)
		}
	})

	service := NewService(&http.Client{Transport: transport})
	service.random = rand.New(rand.NewSource(7))
	service.pokemonBaseURL = "https://test/pokemon/"
	service.countURL = "https://test/species"
	service.now = func() time.Time { return time.Unix(0, 0) }

	for i := 0; i < 2; i++ {
		if _, err := service.GetRandomPokemon(context.Background()); err != nil {
			t.Fatalf("call %d failed: %v", i+1, err)
		}
	}

	if got := atomic.LoadInt32(&countCalls); got != 1 {
		t.Fatalf("expected count endpoint 1 call, got %d", got)
	}
}

func TestGetRandomPokemonZeroCount(t *testing.T) {
	t.Parallel()

	transport := roundTripper(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/species":
			return jsonResponse(http.StatusOK, `{"count":0}`), nil
		default:
			return nil, fmt.Errorf("unexpected request path %s", req.URL.Path)
		}
	})

	service := NewService(&http.Client{Transport: transport})
	service.random = rand.New(rand.NewSource(1))
	service.pokemonBaseURL = "https://test/pokemon/"
	service.countURL = "https://test/species"
	service.now = func() time.Time { return time.Unix(0, 0) }

	_, err := service.GetRandomPokemon(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}

	var apiErr ExternalAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected ExternalAPIError, got %v", err)
	}
}
