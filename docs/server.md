# Server

## Firewall
```cmd
ufw --force reset
ufw allow from X.X.X.X to any port
ufw default deny incoming
ufw default allow outgoing
ufw enable
ufw status verbose
```

- Adjust local VPS firewall settings (remove all rules, add rules as needed)

## Cloudflared
```cmd
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared
chmod +x cloudflared
sudo mv cloudflared /usr/local/bin/
```

```cmd
cloudflared tunnel login
```

- Copy over existing `/etc/cloudflared/{credential}.json`

- Copy over existing `/etc/cloudflared/config.yml`

Verify
```cmd
cloudflared tunnel run albionstats
```

```cmd
cloudflared service install
systemctl enable cloudflared
systemctl start cloudflared
systemctl status cloudflared
```

## SSH Key

- Copy over existing `~/.ssh/authorized_keys`