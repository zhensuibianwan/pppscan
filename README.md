## pppscan

go+wails 编写的漏洞、指纹管理扫描工具

## 免责声明

本工具截图所进行的演示均在本地环境或授权情况下进行，且本工具不包含任何权限级别的漏洞利用poc，仅作为企业或个人资产漏洞自查的安全建设工具。在使用本工具时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权，请勿对非授权目标进行访问。如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。请勿将本项目技术或代码应用在恶意软件制作、软件著作权/知识产权盗取或不当牟利等非法用途中。

## 使用说明

### 初始化

开始使用时，请本地搭建一个Mysql数据库（或者远程的数据库），在设置这填写并点击保存。随后点击初始化，应用会自动创建数据库和表，随后便可正常使用了。

![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291036548.png)
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291110700.png)


### POC管理

#### 添加POC

添加的poc数据库和本地同时都会保存一个。添加POC时建议先去指纹管理里添加对应指纹，这样就可以直接将指纹与该poc绑定，到时指纹识别时，扫到该指纹时会自动联动poc进行扫描。
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291112910.png)


只需要把POC发送的请求包贴到左边即可，请求包格式要对，最后是从bp，yakit中复制出来的，Host头可以不带。
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291116452.png)
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291122684.png)

这里特别说明一下状态码，如果状态码为空（即什么都不填）且为多个请求包的poc，则即使当前请求包通过判断值判断该漏洞不存在，程序也会自动发送下一个请求包，并判断漏洞是否存在。所以可以通过不填状态码，编写平行请求的poc，但是判断值一定要准确，否则容易误报。比如若依弱口令，一般存在两个默认口令，admin/admin123、ry/admin123。这是你就可以这样写，如下图：
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291133036.png)
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291133172.png)

添加值是为了适应多样化的poc编写的，目前内置了4个。

![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291138186.png)

request.url.0 ：poc扫描时会自动替换该值为扫描的url
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291145020.png)

编码：16进制编码，为了适应一些反序列化的漏洞，需要先将playload编码放在poc中，在扫描时程序会自动解码
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291148375.png)

Content-LengthMax和Content-LengthMin一般用于sql注入时间盲注漏洞。（用处不大，聊胜于无）
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291150235.png)

自定义添加值，在设置处添加，目前能自定义添加三种。
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291154866.png)

input：设置漏洞利用的输入值，比如命令执行的命令，这里不方便展开，自行摸索。

request1.body.1 ： request后面的数字指的是第几个请求包，body指的是返回包的body，后面的数字是区分不同的request1.body，添加方式也很简单，只需要将你需要的值的连带前后关键字一起复制粘贴到右边，然后选中你需要的值点击替换即可，程序会自动将poc中的这个值替换成右边“~”所代表的值。最终效果如下：
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291202053.png)
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291204183.png)

request1.header.1 ：用法和request1.body.1差不多，区别就是右边填的是你需要的header头，参考如下：
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291213770.png)
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291214651.png)


#### 删除POC

删除会连同数据库加本地文件一起删除。

#### 导入POC

会显示应用程序同目录下的poc目录里数据库里没有的poc，然后选择导入即可。

#### 保存POC

会将你选择的poc保存到本地，保存位置为应用程序同目录下的poc目录

### 指纹管理

和poc管理差不多，唯一注意的是，指纹中的faviconHash为32位的
![](https://yuexiaduzhuo.oss-cn-nanjing.aliyuncs.com/pppscan/202403291345300.png)


## 结尾

工具刚开发出来没多久，可能存在各种bug，希望大家使用过程中遇到后给我指出，同时也欢迎各位大佬提出修改意见。
最后，感谢AtomsTeam中小伙伴的支持。
