# Yboost — Steam Library Explorer

Application web en Go qui affiche ta bibliothèque Steam et sauvegarde ton Top 5 des jeux les plus joués dans une base de données Supabase.

---

## Stack technique

| Couche | Technologie |
|---|---|
| Langage | Go 1.26 |
| Framework HTTP | `net/http` (stdlib) |
| ORM | GORM + driver `pgx` |
| Base de données | PostgreSQL via Supabase |
| Hébergement | Render |
| API externe | Steam Web API |
| Frontend | HTML/CSS vanilla + Go templates |

---

## Structure du projet

```
.
├── Procfile          # Commande de démarrage pour Render
├── README.md
├── api/
│   └── steam.go      # Appels à l'API Steam
├── database.go       # Modèle GORM + logique CRUD
├── db.go             # Initialisation de la connexion PostgreSQL
├── go.mod
├── go.sum
├── html/
│   ├── index.html    # Template Go + JS vanilla
│   └── style.css     # Thème Steam dark
└── main.go           # Serveur HTTP + handlers
```

---

## Variables d'environnement

À configurer dans **Render → Environment** (ou un fichier `.env` en local) :

| Variable | Description | Exemple |
|---|---|---|
| `API_KEY` | Clé Steam Web API ([obtenir ici](https://steamcommunity.com/dev/apikey)) | `XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX` |
| `STEAM_ID` | SteamID64 par défaut affiché au chargement | `76561199111246227` |
| `DATABASE_URL` | URL de connexion PostgreSQL Supabase (pooler Transaction, port 6543) | `postgresql://postgres.xxx:password@aws-0-eu-central-1.pooler.supabase.com:6543/postgres` |
| `PORT` | Port d'écoute du serveur (Render l'injecte automatiquement) | `8080` |

> ⚠️ `DATABASE_URL` doit pointer vers le **Transaction pooler** de Supabase (port `6543`), pas la connexion directe (port `5432`).

---

## Base de données Supabase

### Table `steamDB`

À créer manuellement dans **Supabase → SQL Editor** :

```sql
-- Ajout des colonnes si la table existe déjà
ALTER TABLE "steamDB" ADD COLUMN IF NOT EXISTS steam_id TEXT;
ALTER TABLE "steamDB" ADD COLUMN IF NOT EXISTS game4 TEXT;
ALTER TABLE "steamDB" ADD COLUMN IF NOT EXISTS game5 TEXT;

-- Contrainte unique nécessaire pour l'upsert
ALTER TABLE "steamDB" ADD CONSTRAINT steamdb_steam_id_unique UNIQUE (steam_id);
```

### Schéma final attendu

| Colonne | Type | Description |
|---|---|---|
| `id` | int8 | Clé primaire auto-incrémentée |
| `steam_id` | text | SteamID64 du joueur (unique) |
| `game1` | text | Jeu le plus joué |
| `game2` | text | 2ème jeu le plus joué |
| `game3` | text | 3ème jeu le plus joué |
| `game4` | text | 4ème jeu le plus joué |
| `game5` | text | 5ème jeu le plus joué |

Chaque jeu est stocké au format `Nom du jeu (Xh)`.

---

## Routes HTTP

| Méthode | Route | Description |
|---|---|---|
| `GET` | `/` | Page principale — affiche la bibliothèque Steam du joueur |
| `GET` | `/sync?iduser=<steamID>` | Sauvegarde le Top 5 dans Supabase (INSERT ou UPDATE) |
| `GET` | `/top?iduser=<steamID>` | Retourne le Top 5 stocké en base (JSON) |
| `GET` | `/delete?iduser=<steamID>` | Supprime la ligne du joueur en base |
| `GET` | `/static/*` | Fichiers statiques (CSS) |

### Paramètres

- `iduser` : SteamID64 du joueur. Si absent, utilise la variable d'environnement `STEAM_ID`.

### Exemple de réponse `/top`

```json
{
  "ID": 3,
  "SteamID": "76561199111246227",
  "Game1": "Hollow Knight (236h)",
  "Game2": "Hollow Knight: Silksong (164h)",
  "Game3": "Stardew Valley (126h)",
  "Game4": "Phasmophobia (90h)",
  "Game5": "Subnautica (85h)"
}
```

---

## API Steam utilisées

| Endpoint | Usage |
|---|---|
| `IPlayerService/GetOwnedGames/v0001` | Récupère tous les jeux possédés avec temps de jeu |
| `ISteamUserStats/GetPlayerAchievements/v0001` | Récupère les succès d'un jeu (profil public requis) |

> ⚠️ Le profil Steam du joueur doit être **Public** pour que l'API retourne des données.

---

## Déploiement sur Render

### Prérequis

- Compte [Render](https://render.com)
- Compte [Supabase](https://supabase.com)
- Clé Steam Web API

### Étapes

1. **Fork ou push** le projet sur GitHub
2. Dans Render → **New Web Service** → connecte ton repo
3. Configure le build :
   - **Build Command** : `go build -tags netgo -ldflags '-s -w' -o app`
   - **Start Command** : `./app` (ou via `Procfile`)
4. Ajoute les **variables d'environnement** (voir section ci-dessus)
5. Dans Supabase, exécute les **requêtes SQL** de création de table
6. **Deploy** 🚀

### Procfile

```
web: ./app
```

---

## Lancement en local

```bash
# Cloner le projet
git clone https://github.com/poketoto45/Yboost
cd Yboost

# Définir les variables d'environnement
export API_KEY="ta_cle_steam"
export STEAM_ID="ton_steamid64"
export DATABASE_URL="postgresql://..."

# Lancer
go run .
```

Le serveur écoute sur `http://localhost:8080`.
