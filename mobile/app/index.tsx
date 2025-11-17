import { View, Text, TouchableOpacity, StyleSheet } from 'react-native'
import { Link } from 'expo-router'

export default function HomeScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>TRADEMASTER</Text>
      <Text style={styles.subtitle}>Field Management & Training</Text>

      <View style={styles.buttonsContainer}>
        <Link href="/projects" asChild>
          <TouchableOpacity style={styles.button}>
            <Text style={styles.buttonText}>My Projects</Text>
          </TouchableOpacity>
        </Link>

        <Link href="/camera" asChild>
          <TouchableOpacity style={[styles.button, styles.primaryButton]}>
            <Text style={[styles.buttonText, styles.primaryButtonText]}>
              Capture Progress
            </Text>
          </TouchableOpacity>
        </Link>

        <Link href="/training" asChild>
          <TouchableOpacity style={styles.button}>
            <Text style={styles.buttonText}>Training</Text>
          </TouchableOpacity>
        </Link>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f9fafb',
    padding: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#0ea5e9',
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: '#6b7280',
    marginBottom: 40,
  },
  buttonsContainer: {
    width: '100%',
    maxWidth: 300,
    gap: 16,
  },
  button: {
    backgroundColor: 'white',
    padding: 16,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: '#e5e7eb',
  },
  primaryButton: {
    backgroundColor: '#0ea5e9',
    borderColor: '#0ea5e9',
  },
  buttonText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#374151',
    textAlign: 'center',
  },
  primaryButtonText: {
    color: 'white',
  },
})
