 ## 狼人杀需求

 ### 服务初始化之数据库服务
 1. user : 用户数据
> 1.1 自增长id <br> 
1.2 openId 微信用户唯一标识 <br>
1.3 用户等级-用户经验条 -- 用户经验可以由前端判断获取,等级初始化1级,50级后升级经验不增加 <br>
1.4  金币 - 初始化为300金币 <br>
1.5 其它数据存储，json转字符串存储 

2. store: 商店
> 1.1 自增长id <br>
1.2 名字 <br>****
1.3 价格 <br>
1.4 属性  是消耗品还是服装 <br>
1.5 数量  消耗品数量 <br>
1.6 服装期限 

3. 好友关系
> 1.1 自增长id <br>
1.2 用户id <br>
1.3 好友id 一对多关系<br>
1.4 删除好友

4. 战绩信息
> 1.1 自增id<br>
1.2 用户id<br>
1.3 游戏开局时间<br>
1.4 身份<br>
1.5 游戏模式<br>
1.6 总场次<br>
1.7 连胜数 - 最高连胜<br>
1.8 逃跑<br>
1.9 胜率<br>
1.10 输赢<br>



### websocket 监听数据格式,接收请求
> socket 服务接收格式为json转成字符串形式,其中存在两个数据内容格式如下： {“name”:"请求地址，要实现的功能","values":{传输的必须数据，请求数据,json转字符串}}
返回数据格式如下:{"status":"ok","mes":"提示信息",...."data":{ 其他json数据}} 为字符串格式

1. login : 用户登录
> 1.1 OpenId: 微信用户标识必须<br>
1.2 NickName： 用户昵称必须<br>
1.3 AvatarURL: 用户头像必须<br>
1.4 Orther: 需要保存的其他信息,可选<br>
示例: {"name":"login","values":"{\"nickName\":\"smal\",\"openID\":\"12345\",\"avatarUrl\":\"https:\"}"}<br>
返回数据: {"status":"ok","mes":"登录成功","data":{"openID":"12345","nickName":"smal","avatarUrl":"https:","level":"1","money":"300","orther":"","id":"1"}}


2. back : 获取背包
> 1.1 OpenId : 微信用户标识必须<br>
示例: {"name":"back","values":"{\"openID\":\"12345\"}"}<br>
返回数据:  {"status":"ok","mes":"获取背包成功","data":[{"name":"基础下装","property":"0","stilTime":""},{"name":"基础下装","property":"0","stilTime":""},{"name":"基础帽子","property":"0","stilTime":""}]}

3. upgrade : 升级和增加金币，或者用户修改昵称，头像
> 1.1 OpenId : 微信用户标识必须<br>
1.2 Money: 金币数量，用于购买减少，对战获取增加,非必须<br>
1.3 nickName: 非必须<br>
1.4 AvatarURL : 非必须<br>
示例:{"name":"upgrade","values":"{\"openID\":\"12345\"}"}<br>
返回数据:{"status":"ok","mes":"升级成功","data":{"openID":"12345","nickName":"smal","avatarUrl":"https:","level":"3","money":"300","orther":"","id":"1"}}

4. addback : 增加背包
> 1.1 User: int格式 用户ID 必须<br>
1.2 Name: 名字必须<br>
1.3 property: int 属性必须<br>
1.4 stilTime : 有效时间 可选<br>
示例: {"name":"addback","values":"{\"User\":1,\"Name\":\"高级上衣\",\"property\":0}"}
返回数据: {"status":"ok","mes":"添加背包成功"}

5. record : 获取最近战绩
> 1.1 OpenId : 微信用户标识必须<br>
示例: {"name":"record","values":"{\"openID\":\"12345\"}"}<br>
返回数据: {"status":"ok","mes":"获取战绩成功","data":[]}


6. buddy : 获取好友列表
> 1.1 OpenId : 微信用户标识必须<br>
示例: {"name":"buddy","values":"{\"openID\":\"12345\"}"}<br>
返回数据:{"status":"ok","mes":"获取好友列表成功","data":[]}

7. recordRate: 获取胜率等
> 1.1 OpenId : 微信用户标识必须<br>
示例:{"name":"recordRate","values":"{\"openID\":\"12345\"}"}<br>
返回数据:{"status":"ok","mes":"获取战绩成功","data":["count":"0","runAway":"0","maxWin":"0","winRate":"0"]}

8. newbuddy : 获取好友申请
> 1.1 OpenId : 微信用户标识必须<br>
示例：{"name":"newbuddy","values":"{\"openID\":\"12345\"}"}<br>
返回数据:{"status":"ok","mes":"获取好友申请成功","data":[]}

9. agreebuddy: 同意好友申请
> 1.1 buddy.Id int 必须

10.rcombuddy: 获取好友推荐
> 1.1 OpenId : 微信用户标识必须<br>
示例:{"name":"rcombuddy","values":"{\"openID\":\"12345\"}"}<br>
返回数据:{"status":"ok","mes":"获取推荐好友","data":[{"openID":"12345","nickName":"","avatarUrl":"","level":"0","id":"0"}]}

11.addbuddy: 添加好友
