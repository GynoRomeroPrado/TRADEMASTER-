'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Upload, Plus, Trash2, Save, Brain } from 'lucide-react'

export default function NewEstimatePage() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [analyzingBlueprint, setAnalyzingBlueprint] = useState(false)

  const [formData, setFormData] = useState({
    projectName: '',
    clientName: '',
    clientEmail: '',
    clientPhone: '',
    type: 'electrical',
    laborHours: '',
    laborRate: '75',
  })

  const [items, setItems] = useState([
    { id: '1', name: '', quantity: '', unit: 'each', unitPrice: '', category: 'material' }
  ])

  const addItem = () => {
    setItems([...items, {
      id: Date.now().toString(),
      name: '',
      quantity: '',
      unit: 'each',
      unitPrice: '',
      category: 'material'
    }])
  }

  const removeItem = (id: string) => {
    setItems(items.filter(item => item.id !== id))
  }

  const updateItem = (id: string, field: string, value: string) => {
    setItems(items.map(item =>
      item.id === id ? { ...item, [field]: value } : item
    ))
  }

  const handleBlueprintUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setAnalyzingBlueprint(true)
    try {
      // TODO: Upload and analyze blueprint
      const formData = new FormData()
      formData.append('file', file)

      const response = await fetch('/api/blueprints/analyze', {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) throw new Error('Analysis failed')

      const analysis = await response.json()
      // Auto-fill items based on analysis
      console.log('Blueprint analysis:', analysis)
    } catch (error) {
      console.error(error)
      alert('Failed to analyze blueprint')
    } finally {
      setAnalyzingBlueprint(false)
    }
  }

  const calculateTotals = () => {
    const laborTotal = parseFloat(formData.laborHours || '0') * parseFloat(formData.laborRate || '0')
    const materialsTotal = items.reduce((sum, item) => {
      const itemTotal = parseFloat(item.quantity || '0') * parseFloat(item.unitPrice || '0')
      return sum + itemTotal
    }, 0)
    const subtotal = laborTotal + materialsTotal
    const tax = subtotal * 0.08 // 8% tax
    const total = subtotal + tax

    return { laborTotal, materialsTotal, subtotal, tax, total }
  }

  const totals = calculateTotals()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const estimate = {
        ...formData,
        items,
        ...totals,
      }

      const response = await fetch('/api/estimates', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(estimate),
      })

      if (!response.ok) throw new Error('Failed to create estimate')

      router.push('/estimates')
    } catch (error) {
      console.error(error)
      alert('Failed to create estimate')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center gap-4">
            <Link href="/estimates" className="p-2 hover:bg-gray-100 rounded-lg">
              <ArrowLeft className="h-5 w-5" />
            </Link>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">New Estimate</h1>
              <p className="text-sm text-gray-600">Create a new project estimate</p>
            </div>
          </div>
        </div>
      </header>

      {/* Form */}
      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* AI Blueprint Analysis */}
          <div className="bg-gradient-to-r from-primary-50 to-blue-50 p-6 rounded-lg border-2 border-primary-200">
            <div className="flex items-center gap-3 mb-4">
              <Brain className="h-6 w-6 text-primary-600" />
              <h2 className="text-lg font-semibold text-gray-900">AI-Powered Blueprint Analysis</h2>
            </div>
            <p className="text-sm text-gray-700 mb-4">
              Upload a PDF blueprint and let our AI automatically calculate materials, labor hours, and costs with 90%+ accuracy.
            </p>
            <div className="flex items-center gap-4">
              <label className="flex-1 cursor-pointer">
                <div className="flex items-center justify-center gap-2 px-4 py-3 bg-white border-2 border-dashed border-primary-300 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors">
                  <Upload className="h-5 w-5 text-primary-600" />
                  <span className="font-medium text-primary-700">
                    {analyzingBlueprint ? 'Analyzing...' : 'Upload Blueprint PDF'}
                  </span>
                </div>
                <input
                  type="file"
                  accept=".pdf"
                  className="hidden"
                  onChange={handleBlueprintUpload}
                  disabled={analyzingBlueprint}
                />
              </label>
              {analyzingBlueprint && (
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              )}
            </div>
          </div>

          {/* Basic Info */}
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Basic Information</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">Project Name</label>
                <input
                  type="text"
                  required
                  value={formData.projectName}
                  onChange={(e) => setFormData({ ...formData, projectName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Client Name</label>
                <input
                  type="text"
                  required
                  value={formData.clientName}
                  onChange={(e) => setFormData({ ...formData, clientName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Client Email</label>
                <input
                  type="email"
                  value={formData.clientEmail}
                  onChange={(e) => setFormData({ ...formData, clientEmail: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Client Phone</label>
                <input
                  type="tel"
                  value={formData.clientPhone}
                  onChange={(e) => setFormData({ ...formData, clientPhone: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Trade Type</label>
                <select
                  value={formData.type}
                  onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                >
                  <option value="electrical">Electrical</option>
                  <option value="hvac">HVAC</option>
                  <option value="welding">Welding</option>
                </select>
              </div>
            </div>
          </div>

          {/* Labor */}
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Labor</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Hours</label>
                <input
                  type="number"
                  step="0.5"
                  value={formData.laborHours}
                  onChange={(e) => setFormData({ ...formData, laborHours: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Rate ($/hr)</label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.laborRate}
                  onChange={(e) => setFormData({ ...formData, laborRate: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Total</label>
                <input
                  type="text"
                  readOnly
                  value={`$${totals.laborTotal.toFixed(2)}`}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 font-semibold"
                />
              </div>
            </div>
          </div>

          {/* Materials */}
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Materials</h2>
              <button
                type="button"
                onClick={addItem}
                className="inline-flex items-center gap-2 px-3 py-1 text-sm bg-primary-600 text-white rounded-lg hover:bg-primary-700"
              >
                <Plus className="h-4 w-4" />
                Add Item
              </button>
            </div>
            <div className="space-y-3">
              {items.map((item) => (
                <div key={item.id} className="grid grid-cols-12 gap-2 items-center">
                  <input
                    type="text"
                    placeholder="Item name"
                    value={item.name}
                    onChange={(e) => updateItem(item.id, 'name', e.target.value)}
                    className="col-span-4 px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                  <input
                    type="number"
                    placeholder="Qty"
                    value={item.quantity}
                    onChange={(e) => updateItem(item.id, 'quantity', e.target.value)}
                    className="col-span-2 px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                  <select
                    value={item.unit}
                    onChange={(e) => updateItem(item.id, 'unit', e.target.value)}
                    className="col-span-2 px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  >
                    <option value="each">Each</option>
                    <option value="ft">Feet</option>
                    <option value="m">Meters</option>
                    <option value="box">Box</option>
                    <option value="roll">Roll</option>
                  </select>
                  <input
                    type="number"
                    step="0.01"
                    placeholder="Unit price"
                    value={item.unitPrice}
                    onChange={(e) => updateItem(item.id, 'unitPrice', e.target.value)}
                    className="col-span-3 px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                  <button
                    type="button"
                    onClick={() => removeItem(item.id)}
                    className="col-span-1 p-2 text-red-600 hover:bg-red-50 rounded-lg"
                  >
                    <Trash2 className="h-5 w-5" />
                  </button>
                </div>
              ))}
            </div>
          </div>

          {/* Totals */}
          <div className="bg-white p-6 rounded-lg shadow-sm">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Summary</h2>
            <div className="space-y-2">
              <div className="flex justify-between text-gray-700">
                <span>Labor</span>
                <span className="font-medium">${totals.laborTotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-gray-700">
                <span>Materials</span>
                <span className="font-medium">${totals.materialsTotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-gray-700 pt-2 border-t">
                <span>Subtotal</span>
                <span className="font-medium">${totals.subtotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-gray-700">
                <span>Tax (8%)</span>
                <span className="font-medium">${totals.tax.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-xl font-bold text-gray-900 pt-2 border-t-2">
                <span>Total</span>
                <span>${totals.total.toFixed(2)}</span>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3">
            <Link
              href="/estimates"
              className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center gap-2 px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50"
            >
              <Save className="h-5 w-5" />
              {loading ? 'Creating...' : 'Create Estimate'}
            </button>
          </div>
        </form>
      </main>
    </div>
  )
}
