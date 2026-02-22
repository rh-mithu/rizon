const BASE_URL = 'https://localhost:8080'; // Placeholder, should be ENV variable later

export const api = {
    requestLoginLink: async (email: string): Promise<void> => {
        try {
            const response = await fetch(`${BASE_URL}/auth/request-link`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Failed to request login link');
            }
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    },
};
