FROM gitpod/workspace-full

RUN sudo apt-get update && \
    sudo apt-get install -y gnome-keyring dbus-x11 build-essential ca-certificates && \
    sudo mkdir -p /github/home/.cache/ && \
    sudo mkdir -p /github/home/.local/share/keyrings/ && \
    sudo chmod 700 -R /github/home/.local/
