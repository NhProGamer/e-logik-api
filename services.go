package elogik

import (
	"context"
	"net/http"
	"net/url"
)

var esc = url.PathEscape

// ArticlesService covers /articles.
type ArticlesService struct{ c *Client }

func (s *ArticlesService) Liste(ctx context.Context, f WebFiltresArticlesDTO) (*ArticleListDTO, error) {
	return call[ArticleListDTO](ctx, s.c, http.MethodPost, "/articles/liste", f)
}

func (s *ArticlesService) Stock(ctx context.Context, f WebFiltresArticlesDTOStock) (*ArticleListDTOStock, error) {
	return call[ArticleListDTOStock](ctx, s.c, http.MethodPost, "/articles/stock", f)
}

func (s *ArticlesService) StockPositif(ctx context.Context, f WebFiltresArticlesDTOStock) (*ArticleListDTOStock, error) {
	return call[ArticleListDTOStock](ctx, s.c, http.MethodPost, "/articles/stockPositif", f)
}

func (s *ArticlesService) Creer(ctx context.Context, a APIFormArticleDTO) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodPost, "/articles/creer", a)
}

func (s *ArticlesService) Modifier(ctx context.Context, refEcommercant string, a APIFormArticleDTO) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodPatch, "/articles/"+esc(refEcommercant)+"/modifier", a)
}

func (s *ArticlesService) AjouterFournisseur(ctx context.Context, refEcommercant string, f APIFormArticleFournisseurDTO) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodPost, "/articles/"+esc(refEcommercant)+"/ajouter-fournisseur", f)
}

func (s *ArticlesService) ModifierFournisseur(ctx context.Context, refEcommercant string, f APIFormArticleFournisseurDTO) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodPut, "/articles/"+esc(refEcommercant)+"/modifier-fournisseur", f)
}

func (s *ArticlesService) RetirerFournisseur(ctx context.Context, refEcommercant, codeFournisseur string) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodDelete, "/articles/"+esc(refEcommercant)+"/retirer-fournisseur/"+esc(codeFournisseur), nil)
}

func (s *ArticlesService) AjouterContenuLot(ctx context.Context, refEcommercant string, l APIFormLotListeArticlesDTO) (*ArticleInfoDTO5, error) {
	return call[ArticleInfoDTO5](ctx, s.c, http.MethodPost, "/articles/"+esc(refEcommercant)+"/ajouter-contenu-lot", l)
}

// ColisService covers /colis.
type ColisService struct{ c *Client }

func (s *ColisService) Liste(ctx context.Context, f WebFiltresColisDTO) (*ColisListDTOSuivi, error) {
	return call[ColisListDTOSuivi](ctx, s.c, http.MethodPost, "/colis/liste", f)
}

func (s *ColisService) DetailsDLC(ctx context.Context, f APIFiltresListeColisDTO) (*ColisListDTO, error) {
	return call[ColisListDTO](ctx, s.c, http.MethodPost, "/colis/details/dlc", f)
}

// CommandesService covers /commandes.
type CommandesService struct{ c *Client }

func (s *CommandesService) Liste(ctx context.Context, f WebFiltresCommandesDTO) (*CommandeInfoListDTO2, error) {
	return call[CommandeInfoListDTO2](ctx, s.c, http.MethodPost, "/commandes/liste", f)
}

func (s *CommandesService) Details(ctx context.Context, f APIFiltresCommandesDTO2) (*CommandeInfoListDTO, error) {
	return call[CommandeInfoListDTO](ctx, s.c, http.MethodPost, "/commandes/details", f)
}

func (s *CommandesService) Colis(ctx context.Context, f APIFiltresCommandesDTO) (*APICommandeListDTO2, error) {
	return call[APICommandeListDTO2](ctx, s.c, http.MethodPost, "/commandes/getColis", f)
}

func (s *CommandesService) ColisEvenements(ctx context.Context, f APIFiltresCommandesDTOEvenements) (*APICommandeListDTO, error) {
	return call[APICommandeListDTO](ctx, s.c, http.MethodPost, "/commandes/getColisEvenements", f)
}

func (s *CommandesService) BonLivraison(ctx context.Context, referenceEKAN string) (*APIDocumentDTO, error) {
	return call[APIDocumentDTO](ctx, s.c, http.MethodGet, "/commandes/"+esc(referenceEKAN)+"/bon-livraison", nil)
}

func (s *CommandesService) Creer(ctx context.Context, cmd APIFormCommandeDTO) (*CommandeInfoDTO5, error) {
	return call[CommandeInfoDTO5](ctx, s.c, http.MethodPost, "/commandes/creer", cmd)
}

