**Запуск миграции:**<br>
`go run cmd/migration/main.go --action <action>`<br>
up - Поднять миграции <br>
down - Откатить миграции
---
**API:** <br>
Auth Header - "Authorization: Bearer ..." <br>
`/api/auth POST body: {"login":"login", "password": "password"}`<br>
`/api/upload-asset/{assetName} GET`<br>
`/api/asset/{assetName} GET` <br>
`/api/asset/get-list POST body: {"limit": 1, "offset": 1}` <br>
`/api/asset/delete/{assetName} DELETE`