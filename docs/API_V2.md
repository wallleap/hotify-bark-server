# API V2

**The V2 version API is switched to the standard REST request, and most of the compatibility
processing has been done for the V1 version API; users should use the new REST API when using
the V2 version.**

> 地址中的 `18080` 是容器对外映射端口（`-p 18080:8080`）。若直接运行二进制（默认监听
> `0.0.0.0:8080`），请使用 `8080`。

- [API V2](#api-v2)
    * [Push](#push)
        + [curl](#curl)
        + [Batch push (curl)](#batch-push-curl)
        + [golang](#golang)
        + [python](#python)
        + [java](#java)
        + [nodejs](#nodejs)
        + [php](#php)
    * [参数说明与优先级](#参数说明与优先级)
    * [响应格式](#响应格式)
    * [认证](#认证)
    * [新增接口（本 fork）](#新增接口本-fork)
    * [Misc](#misc)
        + [Ping](#ping)
        + [Healthz](#healthz)
        + [Info](#info)

## Push

| Field | Type | Description |
| ----- | ---- | ----------- |
| id (optional) | string | Notification collapse id (APNs `apns-collapse-id`) |
| title (optional) | string | Notification title (font size would be larger than the body) |
| subtitle (optional) | string | Notification subtitle |
| body  | string | Notification content |
| device_key | string | The key for each device |
| device_keys (optional) | array | Used for batch pushing |
| level (optional) | string | `'critical'`, `'active'`, `'timeSensitive'`, `'passive'` |
| volume (optional) | string | The ringtone volume for critical alert notification. |
| badge (optional) | integer | The number displayed next to App icon ([Apple Developer](https://developer.apple.com/documentation/usernotifications/unnotificationcontent/1649864-badge)) |
| call (optional) | string | Must be `1`, The ringtone will continue to play for 30 seconds |
| autoCopy (optional) | string | Must be `1` |
| copy (optional) | string |  The value to be copied |
| sound (optional) | string | Value from [here](https://github.com/Finb/Bark/tree/master/Sounds)， and custom ringtones are also available |
| icon (optional) | string | An url to the icon, available only on iOS 15 or later |
| group (optional) | string | The group of the notification |
| ciphertext (optional) | string | The ciphertext of encrypted push notifications |
| isArchive (optional) | string | Value must be `1`. Whether or not should be archived by the app |
| ttl (optional) | integer | Time to live for archived messages, in seconds. Expired archived messages are deleted automatically |
| url (optional) | string | Url that will jump when click notification |
| action (optional) | string | Set to "none", tap notifications do nothing |
| delete (optional) | string | Must be `1`. Silent push: delivers a background push and is not shown on screen |

### curl

```sh
curl -X "POST" "http://127.0.0.1:18080/push" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "body": "Test Bark Server",
  "device_key": "ynJ5Ft4atkMkWeo2PAvFhF",
  "title": "bleem",
  "badge": 1,
  "sound": "minuet",
  "icon": "https://day.app/assets/images/avatar.jpg",
  "group": "test",
  "url": "https://mritd.com"
}'
```

### Batch push (curl)

批量推送：`device_keys` 传数组（或逗号分隔字符串）。每个设备各推送一条，返回逐设备的推送结果。

```bash
curl -X "POST" "http://127.0.0.1:18080/push" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "title": "batch",
  "body": "Test Bark Server",
  "device_keys": ["ynJ5Ft4atkMkWeo2PAvFhF", "nysrshcqielvoxsa"]
}'
```

响应（批量时为 `data` 数组）：

```json
{
  "code": 200,
  "message": "success",
  "timestamp": 1700000000,
  "data": [
    {"code": 200, "device_key": "ynJ5Ft4atkMkWei2PAvFhF"},
    {"code": 410, "device_key": "nysrshcqielvoxsa", "message": "push failed: ..."}
  ]
}
```

### golang

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

func sendPush() {
	// push (POST http://127.0.0.1:18080/push)

	json := []byte(`{"body": "Test Bark Server","device_key": "nysrshcqielvoxsa","title": "bleem", "badge": 1, "icon": "https://day.app/assets/images/avatar.jpg", "group": "test", "url": "https://mritd.com","sound": "minuet"}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://127.0.0.1:18080/push", body)
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
```

### python

```python
# Install the Python Requests library:
# `pip install requests`

import requests
import json


def send_request():
    # push
    # POST http://127.0.0.1:18080/push

    try:
        response = requests.post(
            url="http://127.0.0.1:18080/push",
            headers={
                "Content-Type": "application/json; charset=utf-8",
            },
            data=json.dumps({
                "body": "Test Bark Server",
                "device_key": "nysrshcqielvoxsa",
                "title": "bleem",
                "sound": "minuet",
                "badge": 1,
                "icon": "https://day.app/assets/images/avatar.jpg",
                "group": "test",
                "url": "https://mritd.com"
            })
        )
        print('Response HTTP Status Code: {status_code}'.format(
            status_code=response.status_code))
        print('Response HTTP Response Body: {content}'.format(
            content=response.content))
    except requests.exceptions.RequestException:
        print('HTTP Request failed')
```

### java

```java
import java.io.IOException;
import org.apache.http.client.fluent.*;
import org.apache.http.entity.ContentType;

public class SendRequest
{
  public static void main(String[] args) {
    sendRequest();
  }
  
  private static void sendRequest() {
    
    // push (POST )
    
    try {
      
      // Create request
      Content content = Request.Post("http://127.0.0.1:18080/push")
      
      // Add headers
      .addHeader("Content-Type", "application/json; charset=utf-8")
      
      // Add body
      .bodyString("{\"body\": \"Test Bark Server\",\"device_key\": \"nysrshcqielvoxsa\",\"title\": \"bleem\",\"url\": \"https://mritd.com\", \"group\": \"test\",\"sound\": \"minuet\"}", ContentType.APPLICATION_JSON)
      
      // Fetch request and return content
      .execute().returnContent();
      
      // Print content
      System.out.println(content);
    }
    catch (IOException e) { System.out.println(e); }
  }
}
```

### nodejs

```node
// request push 
(function(callback) {
    'use strict';
        
    const httpTransport = require('http');
    const responseEncoding = 'utf8';
    const httpOptions = {
        hostname: '127.0.0.1',
        port: '18080',
        path: '/push',
        method: 'POST',
        headers: {"Content-Type":"application/json; charset=utf-8"}
    };
    httpOptions.headers['User-Agent'] = 'node ' + process.version;
 
    // Using Basic Auth {"username":"","password":""}
    // Paw Store Cookies option is not supported

    const request = httpTransport.request(httpOptions, (res) => {
        let responseBufs = [];
        let responseStr = '';
        
        res.on('data', (chunk) => {
            if (Buffer.isBuffer(chunk)) {
                responseBufs.push(chunk);
            }
            else {
                responseStr = responseStr + chunk;            
            }
        }).on('end', () => {
            responseStr = responseBufs.length > 0 ? 
                Buffer.concat(responseBufs).toString(responseEncoding) : responseStr;
            
            callback(null, res.statusCode, res.headers, responseStr);
        });
        
    })
    .setTimeout(0)
    .on('error', (error) => {
        callback(error);
    });
    request.write("{\"device_key\":\"nysrshcqielvoxsa\",\"body\":\"Test Bark Server\",\"title\":\"bleem\",\"sound\":\"minuet\",\"url\":\"https://mritd.com\", \"group\":\"test\"}")
    request.end();
    

})((error, statusCode, headers, body) => {
    console.log('ERROR:', error); 
    console.log('STATUS:', statusCode);
    console.log('HEADERS:', JSON.stringify(headers));
    console.log('BODY:', body);
});
```

### php

```php
$curl = curl_init();
curl_setopt_array($curl, [
    CURLOPT_URL => 'http://127.0.0.1:18080/push',
    CURLOPT_CUSTOMREQUEST => 'POST',
    CURLOPT_POSTFIELDS => '{
  "title": "bleem",
  "device_key": "nysrshcqielvoxsa",
  "body": "Test Bark Server",
  "badge": 1,
  "sound": "minuet",
  "icon": "https://day.app/assets/images/avatar.jpg",
  "group": "test",
  "url": "https://mritd.com"
}',
    CURLOPT_HTTPHEADER => [
        'Content-Type: application/json; charset=utf-8',
    ],
]);
$response = curl_exec($curl);
curl_close($curl);
echo $response;
```

## 参数说明与优先级

- `device_keys` 省略时为单设备推送（此时必须提供 `device_key`，否则 `400 device key is empty`）。
- 除上表字段外，其余字段会原样透传为 APNs 自定义字段（`payload.custom`），key 转小写。
- 3 个取值来源按优先级覆盖（低→高）：**请求体（body） → URL 查询参数（query） → URL 路径**。

## 响应格式

所有接口统一返回 `CommonResp`：

| Field | Type | Description |
| ----- | ---- | ----------- |
| code | int | 业务码（成功 `200`） |
| message | string | 提示信息（成功 `"success"`，失败为错误描述） |
| data (optional) | json | 附加数据（如批量推送结果、注册返回的 key） |
| timestamp | int | unix 秒级时间戳 |

- 成功单设备推送：`{"code":200,"message":"success","timestamp":...}`
- 失败：HTTP 状态码与 `code` 一致，如 `400`/`410`/`404` 等。

## 认证

- `/push`、`/:device_key`（V1 兼容推送）、`/mcp*` **默认无独立认证**——`device_key` 本身就是凭证，需保证其私密性。
- 启用 Basic Auth（`--user` / `--password` 或 `BARK_SERVER_BASIC_AUTH_USER/PASSWORD`）后，受保护子路径需携带 Basic Auth 头；白名单放行路径见 [README](../README.md)。
- 推送请求中可通过 query 传递参数，Basic Auth 开启时建议在 `Authorization` 头中携带凭证。

## 新增接口（本 fork 新增）

除上游 V2 `/push` 外，本 fork 额外提供：

| Method | Path | 认证 | 说明 |
|---|---|---|---|
| GET | `/` | 无 | 存活探测，返回 `"ok"` |
| POST | `/register` | 无 | 设备注册（body），返回 `device_key`。详见 [TUTORIAL.md](TUTORIAL.md) |
| GET | `/register` | 无 | 设备注册（兼容旧 query 参数） |
| GET | `/register/:device_key` | 无 | 校验 device_key 是否存在 |
| GET | `/version` | 无 | Gotify 兼容探测，返回 `{"version":...}` |
| GET | `/message` | 客户端 token | Gotify 兼容历史消息查询 |
| DELETE | `/message` / `/message/:id` | 客户端 token | Gotify 兼容删除历史 |
| GET | `/stream` | 客户端 token | Gotify 兼容 WebSocket 订阅 |
| GET/POST | `/mcp` / `/mcp/:device_key` | 无 | MCP 接口，供 AI 代理直接推送 |

Gotify 兼容与 MCP 接口的详细说明见 [GOTIFY_COMPAT.md](GOTIFY_COMPAT.md) 与 [MCP.md](MCP.md)。

## Misc

### Ping

```sh
curl "http://127.0.0.1:18080/ping"
```

### Healthz

```sh
curl "http://127.0.0.1:18080/healthz"
```

### Info

```sh
curl "http://127.0.0.1:18080/info"
```