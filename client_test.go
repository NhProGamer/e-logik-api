package elogik

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c := New("user", "secret")
	c.BaseURL = srv.URL + "/api/ecomm/v2"
	return c
}

// Every operation in openapi.json must be reachable through exactly one
// client method with the right HTTP verb and path.
func TestServicesMatchSpec(t *testing.T) {
	raw, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatal(err)
	}
	param := regexp.MustCompile(`\{[^}]+\}`)
	want := map[string]bool{}
	for p, ops := range spec.Paths {
		for m := range ops {
			want[strings.ToUpper(m)+" "+param.ReplaceAllString(p, "x")] = true
		}
	}

	got := map[string]bool{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got[r.Method+" "+r.URL.Path] = true
		if strings.HasSuffix(r.URL.Path, "/transporteurs/liste") {
			w.Write([]byte("[]"))
			return
		}
		w.Write([]byte("{}"))
	})

	ctx := context.Background()
	calls := []func() error{
		func() error { _, err := c.Articles.Liste(ctx, WebFiltresArticlesDTO{}); return err },
		func() error { _, err := c.Articles.Stock(ctx, WebFiltresArticlesDTOStock{}); return err },
		func() error { _, err := c.Articles.StockPositif(ctx, WebFiltresArticlesDTOStock{}); return err },
		func() error { _, err := c.Articles.Creer(ctx, APIFormArticleDTO{}); return err },
		func() error { _, err := c.Articles.Modifier(ctx, "x", APIFormArticleDTO{}); return err },
		func() error {
			_, err := c.Articles.AjouterFournisseur(ctx, "x", APIFormArticleFournisseurDTO{})
			return err
		},
		func() error {
			_, err := c.Articles.ModifierFournisseur(ctx, "x", APIFormArticleFournisseurDTO{})
			return err
		},
		func() error { _, err := c.Articles.RetirerFournisseur(ctx, "x", "x"); return err },
		func() error {
			_, err := c.Articles.AjouterContenuLot(ctx, "x", APIFormLotListeArticlesDTO{})
			return err
		},
		func() error { _, err := c.Colis.Liste(ctx, WebFiltresColisDTO{}); return err },
		func() error { _, err := c.Colis.DetailsDLC(ctx, APIFiltresListeColisDTO{}); return err },
		func() error { _, err := c.Commandes.Liste(ctx, WebFiltresCommandesDTO{}); return err },
		func() error { _, err := c.Commandes.Details(ctx, APIFiltresCommandesDTO2{}); return err },
		func() error { _, err := c.Commandes.Colis(ctx, APIFiltresCommandesDTO{}); return err },
		func() error {
			_, err := c.Commandes.ColisEvenements(ctx, APIFiltresCommandesDTOEvenements{})
			return err
		},
		func() error { _, err := c.Commandes.BonLivraison(ctx, "x"); return err },
		func() error { _, err := c.Commandes.Creer(ctx, APIFormCommandeDTO{}); return err },
		func() error { _, err := c.Commandes.Modifier(ctx, "x", APIFormCommandeDTO{}); return err },
		func() error {
			_, err := c.Commandes.AjouterArticle(ctx, "x", APIFormCommandeListeArticlesDTO{})
			return err
		},
		func() error { _, err := c.Commandes.ModifierArticle(ctx, "x", APIFormCommandeArticleDTO{}); return err },
		func() error { _, err := c.Commandes.Annuler(ctx, "x", APIFormCommandeAnnulationDTO{}); return err },
		func() error {
			_, err := c.Commandes.AjouterDocument(ctx, "x", APIFormCommandeDocumentDTO{})
			return err
		},
		func() error { _, err := c.Commandes.AjouterFacture(ctx, "x", APIFormCommandeDocumentDTO{}); return err },
		func() error {
			_, err := c.Commandes.AjouterFormulaire(ctx, "x", APIFormCommandeDocumentDTO{})
			return err
		},
		func() error {
			_, err := c.CommandesFournisseur.Liste(ctx, WebFiltresCommandesFournisseurDTO{})
			return err
		},
		func() error { _, err := c.CommandesFournisseur.Details(ctx, "x"); return err },
		func() error { _, err := c.CommandesFournisseur.Creer(ctx, APIFormCommandeFournisseurDTO{}); return err },
		func() error {
			_, err := c.CommandesFournisseur.Modifier(ctx, "x", APIFormCommandeFournisseurDTO{})
			return err
		},
		func() error { return c.CommandesFournisseur.Supprimer(ctx, "x") },
		func() error { _, err := c.CommandesFournisseur.AjouterArticles(ctx, "x", nil); return err },
		func() error {
			_, err := c.CommandesFournisseur.ModifierArticle(ctx, "x", APIFormCommandeFournisseurArticleDTO{})
			return err
		},
		func() error { _, err := c.Fournisseurs.Details(ctx, "x"); return err },
		func() error { _, err := c.Packs.Stock(ctx, "x"); return err },
		func() error { _, err := c.Packs.Stocks(ctx, WebFiltresPacksDTO{}); return err },
		func() error { _, err := c.Retours.Liste(ctx, WebFiltresRetourDTO{}); return err },
		func() error { _, err := c.Retours.Details(ctx, ApiFiltresRetoursDTO{}); return err },
		func() error { _, err := c.Retours.Creer(ctx, APIFormRetourDTO{}); return err },
		func() error { _, err := c.Retours.Modifier(ctx, "x", APIFormRetourDTO2{}); return err },
		func() error { return c.Retours.Supprimer(ctx, "x") },
		func() error { _, err := c.Retours.AjouterArticle(ctx, "x", APIFormRetourListeArticleDTO{}); return err },
		func() error { _, err := c.Retours.ModifierArticle(ctx, "x", APIFormRetourArticleDTO2{}); return err },
		func() error {
			_, err := c.Retours.ModifierAdresseLivraison(ctx, "x", APIFormAdresseRetourDTO{})
			return err
		},
		func() error { _, err := c.Transporteurs.Liste(ctx); return err },
		func() error { _, err := c.Webhooks.Liste(ctx); return err },
		func() error { _, err := c.Webhooks.Creer(ctx, APIFormWebhookDTO{}); return err },
		func() error { return c.Webhooks.Supprimer(ctx, "x") },
	}
	for _, call := range calls {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}

	for k := range want {
		if !got[k] {
			t.Errorf("spec operation not called by the client: %s", k)
		}
	}
	for k := range got {
		if !want[k] {
			t.Errorf("client calls an operation absent from the spec: %s", k)
		}
	}
}

func TestDo(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); !ok || u != "user" || p != "secret" {
			t.Errorf("basic auth = %q %q %v", u, p, ok)
		}
		if r.URL.EscapedPath() != "/api/ecomm/v2/commandes/a%2Fb/bon-livraison" {
			t.Errorf("path = %s", r.URL.EscapedPath())
		}
		w.Header().Set("X-OMS-Api-Call-Limit", "0/400")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"code":"TOO_MANY_REQUESTS","message":"Rate Limit Exceeded"}`))
	})

	if _, _, ok := c.CallLimit(); ok {
		t.Error("CallLimit ok before any call")
	}
	_, err := c.Commandes.BonLivraison(context.Background(), "a/b")
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 || apiErr.Code != "TOO_MANY_REQUESTS" {
		t.Fatalf("err = %#v", err)
	}
	if rem, size, ok := c.CallLimit(); !ok || rem != 0 || size != 400 {
		t.Errorf("CallLimit = %d/%d %v", rem, size, ok)
	}
}

func TestDoNonJSONError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	})
	_, err := c.Webhooks.Liste(context.Background())
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 || apiErr.Message != "boom" {
		t.Fatalf("err = %#v", err)
	}
}
