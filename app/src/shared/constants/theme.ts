export const theme = {
    colors: {
        background: '#0F172A', // Deep dark blue/black
        surface: '#1E293B', // Slightly lighter for cards/inputs
        primary: '#6366F1', // Indigo
        primaryGradient: ['#6366F1', '#8B5CF6'] as const, // Indigo to Violet
        text: {
            primary: '#F8FAFC', // White-ish
            secondary: '#94A3B8', // Gray-ish
            accent: '#818CF8', // Light Indigo
        },
        border: '#334155',
        error: '#EF4444',
        success: '#10B981',
    },
    spacing: {
        xs: 4,
        s: 8,
        m: 16,
        l: 24,
        xl: 32,
        xxl: 48,
    },
    borderRadius: {
        s: 8,
        m: 12,
        l: 16,
        full: 9999,
    },
    typography: {
        size: {
            xs: 12,
            s: 14,
            m: 16,
            l: 20,
            xl: 24,
            xxl: 32,
        },
        weight: {
            regular: '400' as const,
            medium: '500' as const,
            bold: '700' as const,
        },
    },
};
