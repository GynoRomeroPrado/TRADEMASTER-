# TRADEMASTER Frontend

Next.js 14 frontend for TRADEMASTER platform.

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Data Fetching**: TanStack Query (React Query)
- **Forms**: React Hook Form + Zod
- **Icons**: Lucide React
- **Charts**: Recharts

## Getting Started

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
src/
├── app/                    # Next.js app router
│   ├── (auth)/            # Auth pages (login, register)
│   ├── dashboard/         # Dashboard pages
│   ├── projects/          # Project pages
│   ├── estimates/         # Estimation pages
│   ├── training/          # Training pages
│   ├── marketplace/       # Marketplace pages
│   └── layout.tsx         # Root layout
├── components/            # React components
│   ├── ui/               # Base UI components
│   ├── forms/            # Form components
│   └── layouts/          # Layout components
├── lib/                   # Utilities
│   ├── api.ts            # API client
│   ├── auth.ts           # Auth helpers
│   └── utils.ts          # Utility functions
└── types/                 # TypeScript types
```

## Features

### Implemented
- Landing page
- Responsive header & footer
- Feature showcase
- Stats section
- API client setup
- React Query provider

### To Implement
- Authentication flow
- Dashboard
- Project management UI
- Estimation workflow
- Training interface
- Marketplace
- AR/VR components
- Real-time updates (Socket.io)

## Environment Variables

Create a `.env.local` file:

```env
API_URL=http://localhost:8080
ML_API_URL=http://localhost:8000
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

## Development

```bash
# Run with hot reload
npm run dev

# Type checking
npx tsc --noEmit

# Linting
npm run lint

# Format code
npm run format
```

## Deployment

The app can be deployed to Vercel:

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel
```

## Contributing

1. Create feature branch
2. Make changes
3. Test thoroughly
4. Submit PR

## License

MIT
