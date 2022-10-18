User Case
```
go install github.com/gwaycc/bchain-storage/cmd/bchain-storage

mkdir -p /data/zfs
bchain-storage daemon --export-nfs=true &

# reset the admin password
bchain-storage sys resest-passwd admin # will output the new passwd of admin

# add a common user
bchain-storage --passwd=[new passwd] sys adduser user1 # will output the passwd of user1

# upload files
bchain-storage --user=user1 --passwd=[password of user1] upload [local path] [remote path]

# download files
bchain-storage --user=user1 --passwd=[password of user1] download [remote path] [local path] 

# system status
bchain-storage sys status
```

Http api
```shell
# system status
curl -k "https://127.0.0.1:1330/check"

# more pool
curl -k "https://127.0.0.1:1330/check" # 对应openz-0
curl -k "https://127.0.0.1:1340/check" # 对应openz-1
curl -k "https://127.0.0.1:1350/check" # 对应openz-2

# more api
https://github.com/gwaycc/bchain-storage/blob/main/cmd/bchain-storage/client/http_file.go#L81
```

Reset password
```
# Way 1:
# reset user by api
bchain-storage --user=admin --passwd=[admin passwd] sys reset-passwd [user] # output new user password

# Way 2:
# reset user by local mode
# stop the daemon,
# then reset the password. 
bchain-storage sys reset-passwd --local=true --repo=[repo path] [user] 
```

TODO
```
Apache Lisence
implement userspace
implement fuse
implement s3 interface
implement cluster
```
