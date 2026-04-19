package custom_domain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	tenantRepo "invitation-api/internal/repository/tenant"
	"invitation-api/pkg/config"

	"github.com/google/uuid"
)

// Service handles custom domain operations
type Service struct {
	tenantRepo   tenantRepo.Repository
	config       *config.Config
	nginxConfDir string
}

// NewService creates a new custom domain service
func NewService(tenantRepo tenantRepo.Repository, cfg *config.Config) *Service {
	return &Service{
		tenantRepo:   tenantRepo,
		config:       cfg,
		nginxConfDir: "/etc/nginx/conf.d/custom_domains",
	}
}

// NginxServerTemplate is the template for custom domain nginx config
const NginxServerTemplate = `# Custom domain config for tenant: {{.TenantID}}
# Domain: {{.CustomDomain}}
# Generated automatically - DO NOT EDIT MANUALLY

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name {{.CustomDomain}} www.{{.CustomDomain}};

    ssl_certificate /etc/nginx/ssl/live/{{.CustomDomain}}/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/live/{{.CustomDomain}}/privkey.pem;

    add_header Strict-Transport-Security "max-age=31536000" always;

    # Static assets
    location /assets/ {
        alias /var/www/invitation/public/assets/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    location /uploads/ {
        alias /var/www/invitation/uploads/;
        expires 7d;
    }

    # htmx endpoints
    location ~ ^/(gallery|messages|gift)$ {
        rewrite ^/(.*)$ /{{.Subdomain}}/$1 break;
        
        proxy_pass http://backend_api;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Subdomain {{.Subdomain}};
        proxy_set_header X-Custom-Domain {{.CustomDomain}};
    }

    location /rsvp {
        rewrite ^/rsvp$ /{{.Subdomain}}/rsvp break;
        
        proxy_pass http://backend_api;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Subdomain {{.Subdomain}};
        
        limit_req zone=rsvp_limit burst=5 nodelay;
    }

    location / {
        rewrite ^/(.*)$ /{{.Subdomain}}/$1 break;
        
        proxy_pass http://backend_api;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Subdomain {{.Subdomain}};
        proxy_set_header X-Custom-Domain {{.CustomDomain}};
    }
}

# HTTP to HTTPS redirect
server {
    listen 80;
    listen [::]:80;
    server_name {{.CustomDomain}} www.{{.CustomDomain}};

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://{{.CustomDomain}}$request_uri;
    }
}
`

// DomainConfig holds data for nginx config template
type DomainConfig struct {
	TenantID     string
	CustomDomain string
	Subdomain    string
}

// SetupCustomDomain sets up a custom domain for a tenant
func (s *Service) SetupCustomDomain(tenantID uuid.UUID, customDomain string) error {
	// Get tenant
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// Validate domain
	if err := s.validateDomain(customDomain); err != nil {
		return err
	}

	// Update tenant with custom domain
	tenant.CustomDomain = &customDomain
	if err := s.tenantRepo.Update(tenant); err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	// Generate nginx config
	if err := s.generateNginxConfig(tenant.ID.String(), customDomain, tenant.Subdomain); err != nil {
		return fmt.Errorf("failed to generate nginx config: %w", err)
	}

	// Request SSL certificate
	if err := s.requestSSLCertificate(customDomain); err != nil {
		return fmt.Errorf("failed to request SSL certificate: %w", err)
	}

	// Reload nginx
	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// RemoveCustomDomain removes a custom domain from a tenant
func (s *Service) RemoveCustomDomain(tenantID uuid.UUID) error {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	if tenant.CustomDomain == nil {
		return nil // No custom domain to remove
	}

	customDomain := *tenant.CustomDomain

	// Remove nginx config
	configPath := filepath.Join(s.nginxConfDir, fmt.Sprintf("%s.conf", tenant.ID.String()))
	os.Remove(configPath)

	// Update tenant
	tenant.CustomDomain = nil
	if err := s.tenantRepo.Update(tenant); err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	// Reload nginx
	s.reloadNginx()

	_ = customDomain // Could revoke certificate here if needed
	return nil
}

// validateDomain performs basic domain validation
func (s *Service) validateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Check if domain is already in use
	// TODO: Query database to check if domain exists

	return nil
}

// generateNginxConfig generates nginx config file for custom domain
func (s *Service) generateNginxConfig(tenantID, customDomain, subdomain string) error {
	// Ensure directory exists
	os.MkdirAll(s.nginxConfDir, 0755)

	// Parse template
	tmpl, err := template.New("nginx").Parse(NginxServerTemplate)
	if err != nil {
		return err
	}

	// Create config file
	configPath := filepath.Join(s.nginxConfDir, fmt.Sprintf("%s.conf", tenantID))
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	data := DomainConfig{
		TenantID:     tenantID,
		CustomDomain: customDomain,
		Subdomain:    subdomain,
	}

	return tmpl.Execute(file, data)
}

// requestSSLCertificate requests SSL certificate for custom domain
func (s *Service) requestSSLCertificate(domain string) error {
	// Use certbot to get certificate
	cmd := exec.Command("certbot", "certonly",
		"--webroot",
		"--webroot-path=/var/www/certbot",
		"--email", "admin@yourdomain.com",
		"--agree-tos",
		"--no-eff-email",
		"-d", domain,
		"-d", "www."+domain,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("certbot failed: %s - %w", string(output), err)
	}

	return nil
}

// reloadNginx reloads nginx configuration
func (s *Service) reloadNginx() error {
	cmd := exec.Command("nginx", "-s", "reload")
	return cmd.Run()
}

// VerifyDomainDNS verifies that domain DNS is correctly configured
func (s *Service) VerifyDomainDNS(domain string) (*DNSVerificationResult, error) {
	result := &DNSVerificationResult{
		Domain: domain,
	}

	// Check A record
	cmd := exec.Command("dig", "+short", "A", domain)
	output, err := cmd.Output()
	if err == nil {
		result.ARecord = string(output)
	}

	// Check CNAME record
	cmd = exec.Command("dig", "+short", "CNAME", domain)
	output, err = cmd.Output()
	if err == nil {
		result.CNAMERecord = string(output)
	}

	// Determine if properly configured
	// In production, compare against your server IP
	result.IsValid = result.ARecord != "" || result.CNAMERecord != ""

	return result, nil
}

// DNSVerificationResult holds DNS check results
type DNSVerificationResult struct {
	Domain      string `json:"domain"`
	ARecord     string `json:"a_record,omitempty"`
	CNAMERecord string `json:"cname_record,omitempty"`
	IsValid     bool   `json:"is_valid"`
	Message     string `json:"message,omitempty"`
}
