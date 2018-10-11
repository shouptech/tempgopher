#!/bin/bash

INSTALLDIR=/opt/tempgopher
INSTALLBIN=$INSTALLDIR/tempgopher
INSTALLUSER=pi
BINURL='https://gitlab.com/shouptech/tempgopher/-/jobs/artifacts/master/raw/tempgopher?job=build'
CONFIGFILE=$INSTALLDIR/config.yml

# Download binary
sudo mkdir -p $INSTALLDIR
sudo curl -L $BINURL -o $INSTALLBIN
sudo chmod +x $INSTALLBIN
sudo chown -R $INSTALLUSER: $INSTALLDIR

# Create unit file
sudo sh -c "cat > /etc/systemd/system/tempgopher.service" << EOM
[Unit]
Description=Temp Gopher
After=network.target

[Service]
Type=simple
WorkingDirectory=$INSTALLDIR
ExecStart=$INSTALLBIN -c $CONFIGFILE run
ExecReload=/bin/kill -HUP \$MAINPID
User=$INSTALLUSER
Group=$INSTALLUSER

[Install]
WantedBy=multi-user.target
EOM

sudo systemctl daemon-reload
sudo systemctl enable tempgopher.service
sudo systemctl start tempgopher.service
