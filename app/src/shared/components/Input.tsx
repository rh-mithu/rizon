import React, { useState } from 'react';
import { View, TextInput, Text, StyleSheet, TextInputProps, TouchableOpacity } from 'react-native';
import { theme } from '../constants/theme';
import { LucideIcon, Eye, EyeOff } from 'lucide-react-native';

interface InputProps extends TextInputProps {
    label?: string;
    error?: string;
    icon?: LucideIcon;
    isPassword?: boolean;
}

export const Input: React.FC<InputProps> = ({
    label,
    error,
    icon: Icon,
    isPassword = false,
    style,
    ...props
}) => {
    const [isFocused, setIsFocused] = useState(false);
    const [showPassword, setShowPassword] = useState(!isPassword);

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return (
        <View style={[styles.container, style]}>
            {label && <Text style={styles.label}>{label}</Text>}
            <View
                style={[
                    styles.inputContainer,
                    isFocused && styles.inputContainerFocused,
                    !!error && styles.inputContainerError,
                ]}
            >
                {Icon && (
                    <Icon
                        size={20}
                        color={isFocused ? theme.colors.primary : theme.colors.text.secondary}
                        style={styles.icon}
                    />
                )}
                <TextInput
                    style={styles.input}
                    placeholderTextColor={theme.colors.text.secondary}
                    secureTextEntry={isPassword && !showPassword}
                    onFocus={() => setIsFocused(true)}
                    onBlur={() => setIsFocused(false)}
                    selectionColor={theme.colors.primary}
                    {...props}
                />
                {isPassword && (
                    <TouchableOpacity onPress={togglePasswordVisibility} style={styles.eyeIcon}>
                        {showPassword ? (
                            <EyeOff size={20} color={theme.colors.text.secondary} />
                        ) : (
                            <Eye size={20} color={theme.colors.text.secondary} />
                        )}
                    </TouchableOpacity>
                )}
            </View>
            {error && <Text style={styles.errorText}>{error}</Text>}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        marginBottom: theme.spacing.m,
    },
    label: {
        color: theme.colors.text.secondary,
        marginBottom: theme.spacing.xs,
        fontSize: theme.typography.size.s,
        fontWeight: theme.typography.weight.medium,
    },
    inputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: theme.colors.surface,
        borderRadius: theme.borderRadius.m,
        borderWidth: 1,
        borderColor: theme.colors.border,
        height: 56,
        paddingHorizontal: theme.spacing.m,
    },
    inputContainerFocused: {
        borderColor: theme.colors.primary,
        backgroundColor: '#1E293B', // Keep surface color but maybe slightly lighter if needed
    },
    inputContainerError: {
        borderColor: theme.colors.error,
    },
    icon: {
        marginRight: theme.spacing.s,
    },
    input: {
        flex: 1,
        color: theme.colors.text.primary,
        fontSize: theme.typography.size.m,
        height: '100%',
    },
    eyeIcon: {
        marginLeft: theme.spacing.s,
    },
    errorText: {
        color: theme.colors.error,
        fontSize: theme.typography.size.xs,
        marginTop: theme.spacing.xs,
    },
});
