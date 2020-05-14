# Remote Manager

## Setup

### 0. Preparation

You need to install docker and optionally docker compose in order to use this tool.

### 1. Clone project

```bash
git clone git@github.com:lazy-ants/remote-manager.git
cd remote-manager
```

### 2. Initial setup
```bash
make init
```

### 3. Add server connection to the config

Open **config.json** with the editor of your choice and add server connections.
It means "user@domain:port", however the port is optional.

### 4. Running configuration
With docker:
```
docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager bin/console app:validate-config
```
With docker-compose:
```
docker-compose run remote-manager bin/console app:validate-config
```
#### Result

![validate](https://user-images.githubusercontent.com/249065/81821063-ba509200-9531-11ea-85eb-ef735ab42f1f.png)

### 5. Run you first command to see e.g. the server uptime

To test this setup run

With docker
```
docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager bin/console app:uptime
```

with docker composer:
```
docker-compose run remote-manager bin/console app:uptime
```

## Available commands

### app:docker-compose-version  
Get docket compose version

![docker-compose](https://user-images.githubusercontent.com/249065/81821076-bd4b8280-9531-11ea-8a7f-3f5bbf1418bd.png)


### app:docker-prune
Prune old docker data

![docker-prune](https://user-images.githubusercontent.com/249065/81821073-bd4b8280-9531-11ea-940d-cff7a7cd0c61.png)

### app:docker-ps
Show docker process status

![docker-ps](https://user-images.githubusercontent.com/249065/81821072-bcb2ec00-9531-11ea-8fd6-51513f406ab1.png)

### app:kernel
Get server kernels

![kernel](https://user-images.githubusercontent.com/249065/81821068-bb81bf00-9531-11ea-916a-02daef2b77f4.png)

### app:ls
Run ls command on all servers

#### Examples
- `app:ls` default value, with *-lha* argumetns list the current directory
- `app:ls './ -la'` list the current directory
- `app:ls '../ -la'` list the up directory
- `app:ls '/ -la'` list the root directory

![ls](https://user-images.githubusercontent.com/249065/81821048-b6bd0b00-9531-11ea-93a7-e58422c77b2b.png)

### app:os                      
Get server OS

![os](https://user-images.githubusercontent.com/249065/81821066-bb81bf00-9531-11ea-93be-cdd89c1b5f3a.png)

### app:check-reboot                  
Checks whether a reboot is required

![reboot](https://user-images.githubusercontent.com/249065/81821065-bae92880-9531-11ea-9b70-59d6d8d941d1.png)

### app:system-info
- note: [sudo required]
- Get system information

![system-info](https://user-images.githubusercontent.com/249065/81821061-ba509200-9531-11ea-8829-0fb9ca619105.png)

### app:upgrade
- note: [sudo required]
- Upgrade server packages

![upgrade](https://user-images.githubusercontent.com/249065/81821055-b91f6500-9531-11ea-8480-722a2cc65fbe.png)

### app:uptime                  
Get server uptime

![uptime](https://user-images.githubusercontent.com/249065/81821064-bae92880-9531-11ea-989d-e64e5eff0ac7.png)

### app:validate-config         
Validate server instances config

![validate](https://user-images.githubusercontent.com/249065/81821063-ba509200-9531-11ea-85eb-ef735ab42f1f.png)

## FAQ

### Login into console

If you want to login into the docker container:
```
reman-cli bash
```

### Running command needed the sudo password

Prepare servers you want to manage:

First of all, you need a possibility to provide the sudo password in a secure way as an environment variable to your server (for sure in case you need the possibility of runnnig commands with sudo).
Therefore, On each server you want to manage edit 
```
sudo nano /etc/ssh/sshd_config
``` 
add at the end of the config:
```
AcceptEnv PASSWORD
```
reload the sshd server e.g.
```
sudo service sshd reload
```

After that you either need to provide the sudo password in the .env.local file in case the most of your servers use the same sudo password or you can add in the config.json file in each configuration sectiona using the key "sudo-ppassword".  

### Adding aliases
In case you want to simplify using this tool you may want to add aliases.

In case you have OS specific options like "UseKeychain yes" on MacOS in your .ssh/config: 
```
alias reman-cli='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager'
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager ./run'
```

Otherwise you can simple map the whole .ssh directory:
```
alias reman-cli='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager'
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager ./run'
```

Or if you prefer docker-compose:

```
alias reman-cli='docker-compose run remote-manager bash'
alias reman-console='docker-compose run remote-manager ./ru'
```

## Contributing

Remote Manager is an open source project. If you find bugs or have proposal please create [issue](https://github.com/lazy-ants/remote-manager/issues) or Pull Request

## License

Copyright 2020 Lazy Ants

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
