-- TRADEMASTER Seed Data
-- Version: 002
-- Description: Inserts demo users, projects, and material prices

-- Insert demo users (password is 'password123' hashed with bcrypt cost 10)
-- $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
INSERT INTO users (id, email, password_hash, role) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'demo@trademaster.io', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'contractor'),
('660e8400-e29b-41d4-a716-446655440001', 'john.electrician@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'contractor'),
('770e8400-e29b-41d4-a716-446655440002', 'sarah.hvac@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'contractor'),
('880e8400-e29b-41d4-a716-446655440003', 'client1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'client'),
('990e8400-e29b-41d4-a716-446655440004', 'admin@trademaster.io', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin')
ON CONFLICT (id) DO NOTHING;

-- Insert user profiles
INSERT INTO user_profiles (user_id, first_name, last_name, company, trade, phone, city, state, years_experience) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Demo', 'User', 'Demo Electrical', 'electrical', '+1-555-0100', 'San Francisco', 'CA', 5),
('660e8400-e29b-41d4-a716-446655440001', 'John', 'Smith', 'Smith Electrical Services', 'electrical', '+1-555-0101', 'Austin', 'TX', 10),
('770e8400-e29b-41d4-a716-446655440002', 'Sarah', 'Johnson', 'Johnson HVAC Solutions', 'hvac', '+1-555-0102', 'Phoenix', 'AZ', 8),
('880e8400-e29b-41d4-a716-446655440003', 'Mike', 'Client', 'TechCorp Inc', NULL, '+1-555-0103', 'San Francisco', 'CA', NULL),
('990e8400-e29b-41d4-a716-446655440004', 'Admin', 'User', 'TRADEMASTER', NULL, '+1-555-0104', 'San Francisco', 'CA', NULL)
ON CONFLICT DO NOTHING;

-- Insert certifications
INSERT INTO certifications (user_id, name, issuer, issued_date, expiry_date, verified) VALUES
('660e8400-e29b-41d4-a716-446655440001', 'Master Electrician License', 'Texas State Board', '2020-01-15', '2025-01-15', true),
('660e8400-e29b-41d4-a716-446655440001', 'OSHA 30-Hour Construction', 'OSHA', '2021-03-20', NULL, true),
('770e8400-e29b-41d4-a716-446655440002', 'EPA Universal Certification', 'EPA', '2019-05-10', '2024-05-10', true),
('770e8400-e29b-41d4-a716-446655440002', 'NATE Certification', 'NATE', '2020-08-15', '2025-08-15', true)
ON CONFLICT DO NOTHING;

-- Insert demo projects
INSERT INTO projects (id, user_id, name, description, status, trade, address, city, state, zip, start_date, end_date, budget, progress) VALUES
('aa0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'Office Building Rewire', 'Complete electrical rewiring of 5-story office building', 'active', 'electrical', '123 Main St', 'San Francisco', 'CA', '94105', '2024-01-15', '2024-03-30', 125000.00, 45),
('bb0e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', 'Data Center Installation', 'New data center electrical infrastructure', 'planning', 'electrical', '456 Tech Blvd', 'Austin', 'TX', '78701', '2024-03-01', '2024-05-15', 250000.00, 0),
('cc0e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440002', 'Commercial HVAC Upgrade', 'Replace HVAC system in shopping mall', 'completed', 'hvac', '789 Mall Dr', 'Phoenix', 'AZ', '85001', '2023-11-01', '2024-01-15', 180000.00, 100)
ON CONFLICT (id) DO NOTHING;

-- Insert material prices (Electrical)
INSERT INTO material_prices (name, category, unit, price, trade, supplier) VALUES
('12/2 NM-B Cable (250ft roll)', 'wiring', 'roll', 89.99, 'electrical', 'ElectricSupply Co'),
('14/2 NM-B Cable (250ft roll)', 'wiring', 'roll', 65.50, 'electrical', 'ElectricSupply Co'),
('10/3 NM-B Cable (250ft roll)', 'wiring', 'roll', 145.00, 'electrical', 'ElectricSupply Co'),
('15A Circuit Breaker', 'breakers', 'each', 12.50, 'electrical', 'ElectricSupply Co'),
('20A Circuit Breaker', 'breakers', 'each', 15.75, 'electrical', 'ElectricSupply Co'),
('30A Circuit Breaker', 'breakers', 'each', 22.00, 'electrical', 'ElectricSupply Co'),
('Duplex Outlet (15A)', 'outlets', 'each', 0.85, 'electrical', 'ElectricSupply Co'),
('GFCI Outlet (15A)', 'outlets', 'each', 18.50, 'electrical', 'ElectricSupply Co'),
('Single Pole Switch', 'switches', 'each', 1.25, 'electrical', 'ElectricSupply Co'),
('3-Way Switch', 'switches', 'each', 2.50, 'electrical', 'ElectricSupply Co'),
('Wire Nut (100 pack)', 'connectors', 'pack', 8.99, 'electrical', 'ElectricSupply Co'),
('Junction Box (4-inch)', 'boxes', 'each', 2.75, 'electrical', 'ElectricSupply Co'),
('Junction Box (4-inch square)', 'boxes', 'each', 3.25, 'electrical', 'ElectricSupply Co'),
('Light Fixture (LED 4ft)', 'fixtures', 'each', 45.00, 'electrical', 'ElectricSupply Co'),
('Recessed Can Light', 'fixtures', 'each', 28.50, 'electrical', 'ElectricSupply Co')
ON CONFLICT DO NOTHING;

-- Insert material prices (HVAC)
INSERT INTO material_prices (name, category, unit, price, trade, supplier) VALUES
('R-410A Refrigerant (25lb cylinder)', 'refrigerant', 'cylinder', 425.00, 'hvac', 'HVAC Wholesale'),
('3 Ton AC Condenser Unit', 'units', 'each', 1850.00, 'hvac', 'HVAC Wholesale'),
('4 Ton AC Condenser Unit', 'units', 'each', 2250.00, 'hvac', 'HVAC Wholesale'),
('5 Ton AC Condenser Unit', 'units', 'each', 2850.00, 'hvac', 'HVAC Wholesale'),
('Air Handler (3 Ton)', 'units', 'each', 1250.00, 'hvac', 'HVAC Wholesale'),
('Thermostat (Programmable)', 'controls', 'each', 125.00, 'hvac', 'HVAC Wholesale'),
('Smart Thermostat', 'controls', 'each', 245.00, 'hvac', 'HVAC Wholesale'),
('Ductwork (per foot)', 'duct', 'foot', 12.50, 'hvac', 'HVAC Wholesale'),
('Flex Duct (6-inch, 25ft)', 'duct', 'roll', 35.00, 'hvac', 'HVAC Wholesale'),
('Register (4x10)', 'vents', 'each', 8.50, 'hvac', 'HVAC Wholesale'),
('Return Grille (20x20)', 'vents', 'each', 45.00, 'hvac', 'HVAC Wholesale'),
('Condensate Pump', 'accessories', 'each', 95.00, 'hvac', 'HVAC Wholesale'),
('Air Filter (16x25x1)', 'filters', 'each', 8.99, 'hvac', 'HVAC Wholesale'),
('UV Light System', 'accessories', 'each', 485.00, 'hvac', 'HVAC Wholesale')
ON CONFLICT DO NOTHING;

-- Insert material prices (Welding)
INSERT INTO material_prices (name, category, unit, price, trade, supplier) VALUES
('Steel Plate (4x8, 1/4 inch)', 'materials', 'sheet', 185.00, 'welding', 'Metal Supply Co'),
('Steel Pipe (2 inch, 20ft)', 'materials', 'length', 95.00, 'welding', 'Metal Supply Co'),
('Steel Angle Iron (2x2, 20ft)', 'materials', 'length', 45.00, 'welding', 'Metal Supply Co'),
('Welding Rod E7018 (10lb)', 'consumables', 'box', 45.00, 'welding', 'Welding Supply'),
('MIG Wire (11lb spool)', 'consumables', 'spool', 55.00, 'welding', 'Welding Supply'),
('TIG Filler Rod (10lb)', 'consumables', 'box', 85.00, 'welding', 'Welding Supply'),
('Welding Gas (Argon, 125cf)', 'gas', 'cylinder', 125.00, 'welding', 'Welding Supply'),
('Mixed Gas (75/25, 125cf)', 'gas', 'cylinder', 95.00, 'welding', 'Welding Supply'),
('Grinding Disc (10 pack)', 'tools', 'pack', 25.00, 'welding', 'Welding Supply'),
('Cut-off Wheel (10 pack)', 'tools', 'pack', 22.00, 'welding', 'Welding Supply')
ON CONFLICT DO NOTHING;

-- Insert contractor reputation
INSERT INTO contractor_reputation (contractor_id, overall_rating, total_jobs, completed_jobs, completion_rate, on_time_rate, response_rate) VALUES
('550e8400-e29b-41d4-a716-446655440000', 4.5, 12, 11, 0.9167, 0.9091, 0.95),
('660e8400-e29b-41d4-a716-446655440001', 4.8, 25, 24, 0.96, 0.95, 0.98),
('770e8400-e29b-41d4-a716-446655440002', 4.7, 18, 18, 1.0, 0.94, 0.96)
ON CONFLICT (contractor_id) DO NOTHING;

-- Insert demo jobs
INSERT INTO jobs (id, client_id, title, description, trade, budget, budget_type, city, state, urgency, duration, status) VALUES
('dd0e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440003', 'Data Center Electrical Installation', 'Need experienced electrician for complete electrical system installation in new data center facility.', 'electrical', 45000.00, 'fixed', 'San Francisco', 'CA', 'urgent', '6 weeks', 'open'),
('ee0e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440003', 'Office HVAC Maintenance', 'Regular maintenance and inspection of office building HVAC system.', 'hvac', 5000.00, 'fixed', 'San Francisco', 'CA', 'normal', '1 week', 'open')
ON CONFLICT (id) DO NOTHING;
