# RPCM
rpcm(rpc-mesh) 是 mesh service 环境下的 rpc 库。

# CompactT
- rpcx
    + 主要特点：
        - golang 开发，简单易用；
        - 性能远远高于 Dubbo、Motan、Thrift等框架，是gRPC性能的两倍
        - 比较完善的服务发现和服务治理
            + 支持 Failover、Failfast、Failtry、Backup等失败模式，支持 随机、 轮询、权重、网络质量, 一致性哈希,地理位置等路由算法
        - 支持的网络类型丰富：tcp、http、unix、quic、kcp
    + reference: [go-rpc-programming-guide](https://books.studygolang.com/go-rpc-programming-guide/part2/quickstart.html)
- Motan
    + 主要特点：
        - 据说在新浪微博正支撑着千亿次调用
        - 轻量级，开发和学习简单
    + reference： [motan-go](https://github.com/weibocom/motan-go)
