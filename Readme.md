### Messaggio test task

Rest Api microservice for messages processing

Available tunnel for service testing:  https://equipped-colt-large.ngrok-free.app

### API
   - POST "/msg"  (require JSON body : {"content": "some text"})
   - GET "/stats"
   - /swagger — detailed api description

### For local testing
```
make migrate_up // run sql migrations
make local //for run docker compose files
```