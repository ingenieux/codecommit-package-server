# codecommit-package-server

A Golang Package Server for CodeCommit Repositories. Say, if you've got myrepo on us-east-1, you can do this:

```
$ go get codecommit.ingenieux.io/repo/myrepo
```

Or if you've got in a region other than us-east-1, you could do this:

```
$ go get codecommit.ingenieux.io/otherregion/repo/myrepo
```

And thats it! Keep reading if you want to build your own package server.

## Installation

```
$ go get github.com/ingenieux/codecommit-package-server 
```

This will download and install it, next:

  * Setup your systemd scripts such as:

```
$ cat /etc/systemd/system/codecommit-package-server.service 
[Unit]
Description=CodeCommit Import Server
After=network.target

[Service]
User=ubuntu
Group=ubuntu
ExecStart=/home/ubuntu/go/bin/codecommit-package-server
# Fix the path above
 
[Install]
WantedBy=multi-user.target
```
  * Setup nginx:

```
$ cat /etc/nginx/sites-enabled/codecommit.ingenieux.io
server {
  listen 80;
  listen [::]:80;
  server_name codecommit.ingenieux.io;

  access_log /var/log/nginx/nginx.codecommit.ingenieux.io.access.log;
  error_log /var/log/nginx/nginx.codecommit.ingenieux.io.error.log;

  location /.well-known/ {
    root /var/www/letsencrypt;
  }

  location / {
    return 301 https://$host$request_uri;
  }
}

# Append the lines below once TLS keys were installed. Remember to change hosts

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;

    include snippets/ssl-params.conf;

    server_name codecommit.ingenieux.io;

    access_log /var/log/nginx/nginx.codecommit.ingenieux.io.access.log;
    error_log /var/log/nginx/nginx.codecommit.ingenieux.io.error.log;

    ssl_certificate /etc/letsencrypt/live/codecommit.ingenieux.io/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/codecommit.ingenieux.io/privkey.pem;

    location / {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host            $http_host;
        proxy_pass http://localhost:3001/;
    }
}
```

  * Setup SSL with [Certbot](https://github.com/certbot/certbot/)

## Usage

Replace codecommit.ingenieux.io with your package server host. In the example below, myrepo is a :

```
$ go get -d -v codecommit.ingenieux.io/repo/myrepo
```

# What about other AWS Regions?

Err, nice question. Currently this public server defaults on us-east-1, but you prefix the region on the URL instead:

```
$ go get codecommit.ingenieux.io/us-west-2/repo/myrepo
```

# What about SSH?

You can add ?protocol=ssh into your URLS. 

Be wary that you must also change that in all your source code. 

When running a custom install, simply set "defaultProto" to "ssh" and it should work.
