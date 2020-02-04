# backend-coding-challenge

## About the challenge

- Develop a REST microservice that lists languages used by the 100 trending public repos on GitHub.
- For every language, you need to calculate the attributes below 👇:
    - Number of repositories using this language
    - The list of repositories using the language
    - Framework popularity over the 100 repositories

## Lunch the app
1- Docker is required*
2- git clone the current repository 
3- build docker image of API  ``` docker build -t trendingapp . ```
4- run docker container  ``` docker run -p 8000:8100 -d trendingapp ```
5- now you can access the API in http://localhost:8000

## Entry Points
| Entry point | description|
|-------------|-----------|
| /languages  | get list of 100 trending repositories depending on number of stars and forks with an optional query string since={'daily' or 'weekly' or 'monthly'} monthly is set by default|
| /frameworks |get popular frameworks of the trending repositories -> by searching in each repository package file for frameworks of the language used by the repo|