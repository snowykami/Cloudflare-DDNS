Language [zh_CN](README.md) | [EN](README_EN.md)

## Cloudflare DDNS

Cloudflare-DDNS is a program for dynamic domain name resolution that updates domain records using Cloudflare's API.

### Prerequisites

1. Ensure that you have a public IP address.

### How to Use

1. Add the domain to your Cloudflare account.
2. Download the program for your device's corresponding system and architecture from [Releases], then run the program.
3. The first run will automatically generate `config.yml`. Modify it and run the program again.

### Configuration File Explanation

```yaml
api_key: xxxxxxxx # The API key obtained from (https://dash.cloudflare.com/profile/api-tokens)
api_email: user@example.com # Your Cloudflare account email
zone_id: xxxxxxxx # Zone ID obtained from Cloudflare domain homepage
duration: 5 # Update interval in seconds
ddns: # Your domain records. You can add multiple records here, and a subdomain can have both A and AAAA records.
  - name: v4.example.com
    type: A
    ttl: 60
    proxied: false
    comment: "DDNS auto update"

  - name: v6.example.com
    type: AAAA
    ttl: 60
    proxied: false
    comment: "DDNS auto update"
```
### Some Recommendations
1. Considering that most DDNS users have home broadband with public IPs, v4, and v6 addresses may not necessarily be on the same host. Therefore, it is recommended to configure port forwarding locally, and use different records for A and AAAA.
2. You can set up the program as a startup service, making it convenient for automatic domain resolution updates after server power-off and restart.
3. If your host is in a NAT environment, consider using UPnP or NAT-PMP protocols for port mapping.
4. If your network device does not support IPv6, avoid configuring AAAA records.
