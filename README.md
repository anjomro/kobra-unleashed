# Kobra Unleashed

Since the Kobra 2 Pro/Plus/Max do not have a webinterface built-in, I decided to create one myself. This is the result.

These printers sadly feature a closed source firmware, so it isn't easily possible to extend the existing firmware.
Using another firmware is theoretically possible but difficult, since a lot of modifications would be needed to make it
work.


![](https://raw.githubusercontent.com/anjomro/kobra-unleashed/master/img/kobra-unleashed-idle.png)

## How does it work?

This webinterface uses the interface of the firmware that is designed to be used with the proprietary app of the
manufacturer. This interface is not documented and not officially supported. All controls and information used in the
webinterface are the result of reverse engineering.

## Prerequisites

- Root shell access to your Printer
  - Refer to [the guide in this Readme](https://github.com/ultimateshadsform/Anycubic-Kobra-2-Series-Firmware)
- A Linux Server
  - Reachable by IPv4 from the printer
  - Port 8883 has to be opened for MQTT(S)
  - Another port of your choice for the Webinterface
  - (Maybe it's possible to run this locally on something like a Raspberry Pi, but I haven't tested that yet)

## Setup the printer
- Clone this repository
  - `git clone https://github.com/anjomro/kobra-unleashed.git`
- Generate a self signed ca and certificates for the printer and MQTT Server
  - Use the script `certs/generate_certs.sh` for that
  - You can just skip with enter through the fields about the certificate details, they are not checked
- Backup the existing certificates on the printer
  - They are located in `/user/`
  - You can use `cp -r /user /root/user_backup` to save the previous certificates
- Replace the certificates on the printer with the generated ones
  - You need to copy the following files:
    - `certs/ca.pem` to `/user/ca.crt`
    - `certs/client.pem` to `/user/client.crt`
    - `certs/client.key` to `/user/client.key`
- Backup the app on the printer
  - It is located in `/app/app`
  - You can use `cp -r /app/app /root/app_backup` to save the previous app
- Replace the address the printer connects to with one you control
  - My printer connects to `mqtt-universe.anycubic.com`
    - This may change depending on the configured region?
    - Other addresses found in the printer are the following:
      - `mqtt-universe-test.anycubic.com`
      - `mqtt-test.anycubic.com`
      - `mqtt.anycubic.com`
  - Select a hostname that has exactly the same length as the original one
    - You could use [sslip.io](https://sslip.io/) to get a name for your IP
    - Attention: Your router might have `rebind protection` which prevents it from resolving external hostnames to local IPs. 
    - Example:
      - `mqtt-universe.anycubic.com` - 26 characters
      - `paddi.192.168.0.1.sslip.io` - 26 characters
  - Replace the address in the app binary & replace stock app
    - `< /app/app sed 's/mqtt-universe.anycubic.com/paddi.192.168.0.1.sslip.io/g' > /app/mod_app`
    - `chmod 755 /app/mod_app`
    - `mv /app/mod_app /app/app`
- Reboot the printer for the changes to take effect
  - `reboot`

## Setup the server
- Make sure you have Docker installed.
  - See [the official guide](https://docs.docker.com/desktop/install/linux-install/) for that.
- Copy the repo and the certs to your server
- Configure your address in the `docker-compose.yml`
  - `CORS_HOST=http://your-host:5000` or whatever you want to use
- Start the containers
  - `docker-compose up -d`
- Using the default configuration, the webinterface should be reachable at `http://<your-server-ip>:5000`
  - If you can't access it, make sure the port is open in your firewall
  - Also make sure port 8883 is open in order for the printer to connect to the MQTT Server
  - You can test MQTT Connectivity (and read the communication with your printer) using [MQTTX](https://mqttx.app/) 
- You may need to restart your printer again after your MQTT Server is running to establish the connection
- If everything went according to plan, the printer should now appear in the dropdown of the webinterface.