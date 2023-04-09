#!/bin/bash

set -e

check_prerequisites() {
	for c in make git tmux; do
		if ! [[ -x "$(command -v ${c})" ]]; then
			echo "ERR: missing command: '${c}'." >&2
			echo "Please install before continue." >&2
			exit 1
		fi
	done
}

check_install_go() {
	# === Check Go installation
	if ! [[ -x "$(command -v go)" ]]; then
		cpu=`dpkg --print-architecture`
		version=${1:-"1.20.2"}

		# sudo rm -r /usr/local/go

		wget -c https://go.dev/dl/go$version.linux-$cpu.tar.gz -O - | \
			sudo tar -xz -C /usr/local

		ln -n /usr/local/go/bin/go /usr/bin/go
	fi
}

install_libssl() {
	wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.0g-2ubuntu4_amd64.deb
	sudo dpkg -i ./libssl1.1_1.1.0g-2ubuntu4_amd64.deb
	rm libssl1.1_1.1.0g-2ubuntu4_amd64.deb
}

install_indy() {
	INDY_LIB_VERSION="1.16.0"
	UBUNTU_VERSION="bionic"

	sudo apt-get update && \
	    sudo apt-get install -y software-properties-common apt-transport-https && \
	    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 68DB5E88 && \
	    sudo add-apt-repository -y "deb https://repo.sovrin.org/sdk/deb $UBUNTU_VERSION stable" && \
	    sudo add-apt-repository -y "deb https://repo.sovrin.org/sdk/deb xenial stable" && \
	    sudo apt-get update

	sudo apt-get install -y libindy-dev="$INDY_LIB_VERSION-xenial" \
	    libindy="$INDY_LIB_VERSION-$UBUNTU_VERSION"
}

clone_repos() {
	mkdir -p "$install_root"
	pushd "$install_root" > /dev/null

#	git clone https://github.com/findy-network/findy-wrapper-go.git "$install_root/findy-wrapper-go"
	git clone https://github.com/findy-network/findy-agent-auth.git "$install_root/findy-agent-auth"
	git clone https://github.com/findy-network/findy-agent.git "$install_root/findy-agent"
	git clone https://github.com/findy-network/findy-agent-cli.git "$install_root/findy-agent-cli"
	popd > /dev/null
}

build_agency() {
	pushd "$install_root/findy-agent" > /dev/null
	# git checkout <branch>
	make cli
	popd > /dev/null
}

build_cli() {
	pushd "$install_root/findy-agent-cli" > /dev/null
	# git checkout <branch>
	make cli
	popd > /dev/null
}

build_auth() {
	pushd "$install_root/findy-agent-auth" > /dev/null
	git checkout start-scripts
	make acli
	popd > /dev/null
}

install_tmux_conf() {
	cat >"$HOME/.tmux.conf" <<EOF
set -g default-terminal "xterm"

# from neovim info
set-option -sg escape-time 10
set-option -g focus-events on
# set -g default-terminal "screen-256color"
# set -sa terminal-overrides ",alacritty:RGB"

# remap prefix from 'C-b' to 'C-a'
unbind C-b
set-option -g prefix C-a
bind-key C-a send-prefix

# split panes using | and -
bind | split-window -h
bind - split-window -v
unbind '"'
unbind %

# Enable mouse mode (tmux 2.1 and above)
set -g mouse on

# force SHELL ENV variable as shell
set-option -g default-shell ${SHELL}

# Smart pane switching with awareness of Vim splits.
# See: https://github.com/christoomey/vim-tmux-navigator
is_vim="ps -o state= -o comm= -t '#{pane_tty}' \
    | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?(view|n?vim?x?)(diff)?$'"
bind-key -n 'C-h' if-shell "$is_vim" 'send-keys C-h'  'select-pane -L'
bind-key -n 'C-j' if-shell "$is_vim" 'send-keys C-j'  'select-pane -D'
bind-key -n 'C-k' if-shell "$is_vim" 'send-keys C-k'  'select-pane -U'
bind-key -n 'C-l' if-shell "$is_vim" 'send-keys C-l'  'select-pane -R'
tmux_version='$(tmux -V | sed -En "s/^tmux ([0-9]+(.[0-9]+)?).*/\1/p")'
if-shell -b '[ "$(echo "$tmux_version < 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\'  'select-pane -l'"
if-shell -b '[ "$(echo "$tmux_version >= 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\\\'  'select-pane -l'"

bind-key -T copy-mode-vi 'C-h' select-pane -L
bind-key -T copy-mode-vi 'C-j' select-pane -D
bind-key -T copy-mode-vi 'C-k' select-pane -U
bind-key -T copy-mode-vi 'C-l' select-pane -R
bind-key -T copy-mode-vi 'C-\' select-pane -l

bind -r H resize-pane -L 5
bind -r J resize-pane -D 5
bind -r K resize-pane -U 5
bind -r L resize-pane -R 5

bind C-x setw synchronize-panes
bind C-e command-prompt -p "Command:" \
         "run \"tmux list-panes -s -F '##{session_name}:##{window_index}.##{pane_index}' \
                | xargs -I PANE tmux send-keys -t PANE '%1' Enter\""

bind-key r source-file ~/.tmux.conf \; display-message "~/.tmux.conf reloaded"
bind-key M split-window -h "nvim ~/.tmux.conf"

bind X confirm-before kill-session
EOF
}

install_tmuxinator_play() {
	mkdir "$HOME/.config/tmuxinator"
	sudo apt install -y tmuxinator

cat > "$HOME/.config/tmuxinator/play.yml" <<EOF
name: play
root: ~/go/src/github.com/findy-network

windows:
  - editor:
      layout: main-vertical
      panes:
        - # empty shell
        - # empty shell
  - running:
      layout: tiled
      panes:
        - cd findy-agent-auth/scripts; ./mem-dev-run.sh
        - agency:
          - cd findy-agent/scripts/test
          - fa ledger steward create --config create-steward-to-mem-ledger.yaml
          - agency=fa register=findy.json no_clean=1 enclave=MEMORY_enclave.bolt ./mem-server --reset-register --grpc-cert-path ../../grpc/cert
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
        - cd findy-agent-cli/scripts/fullstack; source ./setup-cli-env-local.sh
EOF
}

CURRENT_DIR=$(dirname "$BASH_SOURCE")

check_prerequisites
check_install_go

GOPATH=${GOPATH:-`go env GOPATH`}
gopath=${GOPATH:-"$PWD"}
install_root="$gopath/src/github.com/findy-network"
alias pf='printenv|grep FCLI'

install_libssl
install_indy
clone_repos
build_agency
build_cli
build_auth
install_tmux_conf
install_tmuxinator_play

# steward creation: DONE
# don't use mem/ -dir for findy.json or create it: DONE
# pf alias: DONE 
# -local.sh GOPATH problem. and mem-server auth, export it!: DONE
# auth-server-grpc path problem: needs PR or branch usage.
