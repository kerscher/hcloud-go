package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestFloatingIPClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPGetResponse{
			FloatingIP: schema.FloatingIP{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if floatingIP == nil {
		t.Fatal("no Floating IP")
	}
	if floatingIP.ID != 1 {
		t.Errorf("unexpected ID: %v", floatingIP.ID)
	}
}

func TestFloatingIPClientGetNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: ErrorCodeNotFound,
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if floatingIP != nil {
		t.Fatal("expected no Floating IP")
	}
}

func TestFloatingIPClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.FloatingIPListResponse{
			FloatingIPs: []schema.FloatingIP{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := FloatingIPListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	floatingIPs, _, err := env.Client.FloatingIP.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(floatingIPs) != 2 {
		t.Fatal("expected 2 Floating IPs")
	}
}

func TestFloatingIPClientCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		json.NewEncoder(w).Encode(schema.FloatingIPCreateResponse{
			FloatingIP: schema.FloatingIP{
				ID: 1,
			},
			Action: &schema.Action{
				ID: 1,
			},
		})
	})

	opts := FloatingIPCreateOpts{
		Type:         FloatingIPTypeIPv4,
		Description:  String("test"),
		HomeLocation: &Location{Name: "test"},
		Server:       &Server{ID: 1},
	}

	ctx := context.Background()
	result, _, err := env.Client.FloatingIP.Create(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.FloatingIP.ID != 1 {
		t.Errorf("unexpected Floating IP ID: %d", result.FloatingIP.ID)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
}

func TestFloatingIPClientAssign(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1/actions/assign", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.FloatingIPActionAssignRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Server != 1 {
			t.Errorf("unexpected server ID: %d", reqBody.Server)
		}
		json.NewEncoder(w).Encode(schema.FloatingIPActionAssignResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
		server     = &Server{ID: 1}
	)
	action, _, err := env.Client.FloatingIP.Assign(ctx, floatingIP, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestFloatingIPClientUnassign(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1/actions/unassign", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPActionAssignResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx        = context.Background()
		floatingIP = &FloatingIP{ID: 1}
	)
	action, _, err := env.Client.FloatingIP.Unassign(ctx, floatingIP)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}
