import React from 'react';
import { View, Text, Button } from 'react-native';
import { useAuth } from '../auth/auth.context';

export default function HomeScreen() {
  const { logout } = useAuth();

  return (
    <View>
      <Text>Home Screen</Text>
      <Button title="Logout" onPress={logout} />
    </View>
  );
}
