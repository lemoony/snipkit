FROM gitpod/workspace-postgres

RUN sudo apt-get update && \
    sudo apt-get install -y gnome-keyring dbus-x11 build-essential ca-certificates && \
    mkdir -p /github/home/.cache/ && \
    mkdir -p /github/home/.local/share/keyrings/ && \
    chmod 700 -R /github/home/.local/