func (s *CommandesService) Modifier(ctx context.Context, referenceEKAN string, cmd APIFormCommandeDTO) (*CommandeInfoDTO6, error) {
	return call[CommandeInfoDTO6](ctx, s.c, http.MethodPatch, "/commandes/"+esc(referenceEKAN)+"/modifier", cmd)
}

func (s *CommandesService) AjouterArticle(ctx context.Context, referenceEKAN string, a APIFormCommandeListeArticlesDTO) (*CommandeInfoDTO7, error) {
	return call[CommandeInfoDTO7](ctx, s.c, http.MethodPost, "/commandes/"+esc(referenceEKAN)+"/ajouter-article", a)
}

func (s *CommandesService) ModifierArticle(ctx context.Context, referenceEKAN string, a APIFormCommandeArticleDTO) (*CommandeInfoDTO7, error) {
	return call[CommandeInfoDTO7](ctx, s.c, http.MethodPut, "/commandes/"+esc(referenceEKAN)+"/modifier-article", a)
}

func (s *CommandesService) Annuler(ctx context.Context, referenceEKAN string, a APIFormCommandeAnnulationDTO) (*CommandeInfoDTO8, error) {
	return call[CommandeInfoDTO8](ctx, s.c, http.MethodPost, "/commandes/"+esc(referenceEKAN)+"/annuler", a)
}

func (s *CommandesService) AjouterDocument(ctx context.Context, referenceEKAN string, d APIFormCommandeDocumentDTO) (*CommandeInfoDTO9, error) {
	return call[CommandeInfoDTO9](ctx, s.c, http.MethodPost, "/commandes/"+esc(referenceEKAN)+"/document", d)
}

func (s *CommandesService) AjouterFacture(ctx context.Context, referenceEKAN string, d APIFormCommandeDocumentDTO) (*CommandeInfoDTO9, error) {
	return call[CommandeInfoDTO9](ctx, s.c, http.MethodPost, "/commandes/"+esc(referenceEKAN)+"/facture", d)
}

func (s *CommandesService) AjouterFormulaire(ctx context.Context, referenceEKAN string, d APIFormCommandeDocumentDTO) (*CommandeInfoDTO9, error) {
	return call[CommandeInfoDTO9](ctx, s.c, http.MethodPost, "/commandes/"+esc(referenceEKAN)+"/formulaire", d)
}

// CommandesFournisseurService covers /commandes-fournisseur.
type CommandesFournisseurService struct{ c *Client }

func (s *CommandesFournisseurService) Liste(ctx context.Context, f WebFiltresCommandesFournisseurDTO) (*CommandeFournisseurListDTO, error) {
	return call[CommandeFournisseurListDTO](ctx, s.c, http.MethodPost, "/commandes-fournisseur/liste", f)
}

func (s *CommandesFournisseurService) Details(ctx context.Context, numeroCommande string) (*CommandeFournisseurInfoDTO2, error) {
	return call[CommandeFournisseurInfoDTO2](ctx, s.c, http.MethodGet, "/commandes-fournisseur/"+esc(numeroCommande)+"/details", nil)
}

func (s *CommandesFournisseurService) Creer(ctx context.Context, cmd APIFormCommandeFournisseurDTO) (*CommandeFournisseurInfoDTO2, error) {
	return call[CommandeFournisseurInfoDTO2](ctx, s.c, http.MethodPost, "/commandes-fournisseur/creer", cmd)
}

func (s *CommandesFournisseurService) Modifier(ctx context.Context, numeroCommande string, cmd APIFormCommandeFournisseurDTO) (*CommandeFournisseurInfoDTO2, error) {
	return call[CommandeFournisseurInfoDTO2](ctx, s.c, http.MethodPatch, "/commandes-fournisseur/"+esc(numeroCommande)+"/modifier", cmd)
}

func (s *CommandesFournisseurService) Supprimer(ctx context.Context, numeroCommande string) error {
	return s.c.do(ctx, http.MethodDelete, "/commandes-fournisseur/"+esc(numeroCommande)+"/supprimer", nil, nil)
}

func (s *CommandesFournisseurService) AjouterArticles(ctx context.Context, numeroCommande string, a []APIFormCommandeFournisseurArticleDTO) (*CommandeFournisseurInfoDTO2, error) {
	return call[CommandeFournisseurInfoDTO2](ctx, s.c, http.MethodPost, "/commandes-fournisseur/"+esc(numeroCommande)+"/ajouter-article", a)
}

