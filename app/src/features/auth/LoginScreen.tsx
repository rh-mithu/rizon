import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, KeyboardAvoidingView, Platform, Alert } from 'react-native';
import { useAuth } from './auth.context';
import { ScreenContainer } from '../../shared/components/ScreenContainer';
import { Input } from '../../shared/components/Input';
import { Button } from '../../shared/components/Button';
import { theme } from '../../shared/constants/theme';
import { Mail, ArrowLeft, MailCheck } from 'lucide-react-native';

export default function LoginScreen() {
  const { requestLoginLink } = useAuth();
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isEmailSent, setIsEmailSent] = useState(false);

  const handleLogin = async () => {
    if (!email) {
      Alert.alert('Error', 'Please enter your email address');
      return;
    }

    setIsLoading(true);
    try {
      await requestLoginLink(email);
      setIsEmailSent(true);
    } catch (error) {
      Alert.alert('Error', error instanceof Error ? error.message : 'Failed to send login link');
    } finally {
      setIsLoading(false);
    }
  };

  const handleBackToLogin = () => {
    setIsEmailSent(false);
    setEmail('');
  };

  if (isEmailSent) {
    return (
      <ScreenContainer style={styles.container}>
        <View style={styles.content}>
          <View style={styles.header}>
            <View style={styles.iconContainer}>
              <MailCheck size={48} color={theme.colors.primary} />
            </View>
            <Text style={styles.title}>Check your email</Text>
            <Text style={[styles.subtitle, styles.centerText]}>
              We sent a login link to{'\n'}
              <Text style={styles.emailText}>{email}</Text>
            </Text>
          </View>

          <View style={styles.form}>
            <Button
              title="Open Email App"
              onPress={() => {
                // In a real app, this would open the email client
                Alert.alert('Info', 'Opening email app...');
              }}
              variant="outline"
              style={styles.emailButton}
            />

            <TouchableOpacity onPress={handleBackToLogin} style={styles.backButton}>
              <ArrowLeft size={16} color={theme.colors.text.secondary} style={styles.backIcon} />
              <Text style={styles.backText}>Back to login</Text>
            </TouchableOpacity>
          </View>
        </View>
      </ScreenContainer>
    );
  }

  return (
    <ScreenContainer style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.content}
      >
        <View style={styles.header}>
          <Text style={styles.title}>Welcome Back</Text>
          <Text style={styles.subtitle}>Sign in with your email to continue</Text>
        </View>

        <View style={styles.form}>
          <Input
            label="Email"
            placeholder="hello@example.com"
            value={email}
            onChangeText={setEmail}
            autoCapitalize="none"
            keyboardType="email-address"
            icon={Mail}
          />

          <Button
            title="Send Login Link"
            onPress={handleLogin}
            loading={isLoading}
            style={styles.loginButton}
          />
        </View>
      </KeyboardAvoidingView>
    </ScreenContainer>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingHorizontal: theme.spacing.l,
  },
  content: {
    flex: 1,
    justifyContent: 'center',
  },
  header: {
    marginBottom: theme.spacing.xl,
    alignItems: 'center',
  },
  title: {
    fontSize: theme.typography.size.xxl,
    fontWeight: theme.typography.weight.bold,
    color: theme.colors.text.primary,
    marginBottom: theme.spacing.s,
    textAlign: 'center',
  },
  subtitle: {
    fontSize: theme.typography.size.m,
    color: theme.colors.text.secondary,
    textAlign: 'center',
  },
  form: {
    marginBottom: theme.spacing.xl,
  },
  loginButton: {
    marginTop: theme.spacing.s,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
  },
  footerText: {
    color: theme.colors.text.secondary,
    fontSize: theme.typography.size.s,
  },
  signupText: {
    color: theme.colors.text.accent,
    fontSize: theme.typography.size.s,
    fontWeight: theme.typography.weight.bold,
  },
  iconContainer: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: 'rgba(99, 102, 241, 0.1)',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: theme.spacing.l,
  },
  centerText: {
    textAlign: 'center',
  },
  emailText: {
    color: theme.colors.text.primary,
    fontWeight: theme.typography.weight.bold,
  },
  emailButton: {
    marginBottom: theme.spacing.l,
  },
  backButton: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
  },
  backIcon: {
    marginRight: theme.spacing.xs,
  },
  backText: {
    color: theme.colors.text.secondary,
    fontSize: theme.typography.size.s,
  },
});
