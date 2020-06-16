- Pass telegram bot token to notifier/config/app.yaml
- Pass SerpStack API key to google-domain-checker/config/app.yaml
- Build: ```make```

Add the task: ```POST: localhost:8084/tasks``` with ```application/json``` body: 
```
{ 
    "domain": "example.com", 
    "country": "Country name in english", 
    "requests": [
        { 
            "text": "your text" 
        }
    ]
} 
```

You can add user to your telegram bot by launching ```/start``` command.

TODO:
- Add logging and error handling to goroutines
- Add notifications to telegram if something's wrong with external API
- Add some pause between database queries, if no results returned
- Add opportunity to leave chat with bot