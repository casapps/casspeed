# Admin Panel

The admin panel provides web-based configuration and management for casspeed.

## Access

Access the admin panel at: `http://{your-server}:{port}/admin`

## First-Run Setup

On first run, you'll be prompted to create the primary admin account:

1. Navigate to `/admin/setup`
2. Create admin username and password
3. Configure initial settings
4. Complete setup wizard

## Admin Features

### Server Configuration

Manage all server settings through the web UI:

- Server address and port
- Application mode (production/development)
- Rate limiting rules
- SSL/TLS certificates

### User Management

- View registered users
- Manage user accounts
- Review device registrations
- Monitor API tokens

### Speed Test Analytics

- View test history
- Analyze speed trends
- Export test data
- Manage shared results

### System Monitoring

- Server health status
- Resource usage
- Active connections
- Error logs

### Backup & Restore

- Configure automatic backups
- Manual backup creation
- Restore from backup
- Backup encryption settings

## Security

### Authentication

Admin panel uses session-based authentication:

- Secure session cookies
- Configurable session timeout
- Remember me option (30 days)

### 2FA (Optional)

Enable two-factor authentication:

- TOTP (Google Authenticator, Authy)
- WebAuthn/Passkeys
- Backup codes

### API Access

Admin API available at `/api/v1/admin/`:

- Bearer token authentication
- Create API tokens in admin panel
- Token format: `adm_xxxxxxxxxxxxxxxxxxxxxxxxxx`

## Configuration

All settings in `server.yml` can be edited through the admin UI.

Changes take effect immediately (except port/address which require restart).

## Audit Log

All admin actions are logged:

- Configuration changes
- User modifications
- Security events
- System operations

Access audit log in Admin Panel → Logs → Audit.
