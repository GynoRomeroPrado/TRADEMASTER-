import Link from 'next/link'
import { ArrowRight, Zap, Brain, GraduationCap, TrendingUp } from 'lucide-react'

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-blue-50 to-white">
      {/* Header */}
      <header className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <Zap className="h-8 w-8 text-primary-600" />
            <h1 className="text-2xl font-bold text-gray-900">TRADEMASTER</h1>
          </div>
          <nav className="hidden md:flex gap-6">
            <Link href="/features" className="text-gray-600 hover:text-primary-600">
              Features
            </Link>
            <Link href="/pricing" className="text-gray-600 hover:text-primary-600">
              Pricing
            </Link>
            <Link href="/training" className="text-gray-600 hover:text-primary-600">
              Training
            </Link>
            <Link href="/about" className="text-gray-600 hover:text-primary-600">
              About
            </Link>
          </nav>
          <div className="flex gap-4">
            <Link
              href="/login"
              className="px-4 py-2 text-gray-700 hover:text-primary-600"
            >
              Login
            </Link>
            <Link
              href="/signup"
              className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
            >
              Get Started
            </Link>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-5xl font-bold text-gray-900 mb-6">
            The All-in-One Platform for
            <span className="text-primary-600"> Skilled Trades</span>
          </h2>
          <p className="text-xl text-gray-600 mb-8">
            AI-powered project estimation, management, and training for electricians,
            HVAC technicians, and welders. Built for the future of construction.
          </p>
          <div className="flex gap-4 justify-center">
            <Link
              href="/signup"
              className="px-8 py-4 bg-primary-600 text-white rounded-lg hover:bg-primary-700 flex items-center gap-2 text-lg font-semibold"
            >
              Start Free Trial
              <ArrowRight className="h-5 w-5" />
            </Link>
            <Link
              href="/demo"
              className="px-8 py-4 border-2 border-primary-600 text-primary-600 rounded-lg hover:bg-primary-50 text-lg font-semibold"
            >
              Watch Demo
            </Link>
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section className="container mx-auto px-4 py-20">
        <h3 className="text-3xl font-bold text-center text-gray-900 mb-12">
          Everything you need to run your trade business
        </h3>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          <FeatureCard
            icon={<Brain className="h-12 w-12 text-primary-600" />}
            title="AI Estimation"
            description="Automated takeoffs from blueprints with 90%+ accuracy. Get instant cost predictions."
          />
          <FeatureCard
            icon={<TrendingUp className="h-12 w-12 text-primary-600" />}
            title="Project Management"
            description="Track progress, manage teams, and optimize routes. All in one place."
          />
          <FeatureCard
            icon={<GraduationCap className="h-12 w-12 text-primary-600" />}
            title="AI Training"
            description="AR/VR training modules with personalized AI coaching. Get certified faster."
          />
          <FeatureCard
            icon={<Zap className="h-12 w-12 text-primary-600" />}
            title="Job Marketplace"
            description="Find projects matched to your skills. Build your reputation with blockchain verification."
          />
        </div>
      </section>

      {/* Stats Section */}
      <section className="bg-primary-600 text-white py-20">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8 text-center">
            <StatCard number="90%+" label="Estimation Accuracy" />
            <StatCard number="10hrs" label="Saved per Week" />
            <StatCard number="500+" label="Certified Contractors" />
            <StatCard number="$2M" label="Revenue Tracked" />
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <div className="max-w-3xl mx-auto">
          <h3 className="text-4xl font-bold text-gray-900 mb-6">
            Ready to transform your trade business?
          </h3>
          <p className="text-xl text-gray-600 mb-8">
            Join hundreds of contractors already saving time and money with TRADEMASTER.
          </p>
          <Link
            href="/signup"
            className="px-8 py-4 bg-primary-600 text-white rounded-lg hover:bg-primary-700 inline-flex items-center gap-2 text-lg font-semibold"
          >
            Start Your Free Trial
            <ArrowRight className="h-5 w-5" />
          </Link>
          <p className="text-sm text-gray-500 mt-4">
            No credit card required. 14-day free trial.
          </p>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t bg-gray-50 py-12">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Zap className="h-6 w-6 text-primary-600" />
                <h4 className="font-bold text-gray-900">TRADEMASTER</h4>
              </div>
              <p className="text-gray-600 text-sm">
                The future of skilled trades management and training.
              </p>
            </div>
            <div>
              <h5 className="font-semibold text-gray-900 mb-4">Product</h5>
              <ul className="space-y-2 text-sm text-gray-600">
                <li><Link href="/features">Features</Link></li>
                <li><Link href="/pricing">Pricing</Link></li>
                <li><Link href="/training">Training</Link></li>
                <li><Link href="/marketplace">Marketplace</Link></li>
              </ul>
            </div>
            <div>
              <h5 className="font-semibold text-gray-900 mb-4">Company</h5>
              <ul className="space-y-2 text-sm text-gray-600">
                <li><Link href="/about">About</Link></li>
                <li><Link href="/careers">Careers</Link></li>
                <li><Link href="/contact">Contact</Link></li>
                <li><Link href="/blog">Blog</Link></li>
              </ul>
            </div>
            <div>
              <h5 className="font-semibold text-gray-900 mb-4">Legal</h5>
              <ul className="space-y-2 text-sm text-gray-600">
                <li><Link href="/privacy">Privacy</Link></li>
                <li><Link href="/terms">Terms</Link></li>
                <li><Link href="/security">Security</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t mt-8 pt-8 text-center text-sm text-gray-600">
            © 2025 TRADEMASTER. All rights reserved.
          </div>
        </div>
      </footer>
    </main>
  )
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode
  title: string
  description: string
}) {
  return (
    <div className="bg-white p-6 rounded-xl shadow-sm border hover:shadow-md transition-shadow">
      <div className="mb-4">{icon}</div>
      <h4 className="text-xl font-semibold text-gray-900 mb-2">{title}</h4>
      <p className="text-gray-600">{description}</p>
    </div>
  )
}

function StatCard({ number, label }: { number: string; label: string }) {
  return (
    <div>
      <div className="text-4xl font-bold mb-2">{number}</div>
      <div className="text-primary-100">{label}</div>
    </div>
  )
}
