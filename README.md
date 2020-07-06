## What is it? 
This project allows you to check if domain name of your website
is on first page of Google Search by searching given queries. Service checks each task every 5 minutes. 
If domain is not represented on the first page, you will receive telegram notification.

**Notice:** service uses API of [SerpStack](https://serpstack.com) 
which requires subscription. 
If it does not fit your needs, feel free to implement ```Searcher``` interface in your own way.

## How to start?
- Pass telegram bot token to ```notifier/config/app.yaml```
- Pass SerpStack API key to ```google-domain-checker/config/app.yaml```
- Build: ```make```
- Add user to your telegram bot by launching ```/start``` command

## How to use?
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

#### TODO:
- Add logging and error handling to goroutines
- Add notifications to telegram if something's wrong with external API
- Add tests to checker, repositories and notificator
- Add opportunity to leave chat with bot