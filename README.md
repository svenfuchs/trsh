# Trsh

Experimental Travis CI Shell (tr'sh, pronounced `ˈtrɪʃ`) in Go leveraging API
v3 for demo and educational purposes. Also uses [led-go](https://github.com/svenfuchs/led-go)
as a readline replacement.

![trsh](https://user-images.githubusercontent.com/2208/41823678-84553fb2-7804-11e8-9bd9-45f2105f18a7.gif)

### Build

```
$ go build -o trsh
```

### Interactive shell

```
$ ./trsh
trsh ~ repository find slug=svenfuchs/trsh
GET /repo/svenfuchs%2Ftrsh ...
{
  "id": 18503283,
  "name": "trsh",
  "slug": "svenfuchs/trsh",
  "description": "Playing with a Travis CI API v3 client in Go",
  // ...
}

trsh ~ user find id=8 | .login
"svenfuchs"
```

### Arguments

```
$ ./trsh user find id=8
GET /user/8 ...
{
  "id": 8,
  "login": "svenfuchs",
  "name": "Sven Fuchs",
  // ...
}
```

### Stdin

```
$ echo -e "user find id=8\nrepository find slug=svenfuchs/led-go" | ./trsh
GET /user/8 ...
{
  "id": 8,
  "login": "svenfuchs",
  // ...
}
GET /repo/svenfuchs%2Fled-go ...
{
  "id": 8,
  "name": "led-go",
  // ...
}
```
