動作確認コマンド

全件取得
```
curl -X GET http://localhost:8080/todos
```

個別取得
```
curl -X GET http://localhost:8080/todo/{id}
```

作成
```
curl -X POST http://localhost:8080/todo -H 'Content-Type: application/json' -d '{"title":"タイトル"}'
```

更新
```
curl -X PUT http://localhost:8080/todo/{id} -H 'Content-Type: application/json' -d '{"title":"タイトル"}'
```

削除
```
curl -X DELETE http://localhost:8080/todo/{id}
```