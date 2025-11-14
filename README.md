# Prerise

## NOTE : This is the backend application that supports [my personal website](https://crimson-sunrise.site). This is built specifically for me only, and it may not make sense to you.

Alright. The plan is, I will be writing something when I feel it, and I wanted to save the articles somewhere. I also wanted to integrate some third party APis, just to prove a point sometimes. I don't want to expose any of my keys, so I decided to write a small application that serves this purpose. 

This just uses Go's amazing standard library, even for APIs and authentication. The only third party dependencies are `SQLite3 driver` and `godotenv` to load the env variables. 

# How to run ?

just use the makefile

```bash
$> make run
```
