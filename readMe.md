# 遇到的一些问题
## go:embed 指令解释
go:embed 是 Go 语言 1.16 版本引入的一个编译指令，用于将静态文件（如文本文件、图片、配置文件等）嵌入到 Go 二进制文件中。

具体来说，go:embed 允许你在编译时将文件内容转换为字节切片或字符串，并将其存储在 Go 程序的变量中。
### 使用场景
- 资源文件嵌入：将 HTML、CSS、JavaScript 文件等嵌入到 Web 应用中。
- 配置文件嵌入：将配置文件嵌入到程序中，避免运行时依赖外部文件。
- 模板文件嵌入：将模板文件嵌入到程序中，方便模板渲染。
```go
//go:embed <file-or-pattern>
var <name> <type>
```
<file-or-pattern>：文件路径或模式匹配。可以是单个文件路径，也可以是通配符模式（如 *.txt）。

< name>：变量名，用于存储文件内容。

< type>：变量类型，通常是 string 或 []byte。
## 如何调用go:embed中lua脚本
在redis中调用lua脚本，应该是由不同编程编程语言的实现来调用的，比如lua脚本调用go语言的实现，或者go语言调用lua语言的实现。
在go语言中，我们可以使用redis-go-client库来调用redis服务器，然后使用lua脚本来执行相应的操作。以下是一个简单的例子：

在 Redis 中，EVAL 命令用于执行 Lua 脚本。它允许你在 Redis 服务器端运行 Lua 脚本，从而实现复杂的操作逻辑。下面详细解释一下 EVAL 函数的调用方式及其参数关系。

1. EVAL 命令格式
```shell
EVAL script numkeys key [key ...] arg [arg ...]
script: Lua 脚本的内容，作为字符串传递。
numkeys: 表示接下来的参数中有多少个是键（keys），这些键会被传递给 Lua 脚本中的 KEYS 数组。
key [key ...]: 这些是你要操作的 Redis 键，它们会按顺序被存储在 Lua 脚本中的 KEYS 数组中。
arg [arg ...]: 这些是额外的参数，它们会按顺序被存储在 Lua 脚本中的 ARGV 数组中。
```


lua脚本：
```lua
-- 获取传入的key，用于标识验证码
-- phone_code:login:152xxxx
local key = KEYS[1]
-- 构造计数器的key
-- phone_code:login:152xxxx:cnt
local cntKey = key..":cnt"
-- 你准备的存储的验证码
local val = ARGV[1]

-- 获取key的剩余生存时间
local ttl = tonumber(redis.call("ttl", key))
-- 检查key的生存时间，决定是否发送验证码
if ttl == -1 then
    --    key 存在，但是没有过期时间 
    return -2
elseif ttl == -2 or ttl < 540 then
    --    可以发验证码 , 540 = 600 -60 ，超过一分钟了 ， 并且key不存在，key不存在的时候，ttl == -2 
    redis.call("set", key, val)
    -- 600 秒
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end
-- KEYS[1]: 对应于 EVAL 命令中的第一个键参数，即 key。这个键是用来存储验证码的主键。
-- ARGV[1]: 对应于 EVAL 命令中的第一个额外参数，即 val。这个参数是要存储的验证码值。
```
调用lua
```go
// 使用Lua脚本通过Eval方法执行设置操作，该操作在Redis中实现。
	
res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()

/**
1. ctx
类型: context.Context
作用: 用于控制请求的生命周期，可以用来取消请求或设置超时。通常在并发编程中使用，确保请求在必要时可以被中断。
2. luaSetCode
类型: string
作用: 包含 Lua 脚本的内容。这是一个字符串，表示你希望在 Redis 中执行的 Lua 脚本。例如：
3. []string{c.key(biz, phone)}
类型: []string
作用: 这是一个字符串切片，包含传递给 Lua 脚本的键（keys）。这些键会被存储在 Lua 脚本中的 KEYS 数组中。
具体含义:
c.key(biz, phone) 是一个方法调用，返回一个字符串，表示 Redis 中的键。这个键通常用于标识特定的业务和手机号码。
[]string{c.key(biz, phone)} 将这个键包装成一个切片，传递给 EVAL 命令。
4. code
类型: string
作用: 这是一个额外的参数，传递给 Lua 脚本。这个参数会被存储在 Lua 脚本中的 ARGV 数组中。
具体含义:
code 是要存储的验证码值。
 .Int()
作用: 这是一个方法调用，用于将 Eval 返回的结果转换为 int 类型。
具体含义:
Eval 方法返回一个 redis.Cmd 对象，这个对象包含 Redis 命令的执行结果。
.Int() 方法将这个结果转换为整数。根据 Lua 脚本的返回值，可能的整数结果有：
0: 成功设置验证码。
-1: 发送太频繁。
-2: 键存在但没有过期时间。

*/
```


