# TRADEMASTER Mobile App

React Native mobile application for TRADEMASTER platform.

## Features

- Field data collection
- Progress photo capture with geolocation
- AR training modules
- Offline-first functionality
- Real-time sync

## Tech Stack

- React Native + Expo
- TypeScript
- ARCore / ARKit for AR features
- React Native Maps
- Expo Camera
- AsyncStorage for offline

## Getting Started

```bash
# Install dependencies
npm install

# Start development
npm start

# Run on iOS
npm run ios

# Run on Android
npm run android
```

## Project Structure

```
mobile/
├── app/           # Expo Router pages
├── components/    # React components
├── services/      # API clients
├── hooks/         # Custom hooks
└── utils/         # Utilities
```

## Environment Variables

Create `.env` file:

```
API_URL=https://api.trademaster.io
```
