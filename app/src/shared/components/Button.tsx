import React from 'react';
import { TouchableOpacity, Text, StyleSheet, ActivityIndicator, ViewStyle } from 'react-native';
import { LinearGradient } from 'expo-linear-gradient';
import { theme } from '../constants/theme';

interface ButtonProps {
    title: string;
    onPress: () => void;
    variant?: 'primary' | 'secondary' | 'outline';
    loading?: boolean;
    disabled?: boolean;
    style?: ViewStyle;
}

export const Button: React.FC<ButtonProps> = ({
    title,
    onPress,
    variant = 'primary',
    loading = false,
    disabled = false,
    style,
}) => {
    const isPrimary = variant === 'primary';
    const isOutline = variant === 'outline';

    const content = (
        <>
            {loading ? (
                <ActivityIndicator color={isOutline ? theme.colors.primary : '#FFF'} />
            ) : (
                <Text
                    style={[
                        styles.text,
                        isOutline && styles.textOutline,
                        disabled && styles.textDisabled,
                    ]}
                >
                    {title}
                </Text>
            )}
        </>
    );

    if (isPrimary && !disabled && !isOutline) {
        return (
            <TouchableOpacity
                onPress={onPress}
                activeOpacity={0.8}
                disabled={loading || disabled}
                style={[styles.container, style]}
            >
                <LinearGradient
                    colors={theme.colors.primaryGradient}
                    start={{ x: 0, y: 0 }}
                    end={{ x: 1, y: 1 }}
                    style={styles.gradient}
                >
                    {content}
                </LinearGradient>
            </TouchableOpacity>
        );
    }

    return (
        <TouchableOpacity
            onPress={onPress}
            activeOpacity={0.8}
            disabled={loading || disabled}
            style={[
                styles.container,
                styles.buttonBase,
                isOutline && styles.buttonOutline,
                disabled && styles.buttonDisabled,
                style,
            ]}
        >
            {content}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    container: {
        borderRadius: theme.borderRadius.m,
        overflow: 'hidden',
        height: 56,
    },
    gradient: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    buttonBase: {
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: theme.colors.surface,
    },
    buttonOutline: {
        backgroundColor: 'transparent',
        borderWidth: 1,
        borderColor: theme.colors.primary,
    },
    buttonDisabled: {
        backgroundColor: theme.colors.surface,
        opacity: 0.5,
    },
    text: {
        color: '#FFF',
        fontSize: theme.typography.size.m,
        fontWeight: theme.typography.weight.bold,
    },
    textOutline: {
        color: theme.colors.primary,
    },
    textDisabled: {
        color: theme.colors.text.secondary,
    },
});
