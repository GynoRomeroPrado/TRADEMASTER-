/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    API_URL: process.env.API_URL || 'http://localhost:8080',
    ML_API_URL: process.env.ML_API_URL || 'http://localhost:8000',
  },
  images: {
    domains: ['localhost', 's3.amazonaws.com'],
  },
}

module.exports = nextConfig
