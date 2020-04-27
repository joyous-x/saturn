# saturn
[![license](https://img.shields.io/github/license/joyous-x/saturn.svg)](https://github.com/joyous-x/saturn/blob/master/LICENSE)
![GitHub Actions status](https://github.com/joyous-x/saturn/workflows/code-analyzer/badge.svg)
[![release](https://img.shields.io/github/release/joyous-x/saturn.svg)](https://github.com/joyous-x/saturn/releases/latest)
[![codecov](https://codecov.io/gh/joyous-x/saturn/branch/master/graph/badge.svg)](https://codecov.io/gh/joyous-x/saturn)

## Overview
saturn is a goland universal framework for rapid development of high performance mesh services.

## Function
- ginbox: 启动管理多个http服务
    + 静态文件服务
- ip2region: 判断ip地址所在地区
- cn2pinyin: 将汉字转换为拼音 ：TODO
- alipay: 阿里支付 ： TODO
- gendao：TODO
- rbac: TODO
- TODO
    + https://github.com/facebookgo/inject


## Usage
```
    go get -v github.com/joyous-x/saturn
```
The quick start gives very basic example of running client and server on the same machine. 

For the detailed information about using and developing **saturn**, please jump to [Documents](#Documents). the demo case is in the samples/ directory


## Test
```
    go test -cover ./...
```