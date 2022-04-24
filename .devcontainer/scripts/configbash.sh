#!/bin/bash
set -e
set -x

# Enable bash completion
echo "source /etc/bash_completion" >> "/root/.bashrc"

echo "alias ll='ls -la'" >> "/root/.bashrc"

# Git command prompt
git clone https://github.com/magicmonty/bash-git-prompt.git ~/.bash-git-prompt --depth=1 
echo "if [ -f \"$HOME/.bash-git-prompt/gitprompt.sh\" ]; then GIT_PROMPT_ONLY_IN_REPO=1 && source $HOME/.bash-git-prompt/gitprompt.sh; fi" >> "/root/.bashrc"

echo "export AKC_AUTH_BASE_ADDRESS=https://akc-duende.azurewebsites.net" >> "/root/.bashrc"
echo "export AKC_AUTH_BASE_PATH=/my" >> "/root/.bashrc"
echo "export AKC_AUTH_AUTHORIZATION_TYPE=client_credentials" >> "/root/.bashrc"
echo "export AKC_AUTH_CLIENT_ID=client" >> "/root/.bashrc"
echo "export AKC_AUTH_CLIENT_SECRET=secret" >> "/root/.bashrc"
echo "export AKC_AUTH_SCOPES=IdentityServerApi" >> "/root/.bashrc"
