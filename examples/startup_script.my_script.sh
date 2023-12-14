#!/bin/bash

# Update system packages
apt-get update
apt-get upgrade -y

# Install Apache web server
apt-get install -y apache2

# Enable and start Apache
systemctl enable apache2
systemctl start apache2
