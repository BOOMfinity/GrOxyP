package client

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestUpdateKeepsLastValidDatabase(t *testing.T) {
	for _, body := range []string{"", "# comments only\n", "<html>Upstream failure</html>", "203.0.113.0/24\ninvalid CIDR", "203.0.113.0/24\n" + strings.Repeat("x", 65537)} {
		t.Run(fmt.Sprintf("length_%d", len(body)), func(t *testing.T) {
			var fail atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if fail.Load() {
					fmt.Fprint(w, body)
					return
				}
				fmt.Fprintln(w, "198.51.100.0/24")
			}))
			defer server.Close()
			c, err := NewClient(Config{DatabaseDownloadURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			fail.Store(true)
			if err := c.Update(); err == nil {
				t.Fatal("invalid update accepted")
			}
			if found, _ := c.FindIP(net.ParseIP("198.51.100.1")); !found {
				t.Fatal("last valid database was lost")
			}
			if found, _ := c.FindIP(net.ParseIP("203.0.113.1")); found {
				t.Fatal("partially parsed update was published")
			}
		})
	}
}

func TestUpdateReplacesDatabaseAndSupportsIPv6(t *testing.T) {
	var replace atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if replace.Load() {
			fmt.Fprint(w, "# new database\n 203.0.113.0/24 # comment\n2001:db8::/32\n\n")
			return
		}
		fmt.Fprintln(w, "198.51.100.0/24")
	}))
	defer server.Close()
	c, err := NewClient(Config{DatabaseDownloadURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	replace.Store(true)
	if err := c.Update(); err != nil {
		t.Fatal(err)
	}
	for _, item := range []struct {
		ip   string
		want bool
	}{{"198.51.100.1", false}, {"203.0.113.1", true}, {"2001:db8::1", true}, {"invalid", false}} {
		if found, _ := c.FindIP(net.ParseIP(item.ip)); found != item.want {
			t.Errorf("%s: got %v, want %v", item.ip, found, item.want)
		}
	}
}

func TestAutomaticUpdatesRetryAndClose(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			fmt.Fprintln(w, "198.51.100.0/24")
		case 2:
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			fmt.Fprintln(w, "203.0.113.0/24")
		}
	}))
	defer server.Close()
	c, err := NewClient(Config{DatabaseDownloadURL: server.URL, DatabaseUpdateInterval: 10 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	deadline := time.Now().Add(3 * time.Second)
	for {
		if found, _ := c.FindIP(net.ParseIP("203.0.113.1")); found {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("automatic updates did not recover after HTTP 503")
		}
		time.Sleep(time.Millisecond)
	}
	c.Close()
	before := calls.Load()
	time.Sleep(30 * time.Millisecond)
	if calls.Load() != before {
		t.Fatal("automatic updates continued after Close")
	}
}

func TestZeroIntervalDisablesUpdates(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		fmt.Fprintln(w, "198.51.100.0/24")
	}))
	defer server.Close()
	c, err := NewClient(Config{DatabaseDownloadURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	time.Sleep(30 * time.Millisecond)
	if calls.Load() != 1 {
		t.Fatal("zero interval should only download once")
	}
	if c.HTTPClient.Timeout <= 0 {
		t.Fatal("download timeout is not configured")
	}
}

func TestNewClientRejectsInvalidConfigurationAndDownloads(t *testing.T) {
	if c, err := NewClient(Config{DatabaseUpdateInterval: -time.Second}); err == nil || c != nil {
		t.Fatal("negative interval should fail without creating a client")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	if c, err := NewClient(Config{DatabaseDownloadURL: server.URL}); err == nil || c != nil {
		t.Fatal("failed initial download should not return a usable client")
	}
}

func TestConcurrentUpdatesAndLookups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "198.51.100.0/24\n2001:db8::/32\n")
	}))
	defer server.Close()
	c, err := NewClient(Config{DatabaseDownloadURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			if err := c.Update(); err != nil {
				t.Error(err)
			}
		}
	}()
	for worker := 0; worker < 4; worker++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				for _, ip := range []string{"198.51.100.1", "2001:db8::1"} {
					if found, _ := c.FindIP(net.ParseIP(ip)); !found {
						t.Errorf("%s disappeared during refresh", ip)
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}
