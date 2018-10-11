#!/bin/bash

INSTALLDIR=/opt/tempgopher
INSTALLBIN=$INSTALLDIR/tempgopher
INSTALLUSER=pi
BINURL='https://gitlab.com/shouptech/tempgopher/-/jobs/artifacts/master/raw/tempgopher?job=build'
CONFIGFILE=$INSTALLDIR/config.yml

# Load w1_therm module
sudo /sbin/modprobe w1_therm

# Download binary
sudo mkdir -p $INSTALLDIR
sudo curl -L $BINURL -o $INSTALLBIN
sudo chmod +x $INSTALLBIN
sudo chown -R $INSTALLUSER: $INSTALLDIR

# Generate a configuration file
sudo -u $INSTALLUSER $INSTALLBIN -c $CONFIGFILE config

# Create unit file
sudo sh -c "cat > /etc/systemd/system/tempgopher.service" << EOM
[Unit]
Description=Temp Gopher
After=network.target

[Service]
Type=simple
WorkingDirectory=$INSTALLDIR
PermissionsStartOnly=true
User=$INSTALLUSER
Group=$INSTALLUSER
ExecStartPre=/sbin/modprobe w1_therm
ExecStart=$INSTALLBIN -c $CONFIGFILE run
ExecReload=/bin/kill -HUP \$MAINPID

[Install]
WantedBy=multi-user.target
EOM

# Enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable tempgopher.service
sudo systemctl start tempgopher.service
