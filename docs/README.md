# Axinova MCP Server Documentation

Documentation for the Model Context Protocol (MCP) server implementation.

## Quick Start
- [README](../README.md) - Project overview and quick start guide
- [Deployment Guide](runbooks/DEPLOYMENT.md) - Production deployment steps
- [Testing Guide](runbooks/TESTING.md) - How to test the MCP server
- [Token Generation Walkthrough](ops/TOKEN_GENERATION_WALKTHROUGH.md) - Authentication setup

---

## Operational Guides

Located in [ops/](ops/)

### Configuration & Validation
- [Validation Procedures](ops/VALIDATION.md) - Configuration and service validation
- [Token Generation Walkthrough](ops/TOKEN_GENERATION_WALKTHROUGH.md) - JWT token setup

### Implementation Status
- [Deployment Summary](ops/DEPLOYMENT_SUMMARY.md) - Initial deployment details
- [Deployment Complete Summary](ops/DEPLOYMENT_COMPLETE_SUMMARY.md) - Final deployment status
- [Final Deployment Complete](ops/FINAL_DEPLOYMENT_COMPLETE.md) - Completion report
- [Final Setup Steps](ops/FINAL_SETUP_STEPS.md) - Final configuration steps
- [Final Status and Next Steps](ops/FINAL_STATUS_AND_NEXT_STEPS.md) - Current state and roadmap
- [Status Update](ops/STATUS_UPDATE.md) - Latest status information

### Infrastructure
- [Infrastructure Analysis](ops/INFRASTRUCTURE_ANALYSIS.md) - Infrastructure review and recommendations
- [Security Group Status](ops/SECURITY_GROUP_STATUS.md) - Network security configuration
- [Portainer Agent Setup](ops/PORTAINER_AGENT_SETUP.md) - Portainer integration
- [Dashboard Fixes Final](ops/DASHBOARD_FIXES_FINAL.md) - Dashboard configuration fixes

### Issue Tracking
- [Issues Fixed and Remaining](ops/ISSUES_FIXED_AND_REMAINING.md) - Known issues and resolutions

---

## Runbooks

Located in [runbooks/](runbooks/)

- [Deployment Guide](runbooks/DEPLOYMENT.md) - Step-by-step deployment procedures
- [Testing Guide](runbooks/TESTING.md) - How to test MCP endpoints and functionality

---

## Development

### For AI Agents
- [CLAUDE.md](CLAUDE.md) - AI agent development guide and context
- [AGENTS.md](../AGENTS.md) - Repository guidelines for AI agents (root level)

### Project Structure
```
axinova-mcp-server-go/
├── cmd/                    # Application entrypoints
├── internal/               # Internal packages
├── config/                 # Configuration files
├── scripts/                # Utility scripts
├── docs/                   # This documentation
│   ├── ops/               # Operational guides and status
│   ├── runbooks/          # Step-by-step procedures
│   └── adr/               # Architectural Decision Records
├── Dockerfile             # Container build
├── docker-compose.yml     # Local development stack
├── Makefile              # Build automation
└── README.md             # Main project README
```

---

## Architectural Decision Records

Located in [adr/](adr/)

*No ADRs yet. Use this directory for future architectural decisions.*

**When to create an ADR:**
- Choosing between different MCP protocol versions
- Selecting authentication/authorization strategies
- Deciding on database schema changes
- Major refactoring decisions
- Infrastructure architecture changes

**ADR Template:**
```markdown
# ADR-001: Title

**Status:** Proposed | Accepted | Deprecated | Superseded  
**Date:** YYYY-MM-DD  
**Deciders:** Team members involved

## Context
What is the issue we're facing? What constraints exist?

## Decision
What did we decide to do?

## Consequences
What becomes easier or harder as a result of this decision?

## Alternatives Considered
What other options did we evaluate and why were they rejected?
```

---

## MCP Server Architecture

### Overview
The Axinova MCP Server provides a Model Context Protocol interface for AI agents to interact with the Axinova platform. It exposes tools and resources through a standardized protocol.

### Key Components
1. **MCP Protocol Handler** - Implements MCP specification
2. **Tool Registry** - Available tools for AI agents
3. **Resource Provider** - Exposes platform resources
4. **Authentication** - JWT-based token validation
5. **Logging & Monitoring** - Structured logging and metrics

### Deployment Architecture
- Runs as Docker container
- Exposed through Traefik reverse proxy
- Integrated with Prometheus/Grafana monitoring
- Connected to platform PostgreSQL database

