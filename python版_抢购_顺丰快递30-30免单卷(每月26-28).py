import requests
import json
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime

#QQ1831921455
#抓取微信小程序 或 顺丰app
url = "https://mcs-mimp-web.sf-express.com/mcs-mimp/commonPost/~memberNonactivity~memberDayFreeService~freeCouponPurchase"

# 10个不同的Cookie
cookies = [

    "token",
    "token",

   ]

headers_template = {
    "Host": "mcs-mimp-web.sf-express.com",
    "Accept": "application/json, text/plain, */*",
    "channel": "wxwd26mem1",
    "sysCode": "MCS-MIMP-CORE",
    "sw8": "1-YTkxNzI3ZDEtODdmMS00NzVjLThmMTAtOGYxNjg4MTM2MTJm-NTlhMjI4ODYtZGNlYi00ZTUwLTk0ZmItZTFkM2ZhYjJhMmQ4-0-ZmI0MDgxNzA4NWJlNGUzOThlMGI2ZjRiMDgxNzc3NDY=-d2Vi-L29yaWdpbi9hL21pbXAtYWN0aXZpdHkvbWVtYmVyRGF5-L21jcy1taW1wL2NvbW1vblBvc3Qvfm1lbWJlck5vbmFjdGl2aXR5fm1lbWJlckRheUZyZWVTZXJ2aWNlfmZyZWVDb3Vwb25QdXJjaGFzZQ==",
    "timestamp": "1719460800",
    "Accept-Language": "zh-CN,zh-Hans;q=0.9",
    "Accept-Encoding": "gzip, deflate, br",
    "platform": "MINI_PROGRAM",
    "signature": "7029fc429f82839e7c13901df3103fa2",
    "Origin": "https://mcs-mimp-web.sf-express.com",
    "User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.49(0x18003137) NetType/WIFI Language/zh_CN miniProgram/wxd4185d00bf7e08ac",
    "Referer": "https://mcs-mimp-web.sf-express.com/origin/a/mimp-activity/memberDay?mobile=187****0480&userId=9D2D05BA609249B0A810476634434723&path=/memberDayV2&linkCode=SFAC20230413163415517&supportShare=YES&from=wxwd26mem1",
    "Content-Length": "21",
    "Connection": "keep-alive",
    "Content-Type": "application/json"
}

payload = {
    # "roundTime": "15:00"
    "roundTime": "12:00"
    # "roundTime": "09:00"
}

# 记录所有响应
responses = []

# 定义发送请求的函数
def send_request(cookie):
    headers = headers_template.copy()
    headers["Cookie"] = cookie
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    return {
        "cookie": cookie,
        "status_code": response.status_code,
        "response_body": response.json()
    }



# 设定的目标时间，格式为小时:分钟:秒
# target_time = "14:59:57"
target_time = "11:59:58"
# target_time = "11:59:57"

# 等待直到设定的时间
while True:
    current_time = datetime.now().strftime("%H:%M:%S")
    print(f"Current time is {current_time}")
    if current_time == target_time:
        print(f"Current time is {current_time}. Starting requests.")
        break


# 发送请求的主循环
for _ in range(1):
    with ThreadPoolExecutor(max_workers=10) as executor:
        current_time = datetime.now().strftime("%H:%M:%S")
        print(f"发送请求 Current time is {current_time}")
        futures = [executor.submit(send_request, cookie) for cookie in cookies]
        for future in as_completed(futures):
            responses.append(future.result())
    # 可选：添加延迟以避免发送请求过快

# 打印所有响应
for i, res in enumerate(responses):
    print(f"Request {i + 1}:")
    print("Cookie:", res["cookie"])
    print("Status Code:", res["status_code"])
    print("Response Body:", res["response_body"])
    print("-" * 60)
