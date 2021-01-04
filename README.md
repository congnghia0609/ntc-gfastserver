# ntc-gfastserver
**ntc-gfastserver** is an example golang http server using [fasthttp](https://github.com/valyala/fasthttp)

## Quick start
```bash
# install library dependencies
make deps

# build
make build

# start mode development
make run

# clean build
make clean
```

## Call API Post
### Add New Post
```bash
curl -X POST -i 'http://127.0.0.1:8080/post' \
  -H "Content-Type: application/json" \
  --data '{
    "title": "title1",
    "body": "body1"
  }'
```

### Update Post
```bash
curl -X PUT -i 'http://127.0.0.1:8080/post' \
  -H "Content-Type: application/json" \
  --data '{
  	"id": 1,
    "title": "title1 update",
    "body": "body1 update"
  }'
```

### Get Post
```bash
# Get a post
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/post/1'

# Get all post
curl -X GET -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/posts'
```

### Delete Post
```bash
curl -X DELETE -H 'Content-Type: application/json' \
  -i 'http://127.0.0.1:8080/post/1'
```

## License
This code is under the [Apache License v2](https://www.apache.org/licenses/LICENSE-2.0).  
