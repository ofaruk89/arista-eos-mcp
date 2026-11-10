# Arista EOS MCP Server (sahibinden)

[x86taka/arista-eos-mcp](https://github.com/x86taka/arista-eos-mcp) tabanli, **read-only** Arista EOS MCP sunucusu (eAPI + gNMI).
Orijinal dokumantasyon: [docs/UPSTREAM-README.md](docs/UPSTREAM-README.md).

## Orijinalden farklar
- **Bearer token auth:** `MCP_AUTH_TOKEN` doluysa HTTP endpoint `Authorization: Bearer <token>` ister (401 aksi halde).
- **Ortak credential:** `defaults` blogu envanterdeki cihazlara da uygulanir (orijinalde sadece ad-hoc host'lar icindi).
- **Enable'siz calisma:** `enable_password` bos ise komutlara `enable` eklenmez (enable yetkisi olmayan hesaplar icin).

## Dosyalar
| Dosya | Aciklama |
| --- | --- |
| `docker-compose.yml` | Servis tanimi (imaj Harbor'dan cekilir) |
| `.env` | `MCP_AUTH_TOKEN`, `BIND_ADDR`, `IMAGE_TAG` (**gizli, git'e girmez**) |
| `devices.json` | Switch listesi ve credential'lar (**gizli, git'e girmez**) |
| `.env.example`, `devices.example.json` | Ornek sablonlar |

## Kurulum
```sh
cp .env.example .env && chmod 600 .env                 # token'i doldur: openssl rand -hex 32
cp devices.example.json devices.json && chmod 600 devices.json   # switch + credential'lari doldur
docker compose up -d
```
`defaults` blogundaki `username`/`password` (ve istenirse `enable_password`, `eapi_port` vb.) tum cihazlara uygulanir; cihaz basina sadece `name` ve `host` yazmak yeterlidir. Farkli credential'li cihaza kendi `username`/`password` alanlarini yazin, cihazdaki deger onceliklidir.
Switch'lerde eAPI acik olmali: `management api http-commands` > `protocol https` > `no shutdown`.

## Guncelleme
Compose imaji `latest` etiketiyle ve `pull_policy: always` ile calisir; en son imaji almak icin:
```sh
docker compose up -d      # her seferinde registry'den en son imaji ceker
```
Belirli bir surume sabitlemek icin `.env` icinde `IMAGE_TAG=<surum>` yazin.

## Imaj build / Harbor'a gonderme
```sh
docker build -t harbor.sahibindenlocal.net/nwsec/arista/arista-eos-mcp:<surum> .
docker login harbor.sahibindenlocal.net
docker push harbor.sahibindenlocal.net/nwsec/arista/arista-eos-mcp:<surum>
```

## GitHub Container Registry (ghcr.io)
`.github/workflows/ghcr.yml` imaji GitHub Actions ile build edip `ghcr.io/ofaruk89/arista-eos-mcp` altina yollar
(`master` push'u -> `:latest` ve `:sha-<kisa>`, `v*` tag'i -> `:1.2.3`; Actions > ghcr > Run workflow ile elle de calistirilabilir).
Ek secret gerekmez (`GITHUB_TOKEN`). Ilk calismadan sonra paket **private** olur; gorunurlugu
GitHub > Packages > arista-eos-mcp > Package settings'ten ayarlayin.

ghcr'dan kullanmak icin `.env`:
```
IMAGE_REPO=ghcr.io/ofaruk89/arista-eos-mcp
IMAGE_TAG=latest
```
Private paket icin sunucuda bir kez giris gerekir (token'da `read:packages` yetkisi):
```sh
echo <PAT> | docker login ghcr.io -u ofaruk89 --password-stdin
docker compose pull && docker compose up -d
```

## Istemci config'i
Claude Code:
```sh
claude mcp add --transport http --scope user arista-eos http://10.211.112.12:8080/mcp \
  --header "Authorization: Bearer <TOKEN>"
```
Claude Desktop / baska istemci (mcp-remote):
```json
"arista-eos": {
  "command": "npx",
  "args": ["mcp-remote", "http://10.211.112.12:8080/mcp", "--allow-http",
           "--header", "Authorization: Bearer <TOKEN>"]
}
```
Token `.env` dosyasindaki `MCP_AUTH_TOKEN` degeridir.

## Kontrol
```sh
docker compose logs --tail 20                  # "loaded N device(s)" ve "bearer token authentication enabled"
curl -s -o /dev/null -w '%{http_code}\n' -X POST http://10.211.112.12:8080/mcp   # 401 beklenir
```

## Guvenlik notlari
- Sunucu `read_only: true` ile yazma komutlarini (`configure`, `write`, `reload`, `clear` ...) reddeder; yine de switch'teki hesabi read-only role ile sinirlayin.
- `show` disinda `ping`, `traceroute`, `dir` izinlidir; switch hesabinin rolu de bunlara izin vermelidir.
- Trafik HTTP'dir; token aginda duz metin gider. Sadece guvenilir ic agda kullanin.
- Token'i degistirmek: `.env` icindeki `MCP_AUTH_TOKEN`'i yenile, `docker compose up -d`, istemci config'lerini guncelle.
