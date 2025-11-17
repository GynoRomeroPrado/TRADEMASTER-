# Contributing to TRADEMASTER

Thank you for your interest in contributing to TRADEMASTER!

## Development Setup

### Prerequisites
- Node.js 18+
- Go 1.21+
- Python 3.11+
- Docker Desktop
- PostgreSQL 15+
- MongoDB 7+

### Initial Setup

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/trademaster.git
cd trademaster
```

2. **Install dependencies**
```bash
make install
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Start databases**
```bash
make docker-up
```

5. **Run migrations**
```bash
make migrate
```

6. **Start development servers**
```bash
make dev
```

## Project Structure

```
trademaster/
├── frontend/              # Next.js 14 frontend
├── mobile/                # React Native mobile app
├── backend/               # Go microservices
│   ├── auth-service/
│   ├── estimation-service/
│   ├── project-service/
│   ├── training-service/
│   └── marketplace-service/
├── ml-services/           # Python ML/AI services
├── ar-vr-training/        # Unity AR/VR modules
├── infrastructure/        # Infrastructure as code
└── docs/                  # Documentation
```

## Development Workflow

### Creating a Feature

1. Create a feature branch
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes
3. Write tests
4. Run tests and linters
```bash
make test
make lint
```

5. Commit your changes
```bash
git add .
git commit -m "feat: add your feature description"
```

6. Push to your fork
```bash
git push origin feature/your-feature-name
```

7. Create a Pull Request

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Examples:
```
feat: add blueprint analysis with YOLO
fix: resolve auth token expiration issue
docs: update API documentation
```

## Code Style

### TypeScript/JavaScript
- Use ESLint and Prettier
- Run `npm run format` before committing
- Follow Airbnb style guide

### Go
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `make lint` before committing

### Python
- Use Black for formatting
- Follow PEP 8
- Run `black .` before committing

## Testing

### Frontend
```bash
cd frontend && npm test
```

### Backend
```bash
cd backend && make test
```

### ML Services
```bash
cd ml-services && pytest
```

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add description of changes
4. Request review from maintainers
5. Address review comments
6. Merge after approval

## Code Review Guidelines

- Be respectful and constructive
- Focus on code quality and best practices
- Suggest improvements, don't just criticize
- Approve when satisfied with changes

## Reporting Bugs

1. Check if bug already reported
2. Create detailed issue with:
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots if applicable
   - Environment details

## Feature Requests

1. Check if feature already requested
2. Create issue with:
   - Use case description
   - Proposed solution
   - Alternative approaches
   - Mockups if applicable

## Questions?

- Create a discussion in GitHub Discussions
- Join our Discord community
- Email: dev@trademaster.io

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
