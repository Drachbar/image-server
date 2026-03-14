# Image Server

En enkel HTTP-bildserver skriven i Go för att hantera uppladdning, hämtning och borttagning av bildfiler. Stöder flera applikationer med API-nyckelautentisering.

## Funktioner

- Uppladdning av bilder via multipart/form-data
- Borttagning av bilder
- API-nyckelautentisering per applikation
- Innehållsbaserad lagring med SHA1-hashning (inga duplicerade filer)
- CORS-stöd
- Healthcheck-endpoint

## Lokal utveckling

Starta servern direkt utan att bygga:

```bash
go run . -apikeyfile=./apikeys.json -dir=./images -baseurl=http://localhost:8080/images
```

Öppna sedan `test.html` i webbläsaren för att testa uppladdning mot den lokala servern.

## Bygga

```bash
mkdir -p build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/image-server
```

## Köra

```bash
./build/image-server \
  -apikeyfile=/etc/imageserver/apikeys.json \
  -dir=/var/www/images \
  -baseurl=https://cdn.example.com/images \
  -port=8080
```

### Flaggor

| Flagga | Standard | Beskrivning |
|--------|----------|-------------|
| `-apikeyfile` | `/etc/imageserver/apikeys.json` | Sökväg till JSON-fil med API-nycklar |
| `-dir` | `/var/www/images` | Katalog för uppladdade bilder |
| `-port` | `8080` | Port att lyssna på |
| `-baseurl` | `http://localhost:8080/images` | Publik bas-URL för bildlänkar |

## Lägga till en ny applikation

Kör `add-apikey.sh` på servern med applikationsnamnet som argument:

```bash
sudo ./add-apikey.sh min-app
```

Output:
```
App:        min-app
API-nyckel: a3f9d2e1b8c74f0e5a6d9b2c8e4f7a1d3b6c9e2f5a8d1e4b7c0f3a6d9b2c5e8
Sparad i:   /etc/imageserver/apikeys.json
```

Kopiera API-nyckeln och konfigurera din applikation att skicka den som `X-API-Key`-header. Scriptet startar automatiskt om tjänsten på servern så att nyckeln läses in direkt.

Vill du spara nycklar i en annan fil kan du ange sökvägen som andra argument:

```bash
./add-apikey.sh min-app ./apikeys.json
```

### API-nyckelfil

Filen är en JSON-mappning från API-nyckel till applikationsnamn. Applikationsnamnet används som undermapp i filsystemet och syns i bild-URL:er.

```json
{
  "a3f9d2e1b8c74f0e...": "min-app",
  "f7c0e3b2a9d4f1e8...": "annan-app"
}
```

## API

### POST /upload

Laddar upp en bild. Kräver `X-API-Key`-header.

**Request:** Multipart form med fältet `file` (max 10 MB)

**Response:**
```json
{"url": "https://cdn.example.com/images/app1/ab/cd/abcd1234567890.png"}
```

### DELETE /delete

Tar bort en bild. Kräver `X-API-Key`-header och query-parameter `filename`.

```
DELETE /delete?filename=abcd1234567890.png
```

### GET /healthcheck

Kontrollerar serverstatus (ingen autentisering krävs).

```json
{"status": "ok"}
```

## Filstruktur i lagring

Filer organiseras efter SHA1-hash och applikation:

```
{dir}/{app}/{hash[0:2]}/{hash[2:4]}/{hash}.{extension}
```

Exempel: `/var/www/images/app1/ab/cd/abcd1234567890abcd1234.png`

## Driftsättning på server (Alpine Linux)

### Första gången

Installera Go och Git:

```bash
echo "http://dl-cdn.alpinelinux.org/alpine/v3.21/community" >> /etc/apk/repositories
apk update && apk add git

# Installera senaste Go manuellt (kolla aktuell version på go.dev/dl)
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
export PATH=$PATH:/usr/local/go/bin
rm go1.26.1.linux-amd64.tar.gz
```

Klona repot:

```bash
mkdir -p /root/projects
cd /root/projects
git clone https://github.com/Drachbar/image-server.git
```

Lägg deploy-scriptet på plats:

```bash
cp /root/projects/image-server/deploy.sh /root/deploy.sh
chmod +x /root/deploy.sh
```

### Deploya ny version

```bash
/root/deploy.sh
```

Scriptet hämtar senaste koden från GitHub, bygger en ny binär och startar om tjänsten automatiskt.

## Krav

- Go 1.24+
- Inga externa beroenden
