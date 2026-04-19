#!/bin/bash
# Certificate Renewal Script

DOMAIN="yourdomain.com"

echo "Renewing SSL certificates..."

# Renew certificates
certbot renew --quiet

# Copy new certificates
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ./ssl/live/$DOMAIN/
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ./ssl/live/$DOMAIN/

# Reload nginx
docker-compose exec nginx nginx -s reload

echo "Certificates renewed at $(date)"