func (s *CommandesFournisseurService) ModifierArticle(ctx context.Context, numeroCommande string, a APIFormCommandeFournisseurArticleDTO) (*CommandeFournisseurInfoDTO2, error) {
	return call[CommandeFournisseurInfoDTO2](ctx, s.c, http.MethodPatch, "/commandes-fournisseur/"+esc(numeroCommande)+"/modifier-article", a)
}

// FournisseursService covers /tiers/fournisseurs.
type FournisseursService struct{ c *Client }

func (s *FournisseursService) Details(ctx context.Context, codeFournisseur string) (*FournisseurInfoDTO2, error) {
	return call[FournisseurInfoDTO2](ctx, s.c, http.MethodGet, "/tiers/fournisseurs/"+esc(codeFournisseur)+"/details", nil)
}

// PacksService covers /packs.
type PacksService struct{ c *Client }

func (s *PacksService) Stock(ctx context.Context, reference string) (*PackInfoDTO4, error) {
	return call[PackInfoDTO4](ctx, s.c, http.MethodGet, "/packs/"+esc(reference)+"/stocks", nil)
}

func (s *PacksService) Stocks(ctx context.Context, f WebFiltresPacksDTO) (*PackListDTO, error) {
	return call[PackListDTO](ctx, s.c, http.MethodPost, "/packs/stocks", f)
}

// RetoursService covers /retours.
type RetoursService struct{ c *Client }

func (s *RetoursService) Liste(ctx context.Context, f WebFiltresRetourDTO) (*RetourInfoListDTO, error) {
	return call[RetourInfoListDTO](ctx, s.c, http.MethodPost, "/retours/liste", f)
}

func (s *RetoursService) Details(ctx context.Context, f ApiFiltresRetoursDTO) (*RetourInfoListDTO2, error) {
	return call[RetourInfoListDTO2](ctx, s.c, http.MethodPost, "/retours/details", f)
}

func (s *RetoursService) Creer(ctx context.Context, r APIFormRetourDTO) (*RetourInfoDTO, error) {
	return call[RetourInfoDTO](ctx, s.c, http.MethodPost, "/retours/creer", r)
}

func (s *RetoursService) Modifier(ctx context.Context, reference string, r APIFormRetourDTO2) (*RetourInfoDTO, error) {
	return call[RetourInfoDTO](ctx, s.c, http.MethodPatch, "/retours/"+esc(reference)+"/modifier", r)
}

func (s *RetoursService) Supprimer(ctx context.Context, reference string) error {
	return s.c.do(ctx, http.MethodDelete, "/retours/"+esc(reference)+"/supprimer", nil, nil)
}

func (s *RetoursService) AjouterArticle(ctx context.Context, reference string, a APIFormRetourListeArticleDTO) (*RetourInfoDTO, error) {
	return call[RetourInfoDTO](ctx, s.c, http.MethodPost, "/retours/"+esc(reference)+"/ajouter-article", a)
}

func (s *RetoursService) ModifierArticle(ctx context.Context, reference string, a APIFormRetourArticleDTO2) (*RetourInfoDTO, error) {
	return call[RetourInfoDTO](ctx, s.c, http.MethodPatch, "/retours/"+esc(reference)+"/modifier-article", a)
}

func (s *RetoursService) ModifierAdresseLivraison(ctx context.Context, reference string, a APIFormAdresseRetourDTO) (*RetourInfoDTO, error) {
	return call[RetourInfoDTO](ctx, s.c, http.MethodPatch, "/retours/"+esc(reference)+"/modifier-adresse-livraison", a)
}

// TransporteursService covers /transporteurs.
type TransporteursService struct{ c *Client }

func (s *TransporteursService) Liste(ctx context.Context) ([]TransporteurInfoDTO2, error) {
	var out []TransporteurInfoDTO2
	return out, s.c.do(ctx, http.MethodGet, "/transporteurs/liste", nil, &out)
}

// WebhooksService covers /webhooks.
type WebhooksService struct{ c *Client }

func (s *WebhooksService) Liste(ctx context.Context) (*WebhookListDTO, error) {
	return call[WebhookListDTO](ctx, s.c, http.MethodGet, "/webhooks/liste", nil)
}

func (s *WebhooksService) Creer(ctx context.Context, w APIFormWebhookDTO) (*WebhookInfoDTO, error) {
	return call[WebhookInfoDTO](ctx, s.c, http.MethodPost, "/webhooks/creer", w)
}

func (s *WebhooksService) Supprimer(ctx context.Context, id string) error {
	return s.c.do(ctx, http.MethodDelete, "/webhooks/"+esc(id)+"/supprimer", nil, nil)
}
