# gator

# Requirements
PostgreSQL Golang

- Set up a config file called `.gatorconfig.json` in the home directory `~/`.
- The file should have a structure as the following:
```
{"db_url":"[POSTGRES URL]","current_user_name":"DB USERNAME"}
```
- 

# Installation
Run the command:
```
go install github.com/firerockets/gator@latest
```