#!/bin/env bash

# https://docs.spotifyd.rs/installation/Raspberry-Pi.html
mkdir -p ~/.config/systemd/user/
cp spotifyd.service ~/.config/systemd/user
sudo loginctl enable-linger ${SUDO_USER}
systemctl --user enable spotifyd.service
