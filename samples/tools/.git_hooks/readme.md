# git hook

## git hook分类
Git hook分为客户端hooks（Client-Side Hooks）和服务端hooks（Server-Side Hooks），下面列出了所有可以触发hook的时机，可以在[官方文档](https://git-scm.com/docs/githooks)中查询：

- Client-Side Hooks
```
    pre-commit: 执行git commit命令时触发，常用于检查代码风格
    prepare-commit-msg: commit message编辑器呼起前default commit message创建后触发，常用于生成默认的标准化的提交说明
    commit-msg: 开发者编写完并确认commit message后触发，常用于校验提交说明是否标准
    post-commit: 整个git commit完成后触发，常用于邮件通知、提醒
    applypatch-msg: 执行git am命令时触发，常用于检查命令提取出来的提交信息是否符合特定格式
    pre-applypatch: git am提取出补丁并应用于当前分支后，准备提交前触发，常用于执行测试用例或检查缓冲区代码
    post-applypatch: git am提交后触发，常用于通知、或补丁邮件回复（此钩子不能停止git am过程）
    pre-rebase: 执行git rebase命令时触发
    post-rewrite: 执行会替换commit的命令时触发，比如git rebase或git commit --amend
    post-checkout: 执行git checkout命令成功后触发，可用于生成特定文档，处理大二进制文件等
    post-merge: 成功完成一次 merge行为后触发
    pre-push: 执行git push命令时触发，可用于执行测试用例
    pre-auto-gc: 执行垃圾回收前触发
```

- Server-Side Hooks
```
    pre-receive: 当服务端收到一个push操作请求时触发，可用于检测push的内容
    update: 与pre-receive相似，但当一次push想更新多个分支时，pre-receive只执行一次，而此钩子会为每一分支都执行一次
    post-receive: 当整个push操作完成时触发，常用于服务侧同步、通知
```


如何使用git hook
hook脚本会存放在仓库.git/hooks文件夹中，git提供了一些shell样例脚本以作参考。
