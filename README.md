User Case
```
go install github.com/gwaycc/bcstorage/cmd/bcstorage

mkdir -p /data/zfs
#bcstorage daemon --export-nfs=true &
bcstorage daemon &

# reset the admin password when first used
bcstorage sys resest-passwd admin # will output the new passwd of admin

# add a common user
bcstorage --passwd=[new admin passwd] sys adduser user1 # will output the passwd of user1

# upload files
bcstorage --user=user1 --passwd=[password of user1] upload [local path] [remote path]

# download files
bcstorage --user=user1 --passwd=[password of user1] download [remote path] [local path] 

# system status
bcstorage sys status
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
https://github.com/gwaycc/bcstorage/blob/main/module/client/http_file.go#L81
```

Reset password
```
# Way 1:
# reset user by api
bcstorage --user=admin --passwd=[admin passwd] sys reset-passwd [user] # output new user password

# Way 2:
# reset user by local mode
# stop the daemon,
# then reset the password. 
bcstorage sys reset-passwd --local=true --repo=[repo path] [user] 
```

TODO
```
Apache Lisence
implement userspace
implement fuse
implement s3 interface
implement cluster
```
