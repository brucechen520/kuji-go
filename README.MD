# Description
This is my first go project, so I will add commentation to help me understand how go works.

# Flow
- `docker-compose up -d`
- `go run fixtures` -> generate fake datas
- `go run ./cmd/api`
- install wire: `go install github.com/google/wire/cmd/wire@latest`
- wire generate dependency injection file: `wire ./cmd/api`


# Test
- api: [post] `http://localhost:3000/api/v1/client/login`
	- body:
		```bash
			{
				"email": "admin@kuji.go"
			}
		```
	- resposne
		```bash
			{
				"token": "gen_uuid_here"
			}
		```
