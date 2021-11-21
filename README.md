# eat_with_me_backend

ホットリロード起動
```
$ air
```

リリースコマンド
```
$ heroku login
$ heroku container:login
$ heroku container:push web
$ heroku container:release web
```

ローカルDB接続
```
$ docker exec -it <container id> /bin/bash
$ psql -h localhost -p 5432 -U admin -d mydb
```

本番DBの接続
```
$heroku pg:psql
```

ログの出力
```
$ heroku logs --tail
```