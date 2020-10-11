use 
```shell
export ALIYUN_ACCESSKEYID='your'
export ALIYUN_ACCESSSECRET='your'
export DOMAINNAME='your domain'
export DOMAINID='your domain id'
```

add to corn:
```shell
1 * * * * . $HOME/.profile; cd SOMEPATH; ./ddns > /tmp/alidns
```
