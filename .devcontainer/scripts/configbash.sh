#!/bin/bash
set -e
set -x

# Enable bash completion
echo "source /etc/bash_completion" >> "/root/.bashrc"

echo "alias ll='ls -la'" >> "/root/.bashrc"

# Git command prompt
git clone https://github.com/magicmonty/bash-git-prompt.git ~/.bash-git-prompt --depth=1 
echo "if [ -f \"$HOME/.bash-git-prompt/gitprompt.sh\" ]; then GIT_PROMPT_ONLY_IN_REPO=1 && source $HOME/.bash-git-prompt/gitprompt.sh; fi" >> "/root/.bashrc"

wget https://releases.hashicorp.com/terraform/1.1.9/terraform_1.1.9_linux_amd64.zip
unzip terraform_1.1.9_linux_amd64.zip
mv terraform /usr/local/bin/terraform
rm terraform_1.1.9_linux_amd64.zip