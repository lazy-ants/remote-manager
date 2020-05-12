# Remote Manager

This tool is intended for mass management and monitoring of remote servers.

The main idea is to get information about the status of remote servers, analyze it and provide maintenance as easily as possible.
 
The main goal of the project is also to create an utility that can be quickly extended for your own needs.

Feel free to send pull requests if you have any ideas for improving or extending the functionality.


## Setup

### 0. Preparation

You will need to setup docker in order to use this tool.

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

Open config.json with the editor of your choice and fill it with servers you want to manage.

Here:

name: the alias name of your server (for internal use)
connection-string: connection to you server user@domain:port (port is optional)
sudo-password: sudo password if you want to use commands needed sudo

In case your personal private key has a different name than id_rsa, got to the .env.local and provided it

### 4. Prepare servers you want to manage

We need a possibility to provide the sudo password in a secure way as an environment variable to your server (for sure in case you need the possibility of runnnig commands with sudo).
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

### 5. Add alias
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

Re-login into terminal or reload alias config so changes take effect.

### 6. Validate configuration and server accessibility
Run:
```
reman-console app:validate-config
```

#### Result

<img width="672" alt="image" src="https://user-images.githubusercontent.com/28564/81734095-ac503200-949b-11ea-90fc-1a7a7803aff5.png">

### 7. Run you first command to see e.g. the server uptime

To test this setup run
```
reman-console app:uptime
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

### app:reboot                  
Checks whether a reboot is required

### app:system-info             
Get system information

### app:upgrade                 
Upgrade server packages

### app:uptime                  
Get server uptime

### app:validate-config         
Validate server instances config

## Login into docker

If you want to login into the docker container:
```
reman-cli bash
```
