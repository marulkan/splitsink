#!/bin/env bash

mkdir -p ~/.config/systemd/user/
cp soundbox.service ~/.config/systemd/user
systemctl --user enable soundbox.service
