use 
```shell
export ALIYUN_ACCESSKEYID='your'
export ALIYUN_ACCESSSECRET='your'
export DOMAINNAME='your domain'
```

add to cron:
```shell
1 * * * * . $HOME/.profile; cd SOMEPATH; ./ddns > /tmp/alidns
```
