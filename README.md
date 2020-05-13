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

Open **config.json** with the editor of your choice and and add server connections.
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

<img width="672" alt="image" src="https://user-images.githubusercontent.com/28564/81734095-ac503200-949b-11ea-90fc-1a7a7803aff5.png">

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

### app:docker-prune
Prune old docker data

### app:docker-ps               
Show docker process status

### app:kernel                  
Get server kernels

### app:ls                      
Run ls command on all servers

### app:os                      
Get server OS

### app:check-reboot                  
Checks whether a reboot is required

### app:system-info
- note: [sudo required]
- Get system information

### app:upgrade
- note: [sudo required]
- Upgrade server packages

### app:uptime                  
Get server uptime

### app:validate-config         
Validate server instances config

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
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager bin/console'
```

Otherwise you can simple map the whole .ssh directory:
```
alias reman-cli='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager'
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager bin/console'
```

Or if you prefer docker-compose:
```
Otherwise you can simple map the whole .ssh directory:
```
alias reman-cli='docker-compose run remote-manager bash'
alias reman-console='docker-compose run remote-manager bin/console'
```
