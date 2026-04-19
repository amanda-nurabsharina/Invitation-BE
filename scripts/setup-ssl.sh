#!/bin/bash
# SSL Certificate Setup Script using Let's Encrypt
# Run this script on your production server

DOMAIN="yourdomain.com"
EMAIL="admin@yourdomain.com"

# Create directories
mkdir -p /var/www/certbot
mkdir -p /etc/nginx/ssl

# Install certbot if not installed
if ! command -v certbot &> /dev/null; then
    echo "Installing certbot..."
    apt-get update
    apt-get install -y certbot
fi

# Stop nginx temporarily
docker-compose stop nginx

# Get wildcard certificate (requires DNS validation)
echo "Getting wildcard SSL certificate for *.$DOMAIN"
certbot certonly \
    --manual \
    --preferred-challenges dns \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN" \
    -d "*.$DOMAIN"

# Alternative: Get certificate with HTTP validation (for single domains)
# certbot certonly \
#     --webroot \
#     --webroot-path=/var/www/certbot \
#     --email $EMAIL \
#     --agree-tos \
#     --no-eff-email \
#     -d "$DOMAIN" \
#     -d "api.$DOMAIN" \
#     -d "admin.$DOMAIN"

# Copy certificates to nginx ssl directory
mkdir -p ./ssl/live/$DOMAIN
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ./ssl/live/$DOMAIN/
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ./ssl/live/$DOMAIN/

# Set permissions
chmod 600 ./ssl/live/$DOMAIN/privkey.pem

# Restart nginx
docker-compose up -d nginx

echo "SSL certificate installed successfully!"
echo ""
echo "To auto-renew certificates, add this cron job:"
echo "0 12 * * * /path/to/renew-certs.sh >> /var/log/certbot-renew.log 2>&1"
