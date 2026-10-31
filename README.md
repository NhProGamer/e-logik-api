# e-logik-api

Client Go pour l'API e-commerce v2 de l'OMS e-Logik
([documentation officielle](https://oms.ekan-roissy.fr/api/doc/stoplight)).

## Installation

```sh
go get e-logik-api
```

Le module s'appelle `e-logik-api` et le package `elogik`. Pense à changer le
chemin du module dans `go.mod` quand le dépôt aura un remote.

## Démarrage rapide

```go
import elogik "e-logik-api"

c := elogik.New(os.Getenv("ELOGIK_USER"), os.Getenv("ELOGIK_PASSWORD"))

transporteurs, err := c.Transporteurs.Liste(ctx)
if err != nil {
    return err
}
for _, t := range transporteurs {
    fmt.Println(*t.Code, *t.Nom)
}
```

L'authentification se fait en HTTP Basic avec les identifiants du compte API.
Ne les mets pas en dur dans le code.

## Configuration

`New` renvoie un client prêt à l'emploi. Tous les réglages sont des champs
exportés, modifiables après création :

| Champ        | Défaut                                     | Rôle                                         |
| ------------ | ------------------------------------------ | -------------------------------------------- |
| `BaseURL`    | `https://oms.ekan-roissy.fr/api/ecomm/v2`  | Instance de l'API                            |
| `HTTPClient` | `http.DefaultClient`                       | Client HTTP (timeouts, proxy, transport)     |
| `MaxRetries` | `5`                                        | Nombre de retries sur 429, `0` pour désactiver |
| `RetryDelay` | `1s`                                       | Première attente après un 429, doublée ensuite |

```go
c := elogik.New(user, pass)
c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
c.MaxRetries = 2
```

Construis toujours le client avec `New` : un `elogik.Client{}` littéral n'a ni
services ni retries.

## Endpoints

Les méthodes sont rangées par ressource et suivent les noms de l'API. Toutes
prennent un `context.Context` en premier argument.

| Service                  | Méthodes                                                                                                                              |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------- |
| `c.Articles`             | `Liste`, `Stock`, `StockPositif`, `Creer`, `Modifier`, `AjouterFournisseur`, `ModifierFournisseur`, `RetirerFournisseur`, `AjouterContenuLot` |
| `c.Colis`                | `Liste`, `DetailsDLC`                                                                                                                 |
| `c.Commandes`            | `Liste`, `Details`, `Colis`, `ColisEvenements`, `BonLivraison`, `Creer`, `Modifier`, `AjouterArticle`, `ModifierArticle`, `Annuler`, `AjouterDocument`, `AjouterFacture`, `AjouterFormulaire` |
| `c.CommandesFournisseur` | `Liste`, `Details`, `Creer`, `Modifier`, `Supprimer`, `AjouterArticles`, `ModifierArticle`                                            |
| `c.Fournisseurs`         | `Details`                                                                                                                             |
| `c.Packs`                | `Stock`, `Stocks`                                                                                                                     |
| `c.Retours`              | `Liste`, `Details`, `Creer`, `Modifier`, `Supprimer`, `AjouterArticle`, `ModifierArticle`, `ModifierAdresseLivraison`                 |
| `c.Transporteurs`        | `Liste`                                                                                                                               |
| `c.Webhooks`             | `Liste`, `Creer`, `Supprimer`                                                                                                         |

Les types de requête et de réponse sont générés depuis la spec et gardent ses
noms (`APIFormCommandeDTO`, `CommandeInfoDTO5`...). La signature de chaque
méthode dans `services.go` indique lequel utiliser. Les champs optionnels sont
des pointeurs.

### Exemple : créer une commande

```go
cmd, err := c.Commandes.Creer(ctx, elogik.APIFormCommandeDTO{
    DateCommande:            time.Now(),
    CodeServiceTransporteur: 12,
    // ...
})
if err != nil {
    return err
}
// L'API recommande de vérifier que la commande n'a pas été bloquée
// par ses contrôles de conformité.
if cmd.ListeMotifBlocageLogistique != nil && len(*cmd.ListeMotifBlocageLogistique) > 0 {
    // ...
}
```

Les champs date-seule des requêtes (`dateExpeditionSouhaitee`,
`dateLivraison`, `dlc`, et les dates des commandes fournisseur) sont des
`string` au format `2006-01-02`.

## Erreurs

Toute réponse hors 2xx devient un `*elogik.Error` :

```go
var apiErr *elogik.Error
if errors.As(err, &apiErr) {
    log.Printf("HTTP %d %s: %s", apiErr.StatusCode, apiErr.Code, apiErr.Message)
}
```

| Statut | Signification                                                                              |
| ------ | ------------------------------------------------------------------------------------------ |
| 400    | Requête invalide. `Request limit exceeded` : plus de 400 références d'un coup, à découper. |
| 403    | Fonctionnalité non activée sur le compte (Webhooks, Suivi Colis...).                       |
| 429    | Réserve de jetons toujours vide après tous les retries.                                    |

Si le corps de l'erreur n'est pas du JSON, il est mis tel quel dans `Message`.

## Limitation du nombre d'appels

L'API v2 donne à chaque compte une réserve de **400 jetons**, rechargée d'**un
jeton par seconde**. Un appel coûte 1 à 3 jetons, ou 1 jeton par référence
pour les appels sur une liste (voir la documentation officielle).

Le client gère ça de deux façons :

- **Retry automatique sur 429.** L'API garantit qu'un appel refusé en 429 n'a
  pas été traité, il est donc renvoyé sans risque de doublon, y compris pour
  les créations. Avec les réglages par défaut, il attend 1, 2, 4, 8 puis 16 s
  (31 s au plus) avant d'abandonner. L'attente s'arrête dès que le `ctx` est
  annulé : mets un timeout sur le `ctx` si l'appel ne doit pas bloquer
  longtemps.
- **Suivi de la réserve.** `CallLimit` renvoie l'état lu dans le dernier
  en-tête `X-OMS-Api-Call-Limit` :

  ```go
  if restants, taille, ok := c.CallLimit(); ok && restants < 50 {
      log.Printf("réserve basse : %d/%d jetons", restants, taille)
  }
  ```

Le client ne ralentit pas les appels de lui-même et ne découpe pas les listes
de plus de 400 références.

## Développement

```sh
go test -race ./...
```

`TestServicesMatchSpec` appelle chaque méthode contre un faux serveur et
vérifie que les couples verbe + chemin correspondent exactement à
`openapi.json`. Un endpoint oublié ou une faute de frappe dans un chemin le
fait échouer.

### Mettre à jour la spec

1. Télécharger la nouvelle spec : elle est dans la balise
   `<script id="swagger-data">` de https://oms.ekan-roissy.fr/api/doc/swagger
   (champ `spec`), à enregistrer dans `openapi.json`.
2. Régénérer les types (nécessite `python3`) :

   ```sh
   go generate ./...
   ```

   `scripts/fix_spec.py` corrige d'abord les erreurs de la spec produite par
   NelmioApiDocBundle (`type: bool`, enums donnés en nom de classe PHP,
   `oneOf` à une branche) que oapi-codegen refuse.
3. Lancer les tests. Si `TestServicesMatchSpec` échoue, ajouter ou corriger la
   méthode correspondante dans `services.go`. Si un type a été renuméroté par
   Nelmio (`CommandeInfoDTO5` -> `CommandeInfoDTO6`), la compilation l'indique.

## Limites connues

- Le client n'a pas encore été testé contre la vraie API.
- Les dates des réponses sont décodées en `time.Time`. Si l'API renvoie une
  date seule (`2020-01-01`) dans un champ déclaré `date-time`, le décodage de
  la réponse échoue.
