## 测试新增用户

```shell
curl -X POST 'http://localhost:8888/api/v1/auth/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

返回内容：
{"code":200,"message":"注册成功","data":null}

## 用户登录

```bash
curl -X POST 'http://localhost:8888/api/v1/auth/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

返回内容：
{"code":200,"message":"登录成功","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI"}}

## 分页查询所有用户

```bash
curl -X GET 'http://localhost:8888/api/v1/user/page?page=1&pageSize=10' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI'
```

## 创建文章

```shell  
curl -X POST 'http://localhost:8888/api/v1/post/create' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI' \
  -d '{
    "title": "我的第二篇文章",
    "content": "这是文章内容..."
  }'
```

返回内容：
{"code":200,"message":"创建文章成功","data":null}

## 分页查询文章

```shell  
curl -X GET 'http://localhost:8888/api/v1/post/page?page=1&pageSize=10' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI" \
  -d '{
    "userID": 2
  }'
```

## 获取单一文章

```shell  
curl -X GET 'http://localhost:8888/api/v1/post/byId?postId=1' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI"
```

## 更新单元文章

```shell  
curl -X POST 'http://localhost:8888/api/v1/post/edit' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI' \
  -d '{
    "id": 1,
    "title": "我的第一篇文章",
    "content": "这是文章内容...修改"
  }'
```

## 删除单一文章

```shell  
curl -X GET 'http://localhost:8888/api/v1/post/delete?postId=1' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI"
```

## 创建评论

```shell  
curl -X POST 'http://localhost:8888/api/v1/comment/create' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI' \
  -d '{
    "userId": 2,
    "postId": 2,
    "content": "这是文章评论1"
  }'
```

## 查询评论

```shell  
curl -X GET 'http://localhost:8888/api/v1/comment/byPostId?postId=2' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoidXNlciIsImV4cCI6MTc2MjA2NTU2MH0.WJBACO6i4iboQDe8bw7py3URMwqjwj1Ul_ROFyS-jSI'
```