# 节点转换 NodeConverter
项目实现的目的：
1. 维护节点的格式定义
2. 提供一个api服务用以实现节点格式转换

## 使用
1. `GET /sub`


## 分享链接
通常分享的链接格式为：

`vless://xxx@xxx`

`ss://xxx@xxx`

`vmess://xxx@xxx`

`trojan://xxx@xxx`

收集一些常用（有足够的公信力）的节点分享链接的定义。

### VMess AEAD/VLESS
https://github.com/XTLS/Xray-core/discussions/716

### Trojan
https://p4gefau1t.github.io/trojan-go/developer/url/

### Shadowsocks
https://github.com/shadowsocks/shadowsocks-org/wiki/SIP002-URI-Scheme

## clash配置文档
https://wiki.metacubex.one/

### 收集常用的节点配置
```yaml
proxies:
  - name: vless-reality-vision                  # 可以自定义节点名称
    type: vless
    server: 1.2.3.4                             # 解析的域名或IP
    port: 12345                                 # 自定义端口
    uuid: f897325d-053d-45d1-899c-566692331f8   # 自定义 UUID
    network: tcp
    udp: true
    tls: true
    flow: xtls-rprx-vision
    servername: sega.com                        # 自定义回落域名
    reality-opts:
      public-key: 4CiE7y7ZPBXIZWzMwphuSH7qdZyisNjD3CDQGjmilmI    # Reality public-key
      short-id: a8c031ce                        # Reality short-id
    client-fingerprint: chrome                  # 自定义浏览器指纹

  - name: vless-reality-grpc                      # 可以自定义节点名称
    type: vless
    server: 1.2.3.4                               # 解析的域名或IP
    port: 12345                                   # 自定义端口
    uuid: 335ec5dd-61b1-4413-980e-5e009968f633    # 自定义 UUID
    network: grpc
    tls: true
    udp: true
    flow:
    client-fingerprint: chrome                    # 自定义浏览器指纹
    servername: sega.com                          # 自定义回落域名
    grpc-opts:
      grpc-service-name: "misaka"                 # 自定义的字符
    reality-opts:
      public-key: Aqp9oy2EFi4NNfRMZa3I3HdGhHbOIiSDZ8L28UCF73k    # Reality public-key
      short-id: 24410d1c                          # Reality short-id

  - name: vless-xtls-rprx-vision                 # 可以自定义节点名称
    type: vless
    server: www.bing.com                         # 解析的域名
    port: 12345                                  # 自定义端口
    uuid: 5f74f86b-3ee8-44f4-adc4-6666be3d315    # 自定义 UUID
    network: tcp
    tls: true
    udp: true
    flow: xtls-rprx-vision
    client-fingerprint: chrome

  - name: vless-ws-tls                               # 可以自定义节点名称
    type: vless
    server: www.bing.com                             # 解析的 IP / 域名或优选 IP / 域名
    port: 12345                                      # 自定义端口
    uuid: 3cc9a51c-db76-4ad2-a76b-8cb993bddb73       # 自定义 UUID
    udp: true
    tls: true
    network: ws
    servername: www.bing.com                         # SNI 域名，与下面 Host 一致
    ws-opts:
      path: "/?ed=2048"                              # 自定义 path 路径
      headers:
        Host: www.bing.com                           # Host 域名，与上面 server 字段的地址一致

  - name: vless-ws                                   # 可以自定义节点名称
    type: vless
    server: www.bing.com                             # 解析的 IP / 域名或优选 IP / 域名
    port: 8880                                       # 自定义端口
    uuid: 77a571fb-4fd2-4b37-8596-1b7d9728bb5c       # 自定义 UUID
    udp: true
    tls: false
    network: ws
    servername: www.bing.com                         # SNI 域名，与下面 host 一致
    ws-opts:
      path: "/?ed=2048"                              # 自定义 path 路径
      headers:
        Host: www.bing.com                           # Host 域名，与上面 server 字段的地址一致

  - name: vmess-ws-tls                               # 可以自定义节点名称
    type: vmess
    server: www.bing.com                             # 解析的 IP / 域名或优选 IP / 域名
    port: 12345                                      # 自定义端口
    uuid: 3cc9a51c-db76-4ad2-a76b-8cb993bddb73       # 自定义 UUID
    alterId: 0
    cipher: auto
    udp: true
    tls: true
    network: ws
    servername: www.bing.com                         # SNI 域名，与下面 host 一致
    ws-opts:
      path: "/?ed=2048"                              # 自定义 path 路径
      headers:
        Host: www.bing.com                           # Host 域名，与上面 server 字段的地址一致

  - name: vmess-ws                                   # 可以自定义节点名称
    type: vmess
    server: www.bing.com                             # 解析的 IP / 域名或优选 IP / 域名
    port: 8880                                       # 自定义端口
    uuid: 77a571fb-4fd2-4b37-8596-1b7d9728bb5c       # 自定义 UUID
    alterId: 0
    cipher: auto
    udp: true
    tls: false
    network: ws
    servername: www.bing.com                         # SNI 域名，与下面 Host 一致
    ws-opts:
      path: "/?ed=2048"                              # 自定义 path 路径
      headers:
        Host: www.bing.com                           # Host 域名，与上面 server 字段的地址一致

  - name: trojan-tcp-tls                             # 可以自定义节点名称
    type: trojan
    server: www.bing.com                             # 解析的域名
    port: 12345                                      # 自定义端口
    password: 123456789                              # 自定义认证密码
    client-fingerprint: chrome
    udp: true
    sni: www.bing.com                                # SNI 域名，与上面 server 字段的地址一致
    alpn:
      - h2
      - http/1.1
    skip-cert-verify: false

  - name: shadowsocks                                # 可以自定义节点名称
    type: ss
    server: www.bing.com                             # 解析的 IP / 域名
    port: 443                                        # 自定义端口
    cipher: aes-128-gcm                              # 自定义加密方式，详细请查阅 Clash Meta 文档
    password: password                               # 自定义认证密码
    udp: true
    udp-over-tcp: false
    udp-over-tcp-version: 2
    ip-version: ipv4                                 # IP 协议版本，如节点 IP 为 IPv6 则填写 ipv6
    smux:
      enabled: false

  - name: shadowsocks-shadowtls                      # 可以自定义节点名称
    type: ss
    server: 1.2.3.4                                  # 服务器本地 IP
    port: 443                                        # 自定义端口
    cipher: aes-128-gcm                              # 自定义加密方式，详细请查阅 Clash Meta 文档
    password: password                               # 自定义认证密码
    udp: true
    udp-over-tcp: false
    udp-over-tcp-version: 2
    ip-version: ipv4                                 # IP 协议版本，如节点 IP 为 IPv6 则填写 ipv6
    smux:
      enabled: false
    plugin: shadow-tls
    client-fingerprint: chrome                       # 自定义浏览器指纹
    plugin-opts:
      host: cloud.tencent.com                        # 自签证书的三方域名
      password: shadow_tls_password                  # ShadowTLS 认证密码
      version: 3                                     # ShadowTLS 协议，支持 1 / 2 / 3

  - name: hysteria1                                  # 可以自定义节点名称
    type: hysteria
    server: 1.2.3.4                                  # 服务器本地 IP
    port: 12345                                      # 自定义端口，如使用端口跳跃则改为 ports: 1000,2000-3000
    auth-str: 123456                                 # 自定义认证密码
    alpn:
      - h3
    protocol: udp                                    # 自定义协议：udp / wechat-video / faketcp
    up: 20                                           # 自定义带宽上传限制
    down: 100                                        # 自定义带宽下载限制
    sni: www.bing.com                                # SNI 域名或自签证书的三方域名
    skip-cert-verify: true                           # 使用自签证书请保持此处为 true，如为 CA 证书建议修改为 false
    fast-open: true

  - name: hysteria2                                  # 节点名称
    type: hysteria2
    server: 1.1.1.1                                  # 服务器 IP
    port: 1234                                       # 节点端口，如使用端口跳跃则改为 ports: 2000-3000/1000
    password: aa112233                               # 节点认证密码
    sni: www.bing.com                                # SNI 域名或自签证书的三方域名
    skip-cert-verify: true                           # 使用自签证书请保持此处为 true，如为 CA 证书建议修改为 false

  - name: tuic-V4                                    # 可以自定义节点名称
    server: www.bing.com                             # 解析的域名或 IP
    port: 12345                                      # 自定义端口
    type: tuic
    token: a806923b-737c-4581-8b13-56666f911866      # 自定义 Token
    alpn: [h3]
    disable-sni: true
    reduce-rtt: true
    udp-relay-mode: native
    congestion-controller: bbr

  - name: tuic-V5                                    # 可以自定义节点名称
    server: www.bing.com                             # 解析的域名或 IP
    port: 12345                                      # 自定义端口
    type: tuic
    uuid: a806923b-737c-4581-8b13-56666f911866       # 自定义 UUID
    password: a806923b-737c-4581-8b13-56666f911866   # 自定义认证密码
    alpn: [h3]
    disable-sni: true
    reduce-rtt: true
    udp-relay-mode: native
    congestion-controller: bbr

  - name: warp-wireguard                                       # 可以自定义节点名称
    type: wireguard
    server: 162.159.193.10                                     # 可自定义优选 EndPoint IP，与下方端口相对应
    port: 2408                                                 # 可自定义优选 EndPoint IP，与上方 IP 相对应
    ip: 172.16.0.2
    ipv6: 2606:4700:190:814e:7de3:5ddb:9d3e:9359               # warp 的私有 ipv6 地址，如删除本行，表示仅IPV4
    public-key: bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=
    private-key: gK3C8ijdVlT7sd5fsdf5ssdfgsdfgsdfgobT2U+rgHo=  # 获取 warp 的 私钥
    udp: true


```