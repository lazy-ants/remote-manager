# Remote Manager

## Setup

### 0. Clone project

```bash
git clone git@github.com:lazy-ants/remote-manager.git
cd remote-manager
```

### 1. Run
```bash
docker build -t remote-manager .
docker run -it --rm -v "$PWD":/usr/src/remote-manager remote-manager composer install
```

### 2. Create a config.

Copy config.example.json to config.json 
```
cp config.example.json to config.json
```
and fill it with servers you want to manage.
"sudo-password" is necessary in case you have to run commands with sudo and don't want to enter the sudo password each time
Alternatively if most of your servers use the same sudo password you can set in the .env.local.
In first turn the password from the config.json will be used, otherwise the one from the .env.local if it is provided 

### 3. Setup .env.local
Copy .env to .env.local
```
cp .env .env.local
```

Enter here comma separated names of your private keys in PPK_NAMES (if you have just one - enter only this)

### 4. Prepare servers you want to manage.
On each server you want to manage edit 
```
sudo nano /etc/ssh/sshd_config
``` 
and add at the end of the config:
```
AcceptEnv PASSWORD
```
and reload the sshd server e.g.
```
sudo service sshd reload
```

### 5. Add alias.
You can map each of your private and public keys and the known_hosts file:
```
alias reman-cli='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager'
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh/id_rsa:/root/.ssh/id_rsa -v ~/.ssh/id_rsa.pub:/root/.ssh/id_rsa.pub -v ~/.ssh/id_rsa:/root/.ssh/known_hosts remote-manager bin/console'
```

Or you can just map the whole .ssh directory, however it won't work on MacOS if you have some MacOS specific options like "UseKeychain yes" in the .ssh/config: 
```
alias reman-cli='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager'
alias reman-console='docker run -it --rm -v "$PWD":/usr/src/remote-manager -v ~/.ssh:/root/.ssh remote-manager bin/console'
```

Re-login into terminal or reload alias config so changes take effect

### 6. Running test command.

To test this setup run
```
reman-console app:ls
```

If you want to login into the docker container:
```
reman-cli bash
```
