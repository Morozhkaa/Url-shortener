## Url-shortener

URL shortener in Go. You can choose one of three modes:

* __in-memory__ - all data is stored in memory as map[string]string

* __mongo__ - use MongoDB for persistence

* __cached__ - use MongoDB for persistence and Redis for caching


### API Description
Two methods available:

    1. POST, path: "/api/urls", requestBody: {"url": "<YourLink>"}
       Create a new short link for the given address.

---
    2. GET, path: "/{key}", key - received short link (string with pattern `\w{5}`)
       Executes Redirect.


### Usage

1. Set the desired MODE in environment variables in docker-compose.yml file.

2. Run the application: `docker-compose up --build`.

3. Create a short link, for example: 

        curl -XPOST http://localhost:8080/api/urls -d '{"url": "https://habr.com/ru/companies/ruvds/articles/562878/"}'


Get the response in the following json format:

    {
        "key": "sHZua"
    }

4. Use the resulting key to redirect to the desired page: `http://localhost:8080/sHZua`



### Service features

The service uses a ___rate limiter___, which allows you to limit the number of requests per unit of time.

In this way, we can protect the system from accidental or malicious excess requests, which may cause a delay or denial of service to other clients.