---

## Common Tasks

| Task | Document | Section |
|------|----------|---------|
| Deploy MCP server | [Deployment Guide](runbooks/DEPLOYMENT.md) | Deployment Steps |
| Test MCP endpoints | [Testing Guide](runbooks/TESTING.md) | Testing Procedures |
| Generate auth tokens | [Token Generation](ops/TOKEN_GENERATION_WALKTHROUGH.md) | Token Setup |
| Validate config | [Validation](ops/VALIDATION.md) | Configuration Checks |
| Check infrastructure | [Infrastructure Analysis](ops/INFRASTRUCTURE_ANALYSIS.md) | Review |
| Setup monitoring | [Dashboard Fixes](ops/DASHBOARD_FIXES_FINAL.md) | Grafana Setup |

---

## API Reference

### MCP Endpoints

**Base URL:** `https://mcp.axinova.ai`

#### Health Check
```bash
curl https://mcp.axinova.ai/health
```

#### Tool Discovery
```bash
curl -H "Authorization: Bearer <token>" https://mcp.axinova.ai/tools
```

#### Resource Listing
```bash
curl -H "Authorization: Bearer <token>" https://mcp.axinova.ai/resources
```

See [Testing Guide](runbooks/TESTING.md) for detailed API examples.

---

## Development Workflow

### Local Development
1. Clone repository
2. Copy `.env.example` to `.env`
3. Run `docker-compose up`
4. Access server at `http://localhost:8080`

### Testing
```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Test MCP endpoints
./test_mcp.sh
./test_tool_call.sh
```

### Building
```bash
# Build binary
make build

# Build Docker image
make docker-build

# Build and run locally
make run
```

---

## Troubleshooting

### Common Issues

#### 1. Connection Refused
**Symptom:** Cannot connect to MCP server  
**Solution:** Check Traefik routing configuration and service health

#### 2. Authentication Failed
**Symptom:** 401 Unauthorized responses  
**Solution:** Verify JWT token generation and expiration - see [Token Generation](ops/TOKEN_GENERATION_WALKTHROUGH.md)

#### 3. Database Connection Errors
**Symptom:** Cannot connect to PostgreSQL  
**Solution:** Check database URL and credentials in configuration

#### 4. Tool Execution Failures
**Symptom:** Tools return errors or timeouts  
**Solution:** Check tool implementation logs and resource availability

For more troubleshooting, see [Issues Fixed and Remaining](ops/ISSUES_FIXED_AND_REMAINING.md)

---

## Security Considerations

### Authentication
- All endpoints except `/health` require JWT authentication
- Tokens expire after configured duration
- See [Token Generation](ops/TOKEN_GENERATION_WALKTHROUGH.md) for setup

### Network Security
- Server runs behind Traefik reverse proxy
- HTTPS/TLS termination at proxy
- Security group configuration in [Security Group Status](ops/SECURITY_GROUP_STATUS.md)

### Secrets Management
- Environment variables for sensitive data
- No secrets committed to repository
- Use `.env` file locally (never commit)

---

## Monitoring & Observability

### Metrics
- Prometheus metrics exposed on `/metrics`
- Grafana dashboards for visualization
- See [Dashboard Fixes](ops/DASHBOARD_FIXES_FINAL.md)

### Logging
- Structured JSON logging
- Log aggregation through Loki (if configured)
- Log levels configurable via environment

### Health Checks
- HTTP health endpoint: `/health`
- Readiness probe: `/ready`
- Liveness probe: `/live`

---

## Contributing

### Documentation Standards
1. Use Markdown for all documentation
2. Follow existing structure and organization
3. Update this README when adding new docs
4. Keep operational logs in `docs/ops/`
5. Keep procedures in `docs/runbooks/`
6. Create ADRs for significant decisions

### Code Standards
- Follow Go best practices
- Run `make fmt` before committing
- Write tests for new features
- Update documentation with code changes

---

## Related Projects

- [axinova-deploy](../../axinova-deploy/) - Deployment automation
- [axinova-ai-lab-go](../../axinova-ai-lab-go/) - AI Lab backend
- [axinova-home-go](../../axinova-home-go/) - Platform backend

---

## External References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)
- [Traefik Documentation](https://doc.traefik.io/traefik/)

---

**Last Updated:** January 20, 2026  
**Maintained By:** MCP Server Team  
**Questions?** See [CLAUDE.md](CLAUDE.md) for AI agent context
